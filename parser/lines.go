package parser

import "strings"

//NewLines ...
func NewLines() ILines {
	return lines{
		array: make([]ILine, 0),
	}
}

//ILines is a set of significant lines read from a file
type ILines interface {
	Count() int
	Append(int, string, string) ILines

	//read token from next line
	//use this when looking for single keyword
	Next() string
	//line nr of the next line being read
	LineNr() int

	//TBD: read with dynamic buffer
	//the output buffer can be extended to read further as needed
	//and drop the buffer when done
	//...

	SkipOver(s string) (ILines, bool)
}

//lines implement ILines
type lines struct {
	array []ILine
}

func (l lines) Count() int {
	return len(l.array)
}

func (l lines) Append(nr int, text string, comment string) ILines {
	l.array = append(l.array, line{nr: nr, text: strings.Trim(text, " \t"), comment: comment})
	return l
}

//get next text to parse - from the next line only
func (l lines) Next() string {
	if len(l.array) == 0 {
		return ""
	}
	i := 0
	for i < len(l.array) && len(l.array[i].Text()) == 0 {
		i++
	}
	if i < len(l.array) {
		return l.array[i].Text()
	}
	return ""
}

func (l lines) LineNr() int {
	if len(l.array) == 0 {
		return 0
	}
	i := 0
	for i < len(l.array) && len(l.array[i].Text()) == 0 {
		i++
	}
	if i < len(l.array) {
		return l.array[i].Nr()
	}
	return 0
}

//Skip over one token - should all be on the next line
//if this is the last text in the line, the line is dropped
func (l lines) SkipOver(s string) (ILines, bool) {
	sl := len(s)
	if sl == 0 {
		//nothing to skip
		log.Debugf("NOTHING TO SKIP")
		return l, false
	}

	if len(l.array) == 0 {
		//no more lines of text, cannot skip over
		log.Debugf("NO MORE LINES")
		return l, false
	}

	//skip in first non-empty line
	var ok bool
	i := 0
	for i < len(l.array) && l.array[i].Text() == "" {
		//log.Debugf("Skip empty line %d", l.array[i].Nr())
		i++
	}
	//l.array[i], ok = l.array[i].SkipOver(s)
	modifiedLine, ok := l.array[i].SkipOver(s)
	if !ok {
		//token not present on next line
		//log.Debugf("NOT PRESENT: \"%s\" != \"%s\"", s, l.array[i].Text())
		return l, false
	}

	//skipped
	//warning: do not modify array elements in object
	//it will change the object even though its suppose to be const!
	//first make a copy, apply the change, then return the copy
	//we only need to change the specified item in the array
	//and we also drop all empty lines at the top
	c := lines{array: make([]ILine, 0)}
	if len(modifiedLine.Text()) > 0 {
		//some text remain in the modifier line, start with this
		c.array = append(c.array, modifiedLine)
	} else {
		//modified line emptied, start with next line but trim leading spaces
		//log.Debugf("Modified emptied: i=%d len(array)=%d", i, len(l.array))
		if i < len(l.array)-1 {
			nextLine, _ := l.array[i+1].SkipOver("")
			c.array = append(c.array, nextLine)
			i++
		}
	}
	if i < len(l.array) {
		c.array = append(c.array, l.array[i+1:]...)
	}

	//log.Debugf("SKIPPED(%s): next %5d:%s", s, c.LineNr(), c.Next())
	// log.Debugf("  SRC(#=%5d): %5d:%s", len(l.array), l.array[0].Nr(), l.array[0].Text())
	// log.Debugf("  DST(#=%5d): %5d:%s", len(c.array), c.array[0].Nr(), c.array[0].Text())
	return c, true
}

//ILine from a text file
type ILine interface {
	Nr() int
	Text() string //without leading/trailing spaces
	Comment() string

	//skip over and return line without specified token
	SkipOver(s string) (ILine, bool)
}

//line implements ILine
type line struct {
	nr      int
	text    string
	comment string
}

func (l line) Nr() int {
	return l.nr
}

func (l line) Text() string {
	return l.text
}

func (l line) Comment() string {
	return l.comment
}

func (l line) SkipOver(s string) (ILine, bool) {
	sl := len(s)
	if sl > 0 {
		tl := len(l.text)
		if tl < sl {
			//not enough text to skip
			return l, false
		}

		if l.text[0:sl] != s {
			//no match
			return l, false
		}
	} //if has text to skip

	//skip over (even over "") and trim spaces after the token
	l.text = strings.TrimLeft(l.text[sl:], " \t")
	return l, true
}
