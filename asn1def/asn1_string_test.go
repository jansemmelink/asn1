package asn1def_test

import (
	"testing"

	"bitbucket.org/vservices/dark/logger"
	"bitbucket.org/vservices/dark/logger/level"
	"github.com/jansemmelink/asn1/asn1def"
	"github.com/jansemmelink/asn1/parser"
)

var log = logger.New().WithLevel(level.Info)

func TestCString(t *testing.T) {
	l := parser.NewLines()
	l = l.Append(1, "\"ABC   DEF  ", "Trailing spaces must not be included")
	l = l.Append(2, "\tMore ", "Trim both lead and trail")
	l = l.Append(3, "    GHI\"\"JKL\" and some other text", "Leading spaces must not be included")
	l = l.Append(4, "END", "Last line")

	x := asn1def.CString()
	v, remain, err := x.ParseV(log, l)
	if err != nil {
		t.Fatalf("Failed to parse cstring: %v", err)
	}
	if remain.LineNr() != 2 && remain.Next() != "and some other text" {
		t.Fatalf("Remain line %d: %.32s (expected and some other text)", remain.LineNr(), remain.Next())
	}

	cs := v.String()
	if cs != "ABC   DEFMoreGHI\"JKL" {
		t.Fatalf("cstring not (%s) but (%s)", "ABC   DEFMoreGHI\"JKL", cs)
	}
} //TestCString()

func TestBString(t *testing.T) {
	l := parser.NewLines()
	l = l.Append(1, "'00000000'B some text", "Basic value")
	l = l.Append(2, "   '10000000'B some text", "Basic value")
	l = l.Append(3, "\t'0001\t0010\t0011 0100 0101 0110 0111 1000'B some text", "Long binary with tabs and spaces for clarity are allowed")

	x := asn1def.BString()
	l = testBString(t, log, l, x, 0)
	l = testSkip(t, log, l, "some")
	l = testSkip(t, log, l, "text")
	l = testBString(t, log, l, x, 0x80)
	l = testSkip(t, log, l, "some")
	l = testSkip(t, log, l, "text")
	l = testBString(t, log, l, x, 0x12345678)
	l = testSkip(t, log, l, "some")
	l = testSkip(t, log, l, "text")
} //TestBString()

func TestHString(t *testing.T) {
	l := parser.NewLines()
	l = l.Append(1, "'00000000'H some text", "Basic value")
	l = l.Append(2, "   '80'H some text", "Basic value")
	l = l.Append(3, "\t'12\t34 56\t78'H some text", "Long hex with tabs and spaces for clarity are allowed")

	x := asn1def.HString()
	l = testHString(t, log, l, x, 0)
	l = testSkip(t, log, l, "some")
	l = testSkip(t, log, l, "text")
	l = testHString(t, log, l, x, 0x80)
	l = testSkip(t, log, l, "some")
	l = testSkip(t, log, l, "text")
	l = testHString(t, log, l, x, 0x12345678)
	l = testSkip(t, log, l, "some")
	l = testSkip(t, log, l, "text")
} //TestHString()

//skip over some text to get to next piece to be tested, but fail if not present
func testSkip(t *testing.T, log logger.ILogger, l parser.ILines, textToSkip string) parser.ILines {
	remain, ok := l.SkipOver(textToSkip)
	if !ok {
		t.Fatalf("expected \"%s\", but not found at line %d: %.32s ...", textToSkip, l.LineNr(), l.Next())
	}
	return remain
} //testSkip()

//check that specified binary string value is next
func testBString(t *testing.T, log logger.ILogger, l parser.ILines, x asn1def.IBString, expectedValue uint64) parser.ILines {
	v, remain, err := x.ParseV(log, l)
	if err != nil {
		t.Fatalf("Failed to parse bstring: %v", err)
	}
	if remain.LineNr() != 1 && remain.Next() != "some text" {
		t.Fatalf("Remain line %d: %.32s (some text)", remain.LineNr(), remain.Next())
	}

	bs, ok := v.Value().(uint64)
	if !ok {
		t.Fatalf("v is %T, not uint64", v)
	}
	if bs != expectedValue {
		t.Fatalf("bstring not (0x%08x) but (0x%08x)", expectedValue, bs)
	}
	log.Infof("GOOD: Matching bstring=0x%08x, next is line %d: %.32s ...",
		expectedValue,
		remain.LineNr(),
		remain.Next())
	return remain
} //testBString()

//check that specified hex string value is next
func testHString(t *testing.T, log logger.ILogger, l parser.ILines, x asn1def.IHString, expectedValue uint64) parser.ILines {
	v, remain, err := x.ParseV(log, l)
	if err != nil {
		t.Fatalf("Failed to parse hstring: %v", err)
	}
	if remain.LineNr() != 1 && remain.Next() != "some text" {
		t.Fatalf("Remain line %d: %.32s (some text)", remain.LineNr(), remain.Next())
	}

	hs, ok := v.Value().(uint64)
	if !ok {
		t.Fatalf("v is %T, not uint64", v)
	}
	if hs != expectedValue {
		t.Fatalf("hstring not (0x%08x) but (0x%08x)", expectedValue, hs)
	}
	log.Infof("GOOD: Matching hstring=0x%08x, next is line %d: %.32s ...",
		expectedValue,
		remain.LineNr(),
		remain.Next())
	return remain
} //testHString()
