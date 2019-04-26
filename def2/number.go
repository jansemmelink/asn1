package def2

import (
	"strconv"
	"strings"

	"bitbucket.org/vservices/dark/logger"
	"github.com/jansemmelink/asn1/parser"
)

var tempNextIntValue = 4000

func nrTerm(c rune) bool {
	//accept any numeric chars, even if not for real or signed values
	if c >= '0' && c <= '9' || c == '-' || c == 'e' || c == 'E' || c == '.' {
		return false
	}
	return true
}

//Int is parsable signed integer number (golang type int)
type Int int

//Parse ...
func (i *Int) Parse(log logger.ILogger, l parser.ILines, v IParsable) (parser.ILines, error) {
	next := l.Next()
	posEnd := strings.IndexFunc(next, nrTerm)
	if posEnd < 0 || posEnd > 32 {
		return l, log.Wrapf(nil, "not a valid number on line %d: %.32s", l.LineNr(), next)
	}
	str := next[0:posEnd]
	if len(str) == 0 {
		return l, log.Wrapf(nil, "No Int value at line %d: %.32s ...", l.LineNr(), next)
	}
	//log.Debugf("Parsing int from \"%s\"", str)
	//spacing not allowed between digits:
	// str = strings.Replace(str, " ", "", -1)
	// str = strings.Replace(str, "\t", "", -1)
	value, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return l, log.Wrapf(err, "\"%s\" is invalid int", str)
	}

	var out = v.(*Int)
	*out = Int(value)
	remain, _ := l.SkipOver(str)
	return remain, nil
} //Int.Parse()

//Uint is parsable unsigned integer number (golang type uint)
type Uint uint

//Parse ...
func (i *Uint) Parse(log logger.ILogger, l parser.ILines, v IParsable) (parser.ILines, error) {
	next := l.Next()
	posEnd := strings.IndexFunc(next, nrTerm)
	if posEnd < 0 || posEnd > 32 {
		return l, log.Wrapf(nil, "not a valid number on line %d: %.32s", l.LineNr(), next)
	}
	str := next[0 : posEnd+1]
	if len(str) == 0 {
		return l, log.Wrapf(nil, "No Uint value at line %d: %.32s ...", l.LineNr(), next)
	}

	log.Debugf("Parsing uint from \"%s\"", str)
	//spacing not allowed between digits:
	// str = strings.Replace(str, " ", "", -1)
	// str = strings.Replace(str, "\t", "", -1)
	value, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return l, log.Wrapf(err, "\"%s\" is invalid uint", str)
	}

	var out = v.(*Uint)
	*out = Uint(value)
	remain, _ := l.SkipOver(str)
	return remain, nil
} //Uint.Parse()
