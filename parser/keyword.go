package parser

import (
	// "sync"

	"bitbucket.org/vservices/dark/logger"
)

// var (
// 	keywords      = make(map[string]IKeyword)
// 	keywordsMutex = sync.Mutex{}
// )

//Keyword creates a parser that fails if the specified word is not present
func Keyword(word string) IKeyword {
	// keywordsMutex.Lock()
	// defer keywordsMutex.Unlock()
	// if existing, ok := keywords[word]; ok {
	// 	log.Debugf("Existing(keyword): %s", existing.Name())
	// 	return existing
	// }

	k := &keyword{
		INotation: New("keyword", word),
		word:      word,
	}
	// keywords[word] = k
	return k
}

//IKeyword extends INotation
type IKeyword interface {
	INotation
}

type keyword struct {
	INotation
	word string
}

func (k keyword) ParseV(log logger.ILogger, l ILines) (IValue, ILines, error) {
	remain, ok := l.SkipOver(k.word)
	if !ok {
		return nil, l, log.Wrapf(nil, "Line %d: keyword(%s) not present", l.LineNr(), k.word)
	}
	return StrValue(k.Name(), k.word), remain, nil
} //keyword.Parse()
