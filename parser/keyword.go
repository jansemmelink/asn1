package parser

//Keyword creates a parser that fails if the specified word is not present
func Keyword(word string) IKeyword {
	return keyword{
		INotation: New("keyword:" + word),
		word:      word,
	}
}

//IKeyword extends INotation
type IKeyword interface {
	INotation
}

type keyword struct {
	INotation
	word string
}

func (k keyword) Parse(l ILines) (INotation, ILines, error) {
	remain, ok := l.SkipOver(k.word)
	if !ok {
		return nil, l, log.Wrapf(nil, "Line %d: keyword(%s) not present", l.LineNr(), k.word)
	}
	log.Debugf("Line %5d: keyword(%s), next=%s", l.LineNr(), k.word, remain.Next())
	return k, remain, nil
} //keyword.Parse()
