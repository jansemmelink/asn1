package parser

//Block creates a parser for something between a start and end terminator
func Block(name string, start, end string, inside INotation) IBlock {
	return block{
		INotation:      New(name),
		start:          start,
		end:            end,
		expectedInside: inside,
	}
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

func (b block) Parse(l ILines) (INotation, ILines, error) {
	remain, ok := l.SkipOver(b.start)
	if !ok {
		return nil, l, log.Wrapf(nil, "Line %d: block %s...%s start token not present in line %d: %.32s", l.LineNr(), b.start, b.end, l.LineNr(), l.Next())
	}
	var err error
	if b.item, remain, err = b.expectedInside.Parse(remain); err != nil {
		return nil, l, log.Wrapf(err, "Line %d: block %s..%s content not valid", l.LineNr(), b.start, b.end)
	}
	remain, ok = remain.SkipOver(b.end)
	if !ok {
		return nil, l, log.Wrapf(nil, "block(%s) %s...%s started on line %d: end token expected on line %d: %.32s", b.Name(), b.start, b.end, l.LineNr(), remain.LineNr(), remain.Next())
	}
	log.Debugf("block(%s) %s <%T> %s parsed on line %d, next line %d: %.32s", b.Name(), b.start, b.expectedInside, b.end, l.LineNr(), remain.LineNr(), remain.Next())
	return b, remain, nil
} //block.Parse()
