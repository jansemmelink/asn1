package parser

import (
	"regexp"
	"strings"
)

//Regex creates a regex parser
func Regex(name, pattern string) IRegex {
	compiled, err := regexp.Compile(pattern)
	if err != nil {
		panic(log.Wrapf(err, "Invalid regular expression pattern: %s", pattern))
	}
	return regex{
		INotation:     New(name),
		patternString: pattern,
		pattern:       compiled,
	}

}

//IRegex matches a regular expression
type IRegex interface {
	INotation
	Match() string
}

type regex struct {
	INotation
	patternString string
	pattern       *regexp.Regexp
	match         string
}

func (r regex) Match() string {
	return r.match
}

func (r regex) Parse(l ILines) (INotation, ILines, error) {
	//get first match
	next := l.Next()
	matches := r.pattern.FindAllString(next, 1)
	//log.Debugf("next=%s pattern=%s -> %d matches %+v", next, r.patternString, len(matches), matches)
	if len(matches) >= 1 {
		//matched, but could be later in the line and could be multiple
		//we only consider the first match
		//in which case SkipOver() will fail
		if remain, ok := l.SkipOver(matches[0]); ok {
			r.match = strings.TrimRight(matches[0], " \t")
			log.Debugf("regex(%s) parsed \"%s\" on line %d, next line %d: %.32s", r.Name(), r.match, l.LineNr(), remain.LineNr(), remain.Next())
			return r, remain, nil
		}
	}
	return nil, l, log.Wrapf(nil, "regex(%s) mismatch line %d: %.32s", r.patternString, l.LineNr(), next)
} //regex.Parse()
