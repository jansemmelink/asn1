package def2

import (
	"regexp"
	"strings"
	"sync"

	"bitbucket.org/vservices/dark/logger"
	"github.com/jansemmelink/asn1/parser"
)

var (
	compiled      = make(map[string]*regexp.Regexp)
	compiledMutex = sync.Mutex{}
)

//RegexParse function to use for custom regular expressions
//Example:
//----------------------------------------------------------------------------------------------------------
//	type MyItem = string
//  func (item *MyItem) Parse(log logger.ILogger, l parser.ILines, v def2.IParsable) (parser.ILines, error) {
// 	  return def2.RegexParse(log, l, "[a-zA-Z][a-zA-Z0-9]*", (*string)(item))
//  }
//
//  type MyParsableStruct struct {
//		def2.Seq
//		...
//		item MyItem				<-- use it for example here in a sequence of parsable things
//		...
//	}
//----------------------------------------------------------------------------------------------------------
func RegexParse(log logger.ILogger, l parser.ILines, pattern string, value *string) (parser.ILines, error) {
	compiledMutex.Lock()
	defer compiledMutex.Unlock()
	regex, ok := compiled[pattern]
	if !ok {
		var err error
		regex, err = regexp.Compile(pattern)
		if err != nil {
			panic(log.Wrapf(err, "Regex(\"%s\") compilation failed", pattern))
		}
		log.Debugf("Compiled Regex(%s)", pattern)
		compiled[pattern] = regex
	}

	// log.Debugf("REGEX %T(%s)", out, out.Pattern())
	next := l.Next()
	matches := regex.FindAllString(next, 1)
	if len(matches) >= 1 {
		//matched, but could be later in the line and could be multiple
		//we only consider the first match
		//in which case SkipOver() will fail
		if remain, ok := l.SkipOver(matches[0]); ok {
			*value = strings.TrimRight(matches[0], " \t")
			log.Debugf("regex(%s) parsed \"%s\" on line %d, next line %d: %.32s ...", pattern, *value, l.LineNr(), remain.LineNr(), remain.Next())
			return remain, nil
		}
	}
	return l, log.Wrapf(nil, "regex(%s) mismatch line %d: %.32s", pattern, l.LineNr(), next)
}
