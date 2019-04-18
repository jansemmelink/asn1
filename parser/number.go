package parser

import (
	"strconv"
	"unicode"
)

//Number creates a parser for a decimal integer number
func Number() INumber {
	return number{
		INotation: New("number"),
		value:     0,
	}
}

//INumber ...
type INumber interface {
	INotation
	Value() int
}

type number struct {
	INotation
	value int
}

func (n number) Value() int {
	return n.value
}

func (n number) Parse(l ILines) (INotation, ILines, error) {
	//get digits
	next := l.Next()
	numberString := ""
	for pos, ch := range next {
		if !unicode.IsDigit(ch) {
			if pos > 0 || string(ch) != "-" {
				break
			}
		}
		numberString += string(ch)
	}

	var err error
	if n.value, err = strconv.Atoi(numberString); err == nil {
		remain, _ := l.SkipOver(numberString)
		//log.Debugf("number=%d on line %d", n.value, l.LineNr())
		return n, remain, nil
	}
	return nil, l, log.Wrapf(nil, "number not valid on line %d: %.32s", l.LineNr(), next)
} //number.Parse()
