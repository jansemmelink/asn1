package parser

//INotation is something that can be parsed
type INotation interface {
	Name() string
	Parse(l ILines) (INotation, ILines, error)
}

//New notation
func New(n string) INotation {
	return notation{name: n}
}

type notation struct {
	name string
}

func (n notation) Name() string {
	return n.name
}

func (n notation) Parse(l ILines) (INotation, ILines, error) {
	return nil, l, log.Wrapf(nil, "%T.Parse() not implemented", n)
}
