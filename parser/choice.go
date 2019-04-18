package parser

//Choice creates a parser that succeeds on any of the specified notations
func Choice(name string, options ...INotation) IChoice {
	return choice{
		INotation: New(name),
		options:   options,
	}
}

//IChoice extends INotation
type IChoice interface {
	INotation
	ItemType() INotation
	Item() INotation
}

//choice implements IChoice
type choice struct {
	INotation
	//specification
	options []INotation

	//output
	itemType INotation
	item     INotation
}

func (c choice) ItemType() INotation {
	return c.itemType
}

func (c choice) Item() INotation {
	return c.item
}

func (c choice) Parse(l ILines) (INotation, ILines, error) {
	//log.Debugf("choice(%s): Start from line %d: %.32s", c.Name(), l.LineNr(), l.Next())
	for i, o := range c.options {
		var err error
		remain := l
		if c.item, remain, err = o.Parse(remain); err == nil {
			log.Debugf("choice(%s).[%d]=%s parsed on line %d, next is line %d: %.32s", c.Name(), i, o.Name(), l.LineNr(), remain.LineNr(), remain.Next())
			return c, remain, nil
		}
	}
	return nil, l, log.Wrapf(nil, "choice(%s) failed on line %d: %.32s", c.Name(), l.LineNr(), l.Next())
} //choice.Parse()
