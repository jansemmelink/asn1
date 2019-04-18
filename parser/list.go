package parser

import "fmt"

//List creates a parser for a list of the same types of items
func List(sep string, min int, item INotation) IList {
	return list{
		INotation:    New(fmt.Sprintf("list of %s", item.Name())),
		sep:          sep,
		min:          min,
		expectedItem: item,
		items:        nil,
	}
}

//IList extends INotation
type IList interface {
	INotation
	Items() []INotation
}

//list implements IList
type list struct {
	INotation
	//specification
	sep          string
	min          int
	expectedItem INotation
	//output
	items []INotation
}

func (lst list) Items() []INotation {
	return lst.items
}

func (lst list) Parse(l ILines) (INotation, ILines, error) {
	lst.items = make([]INotation, 0)
	remain := l
	//log.Debugf("list(%s): Start from line %d: %.32s", lst.Name(), remain.LineNr(), remain.Next())
	for {
		//log.Debugf("  list(%s).item[%d] from: %.32s", lst.Name(), len(lst.items), remain.Next())
		//get next item
		var err error
		var parsedItem INotation
		if parsedItem, remain, err = lst.expectedItem.Parse(remain); err != nil {
			break
		}
		lst.items = append(lst.items, parsedItem)
		log.Debugf("list(%s) now has %d items, next=%s", lst.Name(), len(lst.items), remain.Next())

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

	if len(lst.items) < lst.min {
		return nil, l, log.Wrapf(nil, "list(%s) has %d items < min=%d", lst.Name(), len(lst.items), lst.min)
	}
	log.Debugf("list(%s) parsed %d items from line %d, next line %d: %.32s", lst.Name(), len(lst.items), l.LineNr(), remain.LineNr(), remain.Next())
	return lst, remain, nil
} //list.Parse()
