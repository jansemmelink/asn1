package def2

import (
	"reflect"

	"bitbucket.org/vservices/dark/logger"
	"github.com/jansemmelink/asn1/parser"
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

		temp := reflect.New(fieldValueType)
		parsable, ok := temp.Interface().(IParsable)
		if !ok {
			log.Debugf("  Field[%d]=%s type=%v is not parsable - ignored", i, structTypeField.Name, fieldValueType)
			if fieldValueType.Name() == "Identifier" {
				tryParsable := temp.Interface().(IParsable)
				tryParsable.Parse(log, l, nil)
			}
			continue
		}

		var err error
		remain, err = parsable.Parse(log, remain, parsable)
		if err != nil {
			return l, log.Wrapf(err, "Failed to parse %s.%s from line %d: %.32s ...", structType.Name(), structTypeField.Name, remain.LineNr(), remain.Next())
		}
		if structTypeField.Type.Kind() == reflect.Ptr {
			reflect.ValueOf(v).Elem().Field(i).Set(temp)
		} else {
			reflect.ValueOf(v).Elem().Field(i).Set(temp.Elem())
		}
	} //for each field
	return remain, nil
} //Seq.Parse()
