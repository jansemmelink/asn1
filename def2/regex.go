package def2

import (
	"reflect"
	"regexp"
	"strings"

	"bitbucket.org/vservices/dark/logger"
	"github.com/jansemmelink/asn1/parser"
)

//NewRegex makes a parsable regular expression
func NewRegex(pattern string) IRegex {
	compiledRegex, err := regexp.Compile(pattern)
	if err != nil {
		panic(log.Wrapf(err, "Regex(\"%s\") compilation failed", pattern))
	}
	log.Debugf("Compiled Regex(\"%s\")", pattern)
	return Regex{
		patternString: pattern,
		compiled:      compiledRegex,
		parsedValue:   "",
	}
}

//IRegex ...
type IRegex interface {
	IParsable
	Pattern() string
	Compiled() *regexp.Regexp
	SetValue(s string)
	Value() string
}

//Regex implements IRegex
type Regex struct {
	//todo: reuse for multiple instances...
	patternString string
	compiled      *regexp.Regexp

	//set in parser:
	parsedValue string
}

//Pattern ...
func (r Regex) Pattern() string {
	return r.patternString
}

//Compiled  ...
func (r Regex) Compiled() *regexp.Regexp {
	return r.compiled
}

//Value ...
func (r Regex) Value() string {
	return r.parsedValue
}

//SetValue ...
func (r Regex) SetValue(s string) {
	log.Debugf("VALUE = %s", s)
	r.parsedValue = s
}

//Parse ...
func (Regex) Parse(log logger.ILogger, l parser.ILines, v IParsable) (parser.ILines, error) {
	log.Debugf("REGEX %T", v)
	vv := reflect.ValueOf(v).Elem().Interface()
	log.Debugf("vv=%T", vv)
	r := vv.(IRegex)
	compiled := r.Compiled()
	if compiled == nil {
		var err error
		compiled, err = regexp.Compile(r.Pattern())
		if err != nil {
			return l, log.Wrapf(err, "%T(\"%s\") compilation failed", r.Pattern())
		}
		log.Debugf("%T(\"%s\").Compiled", vv, r.Pattern())
	}

	// log.Debugf("REGEX %T(%s)", out, out.Pattern())
	next := l.Next()
	matches := compiled.FindAllString(next, 1)
	if len(matches) >= 1 {
		//matched, but could be later in the line and could be multiple
		//we only consider the first match
		//in which case SkipOver() will fail
		if remain, ok := l.SkipOver(matches[0]); ok {
			r.SetValue(strings.TrimRight(matches[0], " \t"))

			//out.parsedValue = strings.TrimRight(matches[0], " \t")
			log.Debugf("regex(%s) parsed \"%s\" on line %d, next line %d: %.32s ...", r.Pattern(), r.Value(), l.LineNr(), remain.LineNr(), remain.Next())
			return remain, nil
		}
	}
	return l, log.Wrapf(nil, "regex(%s) mismatch line %d: %.32s", r.Pattern(), l.LineNr(), next)
} //Regex.Parse()
