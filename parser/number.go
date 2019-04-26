package parser

import (
	"strconv"
	"unicode"

	"bitbucket.org/vservices/dark/logger"
)

//INumber is a positive integer number
type INumber interface {
	INotation
}

var n INumber

//Number creates a parser for a decimal integer number
func Number() INumber {
	if n == nil {
		n = &number{
			INotation: New("number", "number"),
		}
	}
	return n
}

type number struct {
	INotation
}

func (n number) ParseV(log logger.ILogger, l ILines) (IValue, ILines, error) {
	//get digits
	next := l.Next()
	numberString := ""
	for _, ch := range next {
		if !unicode.IsDigit(ch) {
			break
		}
		numberString += string(ch)
	}

	var err error
	intValue := 0
	if intValue, err = strconv.Atoi(numberString); err == nil {
		remain, ok := l.SkipOver(numberString)
		if !ok {
			return nil, l, log.Wrapf(nil, "Failed to skip number \"%s\" after parsing on line %d: %.32s ...", numberString, l.LineNr(), l.Next())
		}
		return IntValue(n.Name(), intValue), remain, nil
	}
	return nil, l, log.Wrapf(nil, "number not valid on line %d: %.32s ...", l.LineNr(), next)
} //number.ParseV()

//ISignedNumber ...
type ISignedNumber interface {
	INotation
}

var sn ISignedNumber

//SignedNumber ...
func SignedNumber() ISignedNumber {
	if sn == nil {
		sn = &signedNumber{
			INotation: New("signedNumber", "signedNumber"),
		}
	}
	return sn
}

type signedNumber struct {
	INotation
}

func (n signedNumber) ParseV(log logger.ILogger, l ILines) (IValue, ILines, error) {
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
	intValue := 0
	if intValue, err = strconv.Atoi(numberString); err == nil {
		remain, _ := l.SkipOver(numberString)
		return IntValue(sn.Name(), intValue), remain, nil
	}
	return nil, l, log.Wrapf(nil, "signedNumber not valid on line %d: %.32s", l.LineNr(), next)
} //signedNumber.Parse()
