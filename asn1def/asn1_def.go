package asn1def

import (
	"bitbucket.org/vservices/dark/logger"
	"github.com/jansemmelink/asn1/parser"
)

var (
	log = logger.New()
)

//New empty definition
func New() IDefinition {
	return definition{}
}

//IDefinition of ASN.1 data structures
type IDefinition interface {
	LoadFile(filename string) error
	//Unmarshal() error
}

type definition struct {
}

func (d definition) LoadFile(filename string) error {
	log.Debugf("Loading %s ...", filename)

	lines, err := parser.LinesFromFile(filename)
	if err != nil {
		return log.Wrapf(err, "Failed to read lines from file %s", filename)
	}

	log.Debugf("Read %d significant lines from %s:", lines.Count(), filename)

	//silence logger in parser package:
	//logger.Top().New("github.com/jansemmelink/asn1/parser").WithLevel(level.Error)

	//loop several times over remaining lines
	//as one file may contain several asn1def's
	modules := make([]parser.IValue, 0)
	for lines.Count() > 0 {
		moduleDefinition := ModuleDefinition()
		parsedModuleDefinition, remainingLines, err := moduleDefinition.ParseV(log.With("asn1def"), lines)
		if err != nil {
			return log.Wrapf(err, "Failed to parse %s from file %s line %d", moduleDefinition.Name(), filename, lines.LineNr())
		}

		log.Debugf("Parsed ASN.1 Definition: %s", parsedModuleDefinition.Name())
		modules = append(modules, parsedModuleDefinition)
		if remainingLines.LineNr() == lines.LineNr() {
			return log.Wrapf(nil, "After  parsing, same lines remain... nothing used!?!?!")
		}

		lines = remainingLines
	} //for whole file

	for _, n := range modules {
		log.Debugf("module(%s): %+v", n.Name(), n.Value())
	}

	return nil
} //definition.LoadFile()

type handlerDefinition struct {
	defs []string
}

func (h handlerDefinition) Handle(n parser.INotation) error {
	// log.Debugf("HANDLING %T", n)
	// s := n.(parser.ISeq)
	// //expect {identifier, specification}
	// parsedItems := s.Items()
	// if len(parsedItems) != 2 {
	// 	return log.Wrapf(nil, "definition seq has %d != 2 items", len(parsedItems))
	// }

	// id := parsedItems[0].(parser.IRegex).Match()
	// spec := parsedItems[1].(parser.IChoice).Item()
	// switch spec.Name() {
	// case "specType": //::= <type specification>
	// 	return log.Wrapf(nil, "Unknown specification type %s", spec.Name())
	// case "objectId": //OBJECT IDENTIFIER ::= { ... }
	// 	h.Add(id)
	// case "typedValue": //<typeName> ::= <value>
	// 	return log.Wrapf(nil, "Unknown specification type %s", spec.Name())
	// default:
	// 	return log.Wrapf(nil, "Unknown specification type %s", spec.Name())
	// }
	return nil
}

func (h handlerDefinition) Add(n string) {
	//h.defs = append(h.defs, n)
}
