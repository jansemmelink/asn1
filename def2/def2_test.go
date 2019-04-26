package def2_test

import (
	"fmt"
	"reflect"
	"testing"

	"bitbucket.org/vservices/dark/logger"
	"bitbucket.org/vservices/dark/logger/level"
	"github.com/jansemmelink/asn1/def2"
	"github.com/jansemmelink/asn1/parser"
)

var log = logger.New().WithLevel(level.Debug)

func TestToken(t *testing.T) {
	{
		l := parser.NewLines()
		l = l.Append(1, "Name is Jan := Semmelink;", "no comment")

		expectedTokens := []def2.IParsable{
			def2.NewToken("Name"),
			def2.NewToken("is"),
			def2.NewToken("Jan"),
			def2.NewToken("Kobus"), //<<--- not expected - added to list to ensure we get parse error
			def2.NewToken(":="),
			def2.NewToken("Semmelink"),
			def2.NewToken("end"), //<<--- also not expected
		}
		remain := l
		for i, token := range expectedTokens {
			var err error
			remain, err = token.Parse(log, remain, token)
			if err != nil && i != 3 && i != 6 { //index 3=Kobus and 6=end should fail
				t.Fatalf("Failed: %v", err)
			}
			if err == nil && (i == 3 || i == 6) {
				t.Fatalf("Failed: %v", err)
				log.Debugf("  Found token[i]=%s", i, token)
			} else {
				log.Debugf("  Expected missing token[i]=%s", i, token)
			}
		}
		log.Debugf("Got all expected tokens")
	}
}

func TestRegex(t *testing.T) {
	//todo
	return
}

func TestNumbers(t *testing.T) {
	//todo
	return
}

func TestSeq(t *testing.T) {
	//todo
	return
}

func TestList(t *testing.T) {
	//todo
	return
}

func TestOptional(t *testing.T) {
	{
		//optional Name2 present:
		l := parser.NewLines()
		l = l.Append(1, "1234 ABC DEF 5678", "no comment")
		x := NrNameOptNameNr{}
		remain, err := x.Parse(log, l, &x)
		if err != nil || x.Nr1 != 1234 || x.Nr2 != 5678 || x.Name1 != "ABC" || x.OptName2 == nil || *x.OptName2 != "DEF" {
			t.Fatalf("Failed: x=%s, next=%s, %v", x, remain.Next(), err)
		}
	}

	{
		//optional Name2 absent:
		l := parser.NewLines()
		l = l.Append(1, "1234 ABC 5678", "no comment")
		x := NrNameOptNameNr{}
		remain, err := x.Parse(log, l, &x)
		if err != nil || x.Nr1 != 1234 || x.Nr2 != 5678 || x.Name1 != "ABC" || x.OptName2 != nil {
			t.Fatalf("Failed: x=%s, next=%s, %v", x, remain.Next(), err)
		}
	}
} //TestOptional()

//NrNameOptNameNr is a sequence of:
//	Nr + Name + Optional Name + Nr
type NrNameOptNameNr struct {
	def2.Seq
	Nr1      def2.Int
	Name1    Identifier
	OptName2 *Identifier //use pointer type to make it optional
	Nr2      def2.Int
}

func (x NrNameOptNameNr) String() string {
	if x.OptName2 == nil {
		return fmt.Sprintf("%d,%s,%s,%d", x.Nr1, x.Name1, "N/A", x.Nr2)
	}
	return fmt.Sprintf("%d,%s,%s,%d", x.Nr1, x.Name1, *x.OptName2, x.Nr2)
}

func TestAll(t *testing.T) {
	//make some text that can be parsed into a struct with two fields:
	//the first field's value is "Outside" the second in "Inside"
	l := parser.NewLines()
	l = l.Append(1, "1234 -456 789 101112 -8 7 Before ::= outside [ Inside 543 55] n1,n2,n3", "no comment")

	//the schema to parse that is a sequence of one word followed by a {...} block containing another word
	//which should result in struct with fields One=Outside and Two=Inside
	var err error

	//custom parser type: IntValue has its own parser in this module
	if true {
		intValue := IntValue{}
		remain := l
		for i := 0; i < 5; i++ {
			remain, err = def2.Parse(log, remain, &intValue)
			if err != nil {
				t.Fatalf("Failed: %v", err)
			}
			log.Debugf("Choice Value: %+v", intValue)
		}
		log.Debugf("Remain: line %d: %.32s ...", remain.LineNr(), remain.Next())
	}

	//choice
	{
		u := MyUnion{}
		remain := l
		for i := 0; i < 6; i++ {
			remain, err = def2.Parse(log, remain, &u)
			if err != nil {
				t.Fatalf("Failed: %v", err)
			}
			log.Debugf("Parsed Choice: %+v   selection: %+v", u, u.Choice.Option())
			//using a union member by name
			switch u.Choice.Option().Name {
			case "Two":
				//can do this, but not necessary if we code each member by hand:
				//iv := u.Choice.Value().(IntValue)
				//log.Debugf("  TWO = %v", iv)
				log.Debugf("  TWO = %v", u.Two)
			default:
				log.Errorf("Cannot log %+v", u.Choice.Option())
			}

			//using a union member by name
			switch u.Choice.Option().Ptr {
			case &u.Two:
				log.Debugf("  TWO = %v", u.Two)
			default:
				log.Errorf("Cannot log %+v", u.Choice.Option())
			}
		}

		log.Debugf("Remain: line %d: %.32s ...", remain.LineNr(), remain.Next())
	} //scope for choice

	//seq
	if true {
		myStruct := MyStruct{}
		var remain parser.ILines
		remain, err = def2.Parse(log, l, &myStruct)
		if err != nil {
			t.Fatalf("Failed: %v", err)
		}
		log.Debugf("2Val: %+v", myStruct)
		log.Debugf("  One   = %d", myStruct.One)
		log.Debugf("  One2  = %d", *myStruct.One2)
		log.Debugf("  TWO   = %d", myStruct.Two.Int())
		log.Debugf("  THREE = %+v", myStruct.Three.Int())
		log.Debugf("  FOUR  = %+v", myStruct.Four)
		log.Debugf("  FIVE  = %+v", myStruct.Five)
		log.Debugf("  I.I   = \"%s\"", myStruct.I.I)
		log.Debugf("  I.V1  = %v", myStruct.I.V1)
		log.Debugf("  I.V2  = %v", myStruct.I.V2)
		log.Debugf("  NL    = %+v", myStruct.NL)
		log.Debugf("  NL.#  = %d", len(myStruct.NL.Items()))
		log.Debugf("Remain: line %d: %.32s ...", remain.LineNr(), remain.Next())
	} //scope for seq
}

//this is a parsable struct that inherits a Parse() method
type MyUnion struct {
	def2.Choice
	Two def2.Int
	//	As  AssignmentToken
	//	Str Identifier
}

//Assignement Token: "::="
type AssignmentToken struct{ def2.Token }

func (AssignmentToken) String() string { return "::=" }

//Identifier is parsed using a regular expression
type Identifier string

func (id *Identifier) Parse(log logger.ILogger, l parser.ILines, v def2.IParsable) (parser.ILines, error) {
	return def2.RegexParse(log, l, "[a-zA-Z][a-zA-Z0-9]*", (*string)(id))
}

//func (i Identifier) String() string { return i.Value }

//this is a parsable sequence struct that inherits a Parse() method from type Seq
type MyStruct struct {
	//embedded struct "Seq" that implements IParsable:
	//	field.type.	kind = struct
	//				Implements(IParsableInterfaceType) = true
	def2.Seq

	//struct field not implementing parsable:
	//	field.type.	kind = ... (string in this case, but may be struct or not, not relevant...)
	//				Implements(IParsableInterfaceType) = false
	one string
	Z   *string

	X OtherStruct

	//non-struct type fields that also implements IParsable:
	One  def2.Int
	One2 *def2.Int

	//struct field using a parsable type: <name> <someParsableStructType>
	//	field.type.	kind = struct
	//				Implements(IParsableInterfaceType) = false (because not a pointer, so recon to be const, but we can set it)
	//		need to make a ptr of this type and test if THAT implements the interface...
	Two IntValue

	Y *OtherStruct

	//struct field using a POINTER TO parsable type: <name> &<someParsableStructType>
	//	field.type.	kind = ptr
	//				Implements(IParsableInterfaceType) = true (because pointer and writable)
	Three *IntValue

	//sub parsed structs
	Four AlsoValues
	Five *AlsoValues

	//keywords expected after above members
	Token1 Before
	Token2 AssignmentToken
	Token3 outside

	//block in braces: {...}
	I InsideBlock

	//list of items
	NL NameList
}

type CSV struct {
	*def2.List
}

func (csv CSV) Sep() def2.IToken {
	return def2.NewToken(",")
}

type NameList struct {
	*CSV
}

// func (NameList) Sep() def2.IToken {
// 	return def2.NewToken(",")
// }

func (NameList) Min() int {
	return 1
}

func (NameList) Max() int {
	return 10
}

func (NameList) ItemType() reflect.Type {
	var name Identifier
	return reflect.TypeOf(name)
}

//InsideBlock defines our own block start..end notation:
type InsideBlock struct {
	def2.Block
	//I  KeywordInside
	I Identifier
	//I def2.Regex
	V1 def2.Int
	V2 def2.Int
}

func (ib InsideBlock) Start() def2.IToken {
	return def2.NewToken("[")
}

func (ib InsideBlock) End() def2.IToken {
	return def2.NewToken("]")
}

type Before struct{ def2.Token }

type outside struct{ def2.Token }

//if type name starts with "Keyword" - that bit is ignored:
type KeywordInside struct{ def2.Token }

//this is a parsable struct that implements its own Parse() method
type IntValue struct {
	int
}

func (i IntValue) Int() int {
	return i.int
}

var tempNextIntValue = 3000

func (iv *IntValue) Parse(log logger.ILogger, l parser.ILines, v def2.IParsable) (parser.ILines, error) {
	tempNextIntValue++
	iv.int = tempNextIntValue
	//log.Debugf("PARSED %T = %d", *iv, iv.i)
	return l, nil
}

type AlsoValues struct {
	def2.Seq
	One def2.Int
	Two def2.Int
}

//general struct that does not implement parsable interface
type OtherStruct struct{}
