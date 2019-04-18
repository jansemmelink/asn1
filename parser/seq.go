package parser

import (
	"fmt"
)

//Seq creates a sequence parser
func Seq(name string, items ...INotation) ISeq {
	if len(items) == 0 {
		panic("empty sequence")
	}

	logDesc := ""
	for _, i := range items {
		logDesc += fmt.Sprintf(" %T", i)
	}

	return seq{
		INotation:     New(name),
		expectedItems: items,
		logDesc:       logDesc[1:],
	}
}

//ISeq of items
type ISeq interface {
	INotation
	Items() []INotation
}

type seq struct {
	INotation
	expectedItems []INotation
	logDesc       string

	items []INotation
}

func (s seq) Parse(l ILines) (INotation, ILines, error) {
	s.items = make([]INotation, 0)
	remain := l
	//log.Debugf("seq(%s) %d items from all=%.32s", s.Name(), len(s.expectedItems), l.Next())
	for i, e := range s.expectedItems {
		//log.Debugf("seq(%s)[%d].%s from r=%.32s all=%.32s", s.Name(), i, e.Name(), remain.Next(), l.Next())
		var err error
		var parsedItem INotation
		parsedItem, remain, err = e.Parse(remain)
		if err != nil {
			return nil, l, log.Wrapf(err, "seq(%s)[%d].%s failed on line %d: %s", s.Name(), i, e.Name(), l.LineNr(), l.Next())
		}
		s.items = append(s.items, parsedItem)
		log.Debugf("  seq(%s)[%d]=%s parsed, next is line %d: %.32s", s.Name(), i, e.Name(), remain.LineNr(), remain.Next())
	}
	log.Debugf("seq(%s)={%s} parsed on line %d, next is line %d: %.32s", s.Name(), s.logDesc, l.LineNr(), remain.LineNr(), remain.Next())
	return s, remain, nil
} //seq.Parse()

func (s seq) Items() []INotation {
	return s.items
}
