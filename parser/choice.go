package parser

import (
	"bitbucket.org/vservices/dark/logger"
)

//Choice creates a parser that succeeds on any of the specified notations
func Choice(name string, options ...INotation) IChoice {
	c := &choice{
		INotation: New("choice", name),
		options:   make([]INotation, 0),
	}
	for _, o := range options {
		c.options = append(c.options, o)
	}
	return c
}

//IChoice extends INotation
type IChoice interface {
	INotation
	Add(INotation) IChoice
}

//choice implements IChoice
type choice struct {
	INotation
	options []INotation
}

func (c *choice) Add(item INotation) IChoice {
	// choicesMutex.Lock()
	// defer choicesMutex.Unlock()
	c.options = append(c.options, item)
	//	choices[c.Name()] = c
	return c
}

func (c choice) ParseV(log logger.ILogger, l ILines) (IValue, ILines, error) {
	for _, o := range c.options {
		var err error
		remain := l
		var parsedValue IValue
		if parsedValue, remain, err = o.ParseV(log, remain); err == nil {
			return parsedValue, remain, nil
		}
	}
	return nil, l, log.Wrapf(nil, "failed")
} //choice.Parse()
