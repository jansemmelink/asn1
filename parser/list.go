package parser

import (
	"fmt"

	"bitbucket.org/vservices/dark/logger"
)

//List creates a parser for a list of the same types of items
func List(sep string, min int, item INotation) IList {
	name := fmt.Sprintf("list of %s", item.Name())
	lst := &list{
		INotation: New("list", name),
		sep:       sep,
		min:       min,
	}
	lst.expectedItem = item
	return lst
}

//IList extends INotation
type IList interface {
	INotation
}

//list implements IList
type list struct {
	INotation
	sep          string
	min          int
	expectedItem INotation
}

func (lst list) ParseV(log logger.ILogger, l ILines) (IValue, ILines, error) {
	remain := l
	parsedList := ListValue(lst.Name())
	var err error
	for {
		//get next item
		var parsedValue IValue
		if parsedValue, remain, err = lst.expectedItem.ParseV(log.With("list(%s)", lst.Name()), remain); err != nil {
			break
		}
		parsedList = parsedList.With(parsedValue)
		log.Debugf("list(%s) now has %d items, next=%s", lst.Name(), len(parsedList.Items()), remain.Next())

		//needs separater before allowing another item
		//space separators are implicit, so do not check
		//todo: not really safe... need better checks here as last item could be a number followed immediately by a bracket or new-line etc...
		if len(lst.sep) > 0 && lst.sep != " " {
			var ok bool
			if remain, ok = remain.SkipOver(lst.sep); !ok {
				log.Debugf("NO SEP - END LIST: next=%s", remain.Next())
				break
			}
			//log.Debugf("list(%s) got sep=\"%s\", next=%s", lst.Name(), lst.sep, remain.Next())
		}
	}

	if len(parsedList.Items()) < lst.min {
		return nil, l, log.Wrapf(err, "list(%s) has %d items < min=%d", lst.Name(), len(parsedList.Items()), lst.min)
	}
	log.Debugf("list(%s) parsed %d items from line %d, next line %d: %.32s", lst.Name(), len(parsedList.Items()), l.LineNr(), remain.LineNr(), remain.Next())
	return parsedList, remain, nil
} //list.ParseV()
