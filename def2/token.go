package def2

import (
	"reflect"
	"strings"

	"bitbucket.org/vservices/dark/logger"
	"github.com/jansemmelink/asn1/parser"
)

//NewToken makes a parsable token
//	A token is a constant single word of text.
//Examples:
//		::=
//		{
//		OPTIONAL
//		[
//		]
//
//If the token is a normal word, e.g. "OPTIONAL",
//	you can use the following shorthand definitions and include it in your Seq struct
//	as explained here, where IToken.parser will use your type name as the token:
//
//		type OPTIONAL = struct { def2.Token }
//
//		type myStruct struct {
//			def2.Seq
//			...
//			Opt OPTIONAL		<-- used here between as part of sequence
//			...
//		}
//
//If the word is a golang reserved word, e.g. "type", then that won't work,
//	but you can just prefix tour type name with "Keyword", then it will still
//	parse everything after that. This example will still parse the word "OPTIONAL":
//
//		type KeywordOPTIONAL = struct { def2.Token }
//
//If the token is not a normal word, e.g. "::=",
//	you cannot use shorthand because it will not be a valid golang type name.
//	In this case, override the Word() method to return your token text:
//
//		type Assignment = struct { def2.Token }
//		func (Assignment) String() string { return "::=" }
//
//		type myStruct struct {
//			def2.Seq
//			...
//			Assign Assignment		<-- used here between as part of sequence
//			...
//		}
//
func NewToken(t string) IToken {
	tClean := strings.Trim(t, " \t")
	if len(tClean) < 1 {
		panic(log.Wrapf(nil, "Token(\"%s\")", t))
	}
	return Token{t: tClean}
}

//Token struct to embed in customer token struct definitions
type Token struct {
	t string
}

//IToken is used in this package, do not embed it in your struct!
type IToken interface {
	IParsable
	String() string
}

//Token ...
func (t Token) String() string {
	return t.t
}

//Parse ...
func (t Token) Parse(log logger.ILogger, l parser.ILines, v IParsable) (parser.ILines, error) {
	//token field is defined when used NewToken(...)
	token := t.t

	if len(token) == 0 {
		//String() method is used to get the token if user provided a custom method:
		t := reflect.ValueOf(v).Elem().Interface().(IToken)
		token = t.String()
	}

	if len(token) == 0 {
		//typename "CustomToken" is used as token when user defined a new type with:
		//	CustomToken = struct {def2.Token}
		token = reflect.TypeOf(v).Elem().Name()
		if len(token) > 0 {
			log.Debugf("Parsing Token(\"%s\") from line %d: %.32s ...", token, l.LineNr(), l.Next())
		}
		//keyword "Type" is used as token when user defined a new type name starting with KeywordXxx:
		//	KeywordType = struct {def2.Token}
		if len(token) > 7 && token[0:7] == "Keyword" {
			token = token[7:]
		}
	}

	if remain, ok := l.SkipOver(token); ok {
		return remain, nil
	}
	return l, log.Wrapf(nil, "Token(\"%s\") not found at line %d: %.32s ...", token, l.LineNr(), l.Next())
}
