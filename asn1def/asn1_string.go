package asn1def

import (
	"strconv"
	"strings"

	"bitbucket.org/vservices/dark/logger"
	"github.com/jansemmelink/asn1/parser"
)

//ICharacterStringList ...
//CharacterStringList ::= " { " CharSyms " } "
type ICharacterStringList interface {
	parser.IBlock
}

//CharacterStringList ...
func CharacterStringList() ICharacterStringList {
	return &characterStringList{
		IBlock: parser.Block("characterStringList", "{", "}", CharSyms()),
	}
}

type characterStringList struct {
	parser.IBlock
}

//ICharSyms ...
//CharSyms ::=
//		CharsDefn
//	|	CharSyms "," CharsDefn
type ICharSyms interface {
	parser.IList
}

//CharSyms ...
func CharSyms() ICharSyms {
	return &charSyms{
		IList: parser.List(",", 1, CharsDefn()),
	}
}

type charSyms struct {
	parser.IList
}

//ICharsDefn ...
//===== Notation =====
// CharsDefn ::=
// 		cstring
//	|	Quadruple
//	|	Tuple
//	|	DefinedValue
type ICharsDefn interface {
	parser.IChoice
}

//CharsDefn ...
func CharsDefn() ICharsDefn {
	return &charsDefn{
		IChoice: parser.Choice("charsDefn",
			CString(),
			//todo: Quadruple(),
			//todo: Tuple(),
			DefinedValue(),
		),
	}
}

type charsDefn struct {
	parser.IChoice
}

//ICString ...
type ICString interface {
	parser.INotation
}

//CString ...
func CString() ICString {
	return &cstring{
		INotation: parser.New("cstring", "cstring"),
	}
}

type cstring struct {
	parser.INotation
}

func (c cstring) ParseV(log logger.ILogger, l parser.ILines) (parser.IValue, parser.ILines, error) {
	//log.Debugf("%T(%s).V() from line %d: %.32s ...", c, c.Name(), l.LineNr(), l.Next())

	//start with a double quotation char '"', then continue until terminated with another
	//but note that two consequtive double quotations are to represent a quotation in the cstring value
	//Example:
	// "ABCDE    FGH
	//	IJK""XYZ"
	//represents a single cstring of value:
	//	ABCDE   FGHIJK"XYZ
	remain, ok := l.SkipOver("\"")
	if !ok {
		return nil, l, log.Wrapf(nil, "cstring not started with '\"' on line %d: %s", l.LineNr(), l.Next())
	}
	value := ""
	for {
		next := remain.Next()
		nextQuotePos := strings.IndexByte(next, '"')
		//log.Debugf("line %d: %.32s     nextQuotePos=%d len=%d", remain.LineNr(), next, nextQuotePos, len(next))
		if nextQuotePos < 0 {
			//no quote in rest of line, include rest of line (trimming trailing spaces)
			//and continue on next line (trimming leading spaces)
			value += strings.TrimRight(next, " ")
			remain, ok = remain.SkipOver(next)
			if !ok {
				return nil, l, log.Wrapf(nil, "cstring failed to skip to end of line on line %d: %s", l.LineNr(), l.Next())
			}
			//log.Debugf("added rest of line, now cstring=%s, remain line %d: %.32s", c.value, remain.LineNr(), remain.Next())
			continue
		}

		//see of consequtive quote chars to be included in value
		if nextQuotePos < len(next) && next[nextQuotePos+1] == '"' {
			value += next[0 : nextQuotePos+1]
			remain, ok = remain.SkipOver(next[0 : nextQuotePos+2])
			//log.Debugf("added text up to c-quote, now cstring=%s, remain line %d: %.32s", c.value, remain.LineNr(), remain.Next())
			if !ok {
				return nil, l, log.Wrapf(nil, "cstring failed to skip to quote-in-value pos=%d on line %d: %s", nextQuotePos, l.LineNr(), l.Next())
			}
			continue
		}

		//not consequtive quotes, i.e. end of cstring
		value += next[0:nextQuotePos]
		remain, ok = remain.SkipOver(next[0 : nextQuotePos+1])
		//log.Debugf("added text up to terminating quote, now cstring=%s, remain line %d: %.32s", c.value, remain.LineNr(), remain.Next())
		if !ok {
			return nil, l, log.Wrapf(nil, "cstring failed to skip to end-of-value")
		}
		break
	} //for cstring value

	return parser.StrValue(c.Name(), value), remain, nil
} //cstring.ParseV()

//IBString ...
//11.10 Binary strings
//Name of lexical item – bstring
//A "bstring" shall consist of an arbitrary number (possibly zero) of the characters: 0 1
//possibly intermixed with white-space,
//preceded by an APOSTROPHE (39) character ( ' ) and followed by the pair of characters: 'B
//EXAMPLE – '01101100'B
//Occurrences of white-space within a binary string lexical item have no significance.
type IBString interface {
	parser.INotation
}

//BString ...
func BString() IBString {
	return &bstring{
		INotation: parser.New("bstring", "bstring"),
	}
}

type bstring struct {
	parser.INotation
}

func (b bstring) ParseV(log logger.ILogger, l parser.ILines) (parser.IValue, parser.ILines, error) {
	next := l.Next()

	//start with an apostrophe
	posStart := strings.IndexByte(next, '\'')
	if posStart != 0 {
		return nil, l, log.Wrapf(nil, "bstring not started, expecting '...'B")
	}

	//max 64-bits = "0000000000000000" = 16 characters but allow lots of spaces, so check up to 64 chars...
	posEnd := strings.Index(next[1:], "'B")
	if posEnd < 0 || posEnd > 64 {
		return nil, l, log.Wrapf(nil, "bstring not ended, expecting '...'B")
	}

	str := next[1 : posEnd-posStart+1]
	//log.Debugf("Parsing binary from start=%d+1 end=%d \"%s\"", posStart, posEnd, str)
	str = strings.Replace(str, " ", "", -1)
	str = strings.Replace(str, "\t", "", -1)
	value, err := strconv.ParseUint(str, 2, 64)
	if err != nil {
		return nil, l, log.Wrapf(err, "bstring=%s is invalid", str)
	}
	remain, _ := l.SkipOver(next[0 : posEnd+3])
	return parser.ValueV(b.Name(), value), remain, nil
} //bstring.ParseV()

//IHString ...
//11.12 Hexadecimal strings
//Name of lexical item – hstring
//11.12.1 An "hstring" shall consist of an arbitrary number (possibly zero) of the characters:
//A B C D E F 0 1 2 3 4 5 6 7 8 9
//possibly intermixed with white-space, preceded by an APOSTROPHE (39) character ( ' ) and followed by the pair of
//characters: 'H
//EXAMPLE – 'AB0196'H
//Occurrences of white-space within a hexadecimal string lexical item have no significance.
type IHString interface {
	parser.INotation
	Value() uint64
}

//HString ...
func HString() IHString {
	return &hstring{
		INotation: parser.New("hstring", "hstring"),
	}
}

type hstring struct {
	parser.INotation
	value uint64
}

func (h hstring) Value() uint64 {
	return h.value
}

func (h hstring) ParseV(log logger.ILogger, l parser.ILines) (parser.IValue, parser.ILines, error) {
	next := l.Next()

	//start with an apostrophe
	posStart := strings.IndexByte(next, '\'')
	if posStart != 0 {
		return nil, l, log.Wrapf(nil, "hstring not started, expecting '...'H")
	}

	//max 64-bits = "0000000000000000" = 16 characters but allow lots of spaces, so check up to 64 chars...
	posEnd := strings.Index(next[1:], "'H")
	if posEnd < 0 || posEnd > 64 {
		return nil, l, log.Wrapf(nil, "hstring not ended, expecting '...'H")
	}

	str := next[1 : posEnd-posStart+1]
	//log.Debugf("Parsing binary from start=%d+1 end=%d \"%s\"", posStart, posEnd, str)
	str = strings.Replace(str, " ", "", -1)
	str = strings.Replace(str, "\t", "", -1)
	value, err := strconv.ParseUint(str, 16, 64)
	if err != nil {
		return nil, l, log.Wrapf(err, "hstring=\"%s\" is invalid", str)
	}
	remain, _ := l.SkipOver(next[0 : posEnd+3])
	return parser.ValueV(h.Name(), value), remain, nil
} //hstring.ParseV()
