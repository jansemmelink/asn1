package parser

import (
	"regexp"
	"strings"

	"bitbucket.org/vservices/dark/logger"
)

//Regex creates a regex parser
func Regex(name, pattern string) IRegex {
	compiled, err := regexp.Compile(pattern)
	if err != nil {
		panic(log.Wrapf(err, "Invalid regular expression pattern: %s", pattern))
	}

	r := &regex{
		INotation:     New("regex", name),
		patternString: pattern,
		pattern:       compiled,
	}
	log.Debugf("%p(regex) = %s", r, r.Name())
	return r
}

//IRegex matches a regular expression
type IRegex interface {
	INotation
}

type regex struct {
	INotation
	patternString string
	pattern       *regexp.Regexp
}

func (r regex) ParseV(log logger.ILogger, l ILines) (IValue, ILines, error) {
	log.Debugf("%T(%s).Parse() from line %d: %.32s ...", r, r.Name(), l.LineNr(), l.Next())
	//get first match
	next := l.Next()
	matches := r.pattern.FindAllString(next, 1)
	//log.Debugf("next=%s pattern=%s -> %d matches %+v", next, r.patternString, len(matches), matches)
	if len(matches) >= 1 {
		//matched, but could be later in the line and could be multiple
		//we only consider the first match
		//in which case SkipOver() will fail
		if remain, ok := l.SkipOver(matches[0]); ok {
			match := strings.TrimRight(matches[0], " \t")
			log.Debugf("regex(%s) parsed \"%s\" on line %d, next line %d: %.32s", r.Name(), match, l.LineNr(), remain.LineNr(), remain.Next())
			return StrValue(r.Name(), match), remain, nil
		}
	}
	return nil, l, log.Wrapf(nil, "regex(%s) mismatch line %d: %.32s", r.patternString, l.LineNr(), next)
} //regex.Parse()
