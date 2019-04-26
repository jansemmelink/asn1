package parser_test

import (
	"testing"

	"github.com/jansemmelink/asn1/parser"
)

func TestLines1(t *testing.T) {
	l := parser.NewLines()
	l = l.Append(1, "abc", "efg")
	if l.Count() != 1 {
		t.Fatalf("not 1")
	}
	if l.Next() != "abc" {
		t.Fatalf("not abc")
	}
	l, _ = l.SkipOver("a")
	if l.Next() != "bc" {
		t.Fatalf("not bc: %s", l.Next())
	}
	l, _ = l.SkipOver("bc")
	if l.Next() != "" {
		t.Fatalf("not \"\": %s", l.Next())
	}
	if l.Count() > 0 {
		t.Fatalf("not 0 lines")
	}
	return
}
