package def2_test

import (
	"fmt"
	"math/rand"
	"testing"
	"reflect"

	"github.com/jansemmelink/asn1/def2"
	"github.com/jansemmelink/mem"
	"github.com/jansemmelink/asn1/parser"
)

func TestList(t *testing.T) {
	l := parser.NewLines()

	//append CSV with uint nrs:
	nrTests := 5
	nr := make([]uint, nrTests)
	s := ""
	for i := 0; i < nrTests; i++ {
		nr[i] = uint(rand.Uint64())
		if i > 0 {
			s += ","
		}
		s += fmt.Sprintf("%d", nr[i])
	}
	l = l.Append(1, s, "no comment")

	remain := l
	u := mem.NewV(reflect.TypeOf(UintCSV{})).(UintCSV)
	var err error
	_, err = u.Parse(log, remain, &u)
	if err != nil || len(u.Items()) != nrTests {
		t.Fatalf("Failed: count=%d: err=%v", len(u.Items()), err)
	}

	log.Debugf("Success: Parsed %d items", len(u.Items()))
	for i := 0; i < nrTests; i++ {
		item := u.Items()[i].(*def2.Uint)
		log.Debugf("  [%d]: %T = %+v", i, item, *item)
		if uint(*item) != nr[i] {
			t.Fatalf("Parsed[%d]=%d != expected %d", i, *item, nr[i])
		}
	} //for each test nr
	return
}

//List parses CSV of Uints
type UintCSV struct{ *def2.List }

func (b UintCSV) Sep() def2.IToken { return def2.NewToken(",") }
func (b UintCSV) Min() int         { return 0 }
func (b UintCSV) Max() int         { return 10 }
func (b UintCSV) ItemType() reflect.Type {
	u := def2.Uint(1)
	return reflect.TypeOf(u)
}
