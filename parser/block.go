package parser

import "bitbucket.org/vservices/dark/logger"

//Block creates a parser for something between a start and end terminator
func Block(name string, start, end string, inside INotation) IBlock {
	b := &block{
		INotation: New("block", name),
		start:     start,
		end:       end,
	}
	b.expectedInside = inside
	return b
}

//IBlock enclodes between a start and end terminator, e.g. "{...}" or "begin...end"
type IBlock interface {
	INotation
	Item() INotation
}

type block struct {
	INotation
	start          string
	end            string
	expectedInside INotation

	item INotation
}

func (b block) Item() INotation {
	return b.item
}

func (b block) ParseV(log logger.ILogger, l ILines) (IValue, ILines, error) {
	remain, ok := l.SkipOver(b.start)
	if !ok {
		return nil, l, log.Wrapf(nil, "Line %d: block %s...%s start token not present in line %d: %.32s", l.LineNr(), b.start, b.end, l.LineNr(), l.Next())
	}
	var err error
	var parsedValue IValue
	if parsedValue, remain, err = b.expectedInside.ParseV(log.With("block(%s)", b.Name()), remain); err != nil {
		return nil, l, log.Wrapf(err, "Line %d: block %s..%s content not valid", l.LineNr(), b.start, b.end)
	}
	remain, ok = remain.SkipOver(b.end)
	if !ok {
		return nil, l, log.Wrapf(nil, "block(%s) %s...%s started on line %d: end token expected on line %d: %.32s", b.Name(), b.start, b.end, l.LineNr(), remain.LineNr(), remain.Next())
	}
	return parsedValue, remain, nil
} //block.ParseV()
