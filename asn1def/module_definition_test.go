package asn1def_test

import (
	"testing"

//	"bitbucket.org/vservices/dark/logger"
"bitbucket.org/vservices/dark/logger/level"
	"github.com/jansemmelink/asn1/asn1def"
	"github.com/jansemmelink/asn1/parser"
)

//var log = logger.New()

func TestModuleReference(t *testing.T) {
	l := parser.NewLines()
	expectedValue := "MAP-DialogueInformation"
	l = l.Append(l.LineNr()+1, expectedValue+" { one two (2) three (3) }", "")

	x := asn1def.ModuleReference()
	v, remain, err := x.ParseV(log, l)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}
	if v.Value().(string) != expectedValue {
		t.Fatalf("Parsed \"%v\", but expected \"%s\"", v.Value(), expectedValue)
	}
	log.Infof("GOOD: ModuleReference = %+v", v.Value())
	log.Debugf("Remain line %d: %.32s", remain.LineNr(), remain.Next())
}

func TestModuleIdentifier(t *testing.T) {
	l := parser.NewLines()
	expectedValue := "MAP-DialogueInformation"
	l = l.Append(l.LineNr()+1, expectedValue+" { one two (2) three (3) }", "")

	log = log.WithLevel(level.Debug)



	
	x := asn1def.ModuleIdentifier
	v, remain, err := x.ParseV(log, l)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}
	//parsedValue := v.(asn1def.IModuleIdentifier).Reference().Value()
	// log.Debugf("Success: %+v", parsedValue)
	// if parsedValue != expectedValue {
	// 	t.Fatalf("Parsed \"%s\", but expected \"%s\"", parsedValue, expectedValue)
	// }
	log.Infof("GOOD: ModuleIdentifier = %+v", v.Value())
	log.Debugf("Remain line %d: %.32s", remain.LineNr(), remain.Next())
}

func Test1(t *testing.T) {
	l := parser.NewLines()
	l = l.Append(l.LineNr()+1, "MAP-DialogueInformation { one two (2) three (3) }", "")
	l = l.Append(l.LineNr()+1, "DEFINITIONS", "")
	l = l.Append(l.LineNr()+1, "IMPLICIT TAGS", "")
	l = l.Append(l.LineNr()+1, "::=", "")
	l = l.Append(l.LineNr()+1, "BEGIN", "")
	// EXPORTS oneOid, TwoChoice;
	// IMPORTS extNetworkId, extAsId FROM MyModuleReference { one (1) two (2) };
	// oneOid OBJECT IDENTIFIER ::= {extNetworkId extAsId someNr (2) ver (4)}
	// TwoChoice ::= CHOICE { map-open [0] MAP-OpenInfo, map-accept [1] MAP-AcceptInfo }
	l = l.Append(l.LineNr()+1, "END", "")

	x := asn1def.ModuleDefinition()
	v, remain, err := x.ParseV(log, l)
	if err != nil {
		t.Fatalf("Failed to parse module definition")
	}

	md := v.(parser.IValueV)
	log.Debugf("Success: %+v", md.Value())
	log.Debugf("Remain line %d: %.32s", remain.LineNr(), remain.Next())
}
