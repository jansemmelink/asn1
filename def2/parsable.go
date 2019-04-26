package def2

import (
	"reflect"

	"bitbucket.org/vservices/dark/logger"
	"github.com/jansemmelink/asn1/parser"
)

//IParsable can be parsed
type IParsable interface {
	Parse(log logger.ILogger, l parser.ILines, v IParsable) (parser.ILines, error)
}

//==============================================================
//to get an interface type in a variable that one
//can use in SomeOtherType.Implements(interfaceType), use this pattern:
//==============================================================
//	interfaceType = reflect.TypeOf((*IMyInterface)(nil)).Elem(IMyInterface)
//==============================================================
//var ParsableInterfaceType = reflect.TypeOf((*IParsable)(nil)).Elem()

//Parse into a parsable variable
func Parse(log logger.ILogger, l parser.ILines, v IParsable) (parser.ILines, error) {
	if reflect.TypeOf(v).Kind() != reflect.Ptr {
		return l, log.Wrapf(nil, "Parse(log,lines,%T) param 3 is not pointer", v)
	}

	// valueType := reflect.TypeOf(v).Elem()
	// if !valueType.Implements(ParsableInterfaceType) {
	// 	return l, log.Wrapf(nil, "Parse(log,lines,%T) param 3 is not parsable", v)
	// }

	remain, err := v.Parse(log, l, v)
	if err != nil {
		return l, log.Wrapf(err, "Failed to parse %s from line %d: %.32s ...", reflect.TypeOf(v).Elem().Name(), l.LineNr(), l.Next())
	}
	return remain, nil
} //Parse()
