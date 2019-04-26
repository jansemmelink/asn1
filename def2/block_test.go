package def2_test

import (
	"testing"
	"math/rand"
	"fmt"
	//"reflect"

	"github.com/jansemmelink/asn1/def2"
	"github.com/jansemmelink/asn1/parser"
	//"github.com/jansemmelink/mem"
)

//var log = logger.New().WithLevel(level.Debug)

func TestBlock(t *testing.T) {
	l := parser.NewLines()

	//append a few consecutive braced uint nrs:
	nrTests := 5
	nr := make([]uint,nrTests)
	for i := 0; i < nrTests; i ++ {
		nr[i] = uint(rand.Uint64())
		l = l.Append(i+1, fmt.Sprintf("{%d}", nr[i]), "no comment")
	}

	remain := l
	u := BracedUint{}
	for i := 0; i < nrTests; i ++ {
		var err error
		remain, err = u.Parse(log, remain, &u)
		if err != nil || uint(u.Nr) != nr[i] {
			t.Fatalf("Failed[%d]: u=%+v  (expected %d), next=%s, %v", i, u, nr[i], remain.Next(), err)
		}
		log.Debugf("Success: Parsed[%d]: {%d}", i, u.Uint())
	}//for each test nr
	return
}

//MyBlock parses a uint in braces, e.g. {5}
type BracedUint struct {
	def2.Block
	Nr def2.Uint
}

func (b BracedUint) Start() def2.IToken { return def2.NewToken("{") }
func (b BracedUint) End() def2.IToken   { return def2.NewToken("}") }
func (b BracedUint) Uint() uint         { return uint(b.Nr) }


func TestBlockValue(t *testing.T) {
	tests := []struct {text string; sel int; u uint; i int; n string; } {
		{"1234", 1, 1234, 0, ""},
		{"-54321", 2, 0, -54321, ""},
		{"Jan", 3, 0, 0, "Jan"},
		{"-InvalidID", -1, 0, 0, ""},
	}

	for i,tt := range tests {
		//write the block as "[ <value> ]":
		l := parser.NewLines()
		l = l.Append(i+1, fmt.Sprintf("[%s]", tt.text), "test data")

		remain := l
		x := &BlockValue{}
		//x := mem.NewX(reflect.TypeOf(BlockValue{})).(*BlockValue)
		var err error
		remain, err = x.Parse(log, remain, x)
		if err != nil {
			//parsing failed: expected when sel==-1
			if tt.sel < 0 {
				log.Debugf("test[%d]: Parsing failed as expected: %v", i, err)
			} else {
				t.Fatalf("test[%d]: Parsing failed unexpectedly: %v", i, err)
			}
		} else {
			//parsing succeeded: not expexted when sel==-1
			if tt.sel < 0 {
				t.Fatalf("test[%d]: Parsing succeded unexpectedly: %v", i, err)
			} else {
				//parsed as expected: check the parsed selector and value:
				log.Debugf("test[%d]: Parsed \"%s\" -> x = %+v", i, tt.text, x)
				if x.V.Choice.Option().Index != tt.sel {
					t.Fatalf("Parsed selector=%d != tt.sel=%d", x.V.Choice.Option().Index, tt.sel)
				}
				switch tt.sel {
				case 1:
					parsedUint := uint(x.V.U)
					if parsedUint != tt.u || tt.u == 0 {
						t.Fatalf("test[%d]: parsed uint %d expected %d (must not be 0)", i, parsedUint, tt.u)
					}
				case 2:
					parsedInt := int(x.V.I)
					if parsedInt != tt.i || tt.i == 0 {
						t.Fatalf("test[%d]: parsed int %d expected %d (must not be 0)", i, parsedInt, tt.i)
					}
				case 3:
					parsedStr := string(x.V.N)
					if parsedStr != tt.n || tt.n == "" {
						t.Fatalf("test[%d]: parsed name \"%s\" expected \"%s\" (must not be \"\")", i, parsedStr, tt.n)
					}
				default:
					t.Fatalf("Parsed selector=%d", tt.sel)
				}
			}
		}
	}//for each test nr
	return

	// remain := l
	// u := BracedUint{}
	// for i := 0; i < nrTests; i ++ {
	// 	var err error
	// 	remain, err = u.Parse(log, remain, &u)
	// 	if err != nil || uint(u.Nr) != nr[i] {
	// 		t.Fatalf("Failed[%d]: u=%+v  (expected %d), next=%s, %v", i, u, nr[i], remain.Next(), err)
	// 	}
	// 	log.Debugf("Success: Parsed[%d]: {%d}", i, u.Uint())
	// }//for each test nr
	// return
}

//BlockValue parses [ <Value> ]
//which is defined in choice_test.go
type BlockValue struct {
	def2.Block
	V MyValue
}

func (b BlockValue) Start() def2.IToken { return def2.NewToken("[") }
func (b BlockValue) End() def2.IToken   { return def2.NewToken("]") }
