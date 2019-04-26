package parser

import (
	"bitbucket.org/vservices/dark/logger"
)

//Seq creates a sequence of consecutive items to parse
func Seq(name string, items ...INotation) ISeq {
	if len(items) == 0 {
		panic("empty sequence")
	}
	s := &seq{
		INotation:     New("seq", name),
		expectedItems: make([]INotation, 0),
	}
	for index, i := range items {
		s.expectedItems = append(s.expectedItems, i)
		log.Debugf("%p(%s)[%d]=%p(%s)", s, s.Name(), index, i, i.Name())
	}
	return s
}

//ISeq of items
type ISeq interface {
	INotation
}

type seq struct {
	INotation
	expectedItems []INotation
}

func (s seq) ParseV(log logger.ILogger, l ILines) (IValue, ILines, error) {
	log = log.With("seq(%s)", s.Name())
	log.Debugf("line %d: %.32s ...", l.LineNr(), l.Next())

	parsedSeq := ListValue(s.Name())
	remain := l
	for i, e := range s.expectedItems {
		var err error
		var parsedValue IValue
		parsedValue, remain, err = e.ParseV(log.With("[%d]", i), remain)
		if err != nil {
			return nil, l, log.Wrapf(err, "seq(%s)[%d].%s failed on line %d: %s", s.Name(), i, e.Name(), l.LineNr(), l.Next())
		}
		parsedSeq = parsedSeq.With(parsedValue)
		log.Debugf("  seq(%s)[%d]=%s parsed, next is line %d: %.32s", s.Name(), i, e.Name(), remain.LineNr(), remain.Next())
	}
	log.Debugf("seq(%s) parsed on line %d, next is line %d: %.32s", s.Name(), l.LineNr(), remain.LineNr(), remain.Next())
	return parsedSeq, remain, nil
} //seq.ParseV()

