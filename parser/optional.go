package parser

import (
	"bitbucket.org/vservices/dark/logger"
)

//Optional creates a parser that detects presents of an item
func Optional(item INotation) IOptional {
	name := item.Name()
	if len(name) == 0 {
		panic("Optional without name")
	}
	o := &optional{
		INotation: New("optional", name),
	}
	o.expectedItem = item
	return o
}

//IOptional extends INotation
type IOptional interface {
	INotation
}

//optional implements IOptional
type optional struct {
	INotation
	expectedItem INotation
}

func (o optional) ParseV(log logger.ILogger, l ILines) (IValue, ILines, error) {
	log.Debugf("%T(%s).Parse() from line %d: %.32s ...", o, o.Name(), l.LineNr(), l.Next())
	var err error
	var remain ILines
	var parsedValue IValue
	parsedValue, remain, err = o.expectedItem.ParseV(log.With("opt(%s)", o.Name()), l)
	if err == nil {
		log.Debugf("optional(%s) parsed on line %d, next is line %d: %.32s", o.expectedItem.Name(), l.LineNr(), remain.LineNr(), remain.Next())
		return parsedValue, remain, nil
	}

	//absent: not an error!
	log.Debugf("optional(%s) is absent on line %d: %.32s", o.expectedItem.Name(), l.LineNr(), l.Next())
	return nil, l, nil
} //optional.Parse()
