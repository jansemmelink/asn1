package parser

import (
	"fmt"
	"sync"

	"bitbucket.org/vservices/dark/logger"
)

//INotation is a named entity that can be parsed
type INotation interface {
	Name() string
	ParseV(log logger.ILogger, l ILines) (IValue, ILines, error)
}

var (
	notations      = make(map[string]INotation)
	notationsMutex = sync.Mutex{}
)

//New notation
func New(typeName, name string) INotation {
	notationsMutex.Lock()
	defer notationsMutex.Unlock()

	if existing, ok := notations[name]; ok {
		panic(log.Wrapf(nil, "Duplicate Notation(%s): %s = %p", typeName, existing.Name(), existing))
		//return existing
	}

	n := &notation{
		parent: nil,
		name:   name,
	}
	log.Debugf("%p(%s) = %s", n, typeName, n.Name())
	notations[name] = n
	return n
}

type notation struct {
	parent INotation
	name   string
}

func (n notation) Name() string {
	if n.parent != nil {
		return n.parent.Name() + "." + n.name
	}
	return n.name
}

func (n notation) Parse(log logger.ILogger, l ILines) (INotation, ILines, error) {
	panic(log.Wrapf(nil, "%T.Parse() not implemented", n))
}

func (n notation) ParseV(log logger.ILogger, l ILines) (IValue, ILines, error) {
	panic(log.Wrapf(nil, "%T.ParseV() not implemented", n))
}

//IValue from a parser
type IValue interface {
	Name() string
	Value() interface{}
	String() string
	Int() int
	Items() []IValue
}

//NewValue ...
func NewValue(name string) IValue {
	return value{
		name: name,
	}
}

type value struct {
	name string
}

func (v value) Name() string {
	return v.name
}

func (v value) Value() interface{} {
	panic(log.Wrapf(nil, "%T(%s).Value() not implemented", v, v.name))
}

func (v value) String() string {
	panic(log.Wrapf(nil, "%T(%s).String() not implemented", v, v.name))
}

func (v value) Int() int {
	panic(log.Wrapf(nil, "%T(%s).Int() not implemented", v, v.name))
}

func (v value) Items() []IValue {
	panic(log.Wrapf(nil, "%T(%s).Items() not implemented", v, v.name))
}

//IIntValue ...
type IIntValue interface {
	IValue
}

//IntValue ...
func IntValue(name string, value int) IIntValue {
	return intValue{
		IValue: NewValue(name),
		i:      value,
	}
}

type intValue struct {
	IValue
	i int
}

func (iv intValue) Value() interface{} {
	return iv.i
}

func (iv intValue) Int() int {
	return iv.i
}

func (iv intValue) String() string {
	return fmt.Sprintf("%d", iv.i)
}

//IStrValue ...
type IStrValue interface {
	IValue
}

//StrValue ...
func StrValue(name, value string) IStrValue {
	return strValue{
		IValue: NewValue(name),
		s:      value,
	}
}

type strValue struct {
	IValue
	s string
}

func (sv strValue) Value() interface{} {
	return sv.s
}

func (sv strValue) String() string {
	return sv.s
}

//IListValue ...
type IListValue interface {
	IValue
	With(v IValue) IListValue
}

//ListValue ...
func ListValue(name string) IListValue {
	return listValue{
		IValue: NewValue(name),
		items:  make([]IValue, 0),
	}
}

type listValue struct {
	IValue
	items []IValue
}

func (lv listValue) Value() interface{} {
	return lv.items
}

func (lv listValue) Items() []IValue {
	return lv.items
}

func (lv listValue) String() string {
	return fmt.Sprintf("LIST:%s", lv.Name())
}

func (lv listValue) With(v IValue) IListValue {
	lv.items = append(lv.items, v)
	return lv
}

//IValueV ...
type IValueV interface {
	IValue
}

//ValueV ... (for any value)
func ValueV(name string, value interface{}) IValueV {
	return valueV{
		IValue: NewValue(name),
		value:  value,
	}
}

type valueV struct {
	IValue
	value interface{}
}

func (v valueV) Value() interface{} {
	return v.value
}
