package parser_test

import (
	"testing"

	"bitbucket.org/vservices/dark/logger"
	"bitbucket.org/vservices/dark/logger/level"
	"github.com/jansemmelink/asn1/parser"
)

var log = logger.New().WithLevel(level.Debug)

func TestNumber(t *testing.T) {
	l := parser.NewLines()
	l = l.Append(1, "\t 0 some text", "Zero value")
	l = l.Append(2, " 12 some ", "other value")
	l = l.Append(3, " \t text 12345678 some text", "Bigger value")
	l = l.Append(4, "-12 must fail", "Negative not allowed in parser.Number")

	x := parser.Number()
	l = testNumber(t, log, l, x, 0, false)
	l = testSkip(t, log, l, "some")
	l = testSkip(t, log, l, "text")
	l = testNumber(t, log, l, x, 12, false)
	l = testSkip(t, log, l, "some")
	l = testSkip(t, log, l, "text")
	l = testNumber(t, log, l, x, 12345678, false)
	l = testSkip(t, log, l, "some")
	l = testSkip(t, log, l, "text")
	l = testNumber(t, log, l, x, -12, true)

	y := parser.SignedNumber()
	l = testNumber(t, log, l, y, -12, false)
	l = testSkip(t, log, l, "must")
	l = testSkip(t, log, l, "fail")
} //TestCString()

//skip over some text to get to next piece to be tested, but fail if not present
func testSkip(t *testing.T, log logger.ILogger, l parser.ILines, textToSkip string) parser.ILines {
	remain, ok := l.SkipOver(textToSkip)
	if !ok {
		t.Fatalf("expected \"%s\", but not found at line %d: %.32s ...", textToSkip, l.LineNr(), l.Next())
	}
	return remain
} //testSkip()

//check that specified number value is next
func testNumber(t *testing.T, log logger.ILogger, l parser.ILines, x parser.INumber, expectedValue int, mustFail bool) parser.ILines {
	v, remain, err := x.ParseV(log, l)
	if mustFail {
		//expect to fail
		if err == nil {
			t.Fatalf("Parser successed on invalid number on line %d: %.32s", l.LineNr(), l.Next())
		}
		return remain
	}

	//expect to succeed
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}
	vv, ok := v.(parser.IIntValue)
	if !ok {
		t.Fatalf("v is %T, not %T", v, x)
	}
	if vv.Int() != expectedValue {
		t.Fatalf("%T = %d, not expected %d", vv, vv.Value(), expectedValue)
	}
	log.Debugf("GOOD: Matching %T=%d, next is line %d: %.32s ...",
		vv,
		vv.Value(),
		remain.LineNr(),
		remain.Next())
	return remain
} //testNumber()
