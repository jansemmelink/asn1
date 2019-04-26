package def2

import (
	"reflect"

	"bitbucket.org/vservices/dark/logger"
	"github.com/jansemmelink/asn1/parser"
	"github.com/jansemmelink/mem"
)

//Seq is a sequence of parsable fields
type Seq struct{}

//Parse ...
func (Seq) Parse(log logger.ILogger, l parser.ILines, v IParsable) (parser.ILines, error) {
	remain := l
	structType := reflect.TypeOf(v).Elem()
	log.Debugf("=====[ PARSING SEQ %v ]=====", structType)
	for i := 0; i < structType.NumField(); i++ {
		structTypeField := structType.Field(i)
		if structTypeField.Anonymous {
			//log.Debugf("  %s.Field[%d]=%s is embedded - ignored", structType.Name(), i, structTypeField.Name)
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

		parsable, ok := mem.NewX(fieldValueType).(IParsable)
		if !ok {
			//log.Debugf("  Field[%d]=%s type=%v is not parsable - ignored", i, structTypeField.Name, fieldValueType)
			continue
		}

		//log.Debugf("  Field[%d]=%s type=%v parsing ...", i, structTypeField.Name, fieldValueType)
		var err error
		remain, err = parsable.Parse(log, remain, parsable)
		if err != nil {
			return l, log.Wrapf(err, "Failed to parse %s.%s from line %d: %.32s ...", structType.Name(), structTypeField.Name, remain.LineNr(), remain.Next())
		}
		if structTypeField.Type.Kind() == reflect.Ptr {
			//log.Debugf("  Field[%d]=%s type=%v set ptr ...", i, structTypeField.Name, fieldValueType)
			reflect.ValueOf(v).Elem().Field(i).Set(reflect.ValueOf(parsable))
		} else {
			//log.Debugf("  Field[%d]=%s type=%v set value ...", i, structTypeField.Name, fieldValueType)
			reflect.ValueOf(v).Elem().Field(i).Set(reflect.ValueOf(parsable).Elem())
		}
	} //for each field
	return remain, nil
} //Seq.Parse()
