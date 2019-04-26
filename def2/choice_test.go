package def2_test

import (
	"testing"
	"reflect"
//	"fmt"

	"github.com/jansemmelink/asn1/def2"
	"github.com/jansemmelink/asn1/parser"
	"github.com/jansemmelink/mem"
)

func TestChoice(t *testing.T) {
	tests := []struct {text string; sel int; u uint; i int; n string; } {
		{"1234", 1, 1234, 0, ""},
		{"-54321", 2, 0, -54321, ""},
		{"Jan", 3, 0, 0, "Jan"},
		{"-InvalidID", -1, 0, 0, ""},
	}
	for i,tt := range tests {
		//write the textual test data:
		l := parser.NewLines()
		l = l.Append(i+1, tt.text, "test data")

		remain := l
		//x := MyValue{}
		x := mem.NewX(reflect.TypeOf(MyValue{})).(*MyValue)
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
				if x.Choice.Option().Index != tt.sel {
					t.Fatalf("Parsed selector=%d != tt.sel=%d", x.Choice.Option().Index, tt.sel)
				}
				switch tt.sel {
				case 1:
					parsedUint := uint(x.U)
					if parsedUint != tt.u || tt.u == 0 {
						t.Fatalf("test[%d]: parsed uint %d expected %d (must not be 0)", i, parsedUint, tt.u)
					}
				case 2:
					parsedInt := int(x.I)
					if parsedInt != tt.i || tt.i == 0 {
						t.Fatalf("test[%d]: parsed int %d expected %d (must not be 0)", i, parsedInt, tt.i)
					}
				case 3:
					parsedStr := string(x.N)
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
}

//MyValue parses a uint or int or string:
type MyValue struct {
	def2.Choice
	U def2.Uint
	I def2.Int
	N Identifier
}
