package parser

import "fmt"

//Optional creates a parser that detects presents of an item
func Optional(item INotation) IOptional {
	return optional{
		INotation:    New(fmt.Sprintf("Optional %s", item.Name())),
		expectedItem: item,
		item:         nil,
	}
}

//IOptional extends INotation
type IOptional interface {
	INotation

	//Item returns the parsed item if present or nil if absent
	Item() INotation
}

//optional implements IOptional
type optional struct {
	INotation
	expectedItem INotation
	item         INotation
}

func (o optional) Item() INotation {
	return o.item
}

func (o optional) Parse(l ILines) (INotation, ILines, error) {
	var err error
	var remain ILines
	o.item, remain, err = o.expectedItem.Parse(l)
	if err == nil {
		log.Debugf("optional(%s) parsed on line %d, next is line %d: %.32s", o.expectedItem.Name(), l.LineNr(), remain.LineNr(), remain.Next())
		return o, remain, nil
	}

	//absent: not an error!
	log.Debugf("optional(%s) is absent on line %d: %.32s", o.expectedItem.Name(), l.LineNr(), l.Next())
	o.item = nil
	return o, l, nil
} //optional.Parse()
