package def2

import (
	"reflect"

	"bitbucket.org/vservices/dark/logger"
	"github.com/jansemmelink/asn1/parser"
	"github.com/jansemmelink/mem"
)

//Choice parses until one of the struct members succeeds
type Choice struct {
	//index in parent struct type of field currently used
	//0 = none, which would be this embedded struct if used
	sf ChoiceSF
}

//ChoiceSF describes the selected field in a choice
type ChoiceSF struct {
	Index int
	Name  string
	Type  reflect.Type
	Value reflect.Value
	Ptr   IParsable
}

//Parse ...
func (Choice) Parse(log logger.ILogger, l parser.ILines, v IParsable) (parser.ILines, error) {
	structType := reflect.TypeOf(v).Elem()
	log.Debugf("=====[ PARSING CHOICE %v ]=====", structType)
	for i := 0; i < structType.NumField(); i++ {
		structTypeField := structType.Field(i)
		if structTypeField.Anonymous {
			//log.Debugf("  Field[%d]=%s is embedded - ignored", i, structTypeField.Name)
			continue
		}
		if !reflect.ValueOf(v).Elem().Field(i).CanSet() {
			//log.Debugf("  Field[%d]=%s is private - ignored", i, structTypeField.Name)
			continue
		}

		//get field value type
		//if its a pointer field, dereference it:
		fieldValueType := structTypeField.Type
		if structTypeField.Type.Kind() == reflect.Ptr {
			fieldValueType = fieldValueType.Elem()
		}

		if parsable, ok := mem.NewX(fieldValueType).(IParsable); ok {
			remain, err := parsable.Parse(log, l, parsable)
			if err != nil {
				//parsing failed - try another option
				log.Debugf("%s.%s %s is Absent %v", structType.Name(), structTypeField.Name, fieldValueType.Name(), err)
				continue
			}

			//this option parsed successfully
			//no need to continue
			if structTypeField.Type.Kind() == reflect.Ptr {
				reflect.ValueOf(v).Elem().Field(i).Set(reflect.ValueOf(parsable))
			} else {
				reflect.ValueOf(v).Elem().Field(i).Set(reflect.ValueOf(parsable).Elem())
			}

			//indicate the choice selection
			var c *Choice
			c = reflect.ValueOf(v).Elem().Field(0).Addr().Interface().(*Choice)
			c.sf = ChoiceSF{
				Index: i,
				Name:  structTypeField.Name,
				Type:  structTypeField.Type,
				Value: reflect.ValueOf(v).Elem().Field(i),
				Ptr:   reflect.ValueOf(v).Elem().Field(i).Addr().Interface().(IParsable),
			}

			return remain, nil
		} //if parsable option in the choice
	} //for each field
	return l, log.Wrapf(nil, "%v none of the options matched in line %d: %.32s ...", structType.Name(), l.LineNr(), l.Next())
} //Choice.Parse()

//Option ...
func (c Choice) Option() ChoiceSF {
	return c.sf
}

//Value ...
func (c Choice) Value() interface{} {
	return c.sf.Value.Interface()
}
