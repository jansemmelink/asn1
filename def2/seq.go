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
	if reflect.TypeOf(v).Kind() != reflect.Ptr {
		return l, log.Wrapf(nil, "%T.Parse(log,lines,%T): output arg 3 is not a pointer!", v, v)
	}
	structType := reflect.TypeOf(v).Elem()
	log.Debugf("=====[ PARSING SEQ %v ]=====", structType)
	for i := 0; i < structType.NumField(); i++ {
		structTypeField := structType.Field(i)
		if structTypeField.Anonymous {
			log.Debugf("  %s.Field[%d]=%s is embedded - ignored", structType.Name(), i, structTypeField.Name)
			continue
		}
		if !reflect.ValueOf(v).Elem().Field(i).CanSet() {
			log.Debugf("  %s.Field[%d]=%s is private - ignored", structType.Name(), i, structTypeField.Name)
			continue
		}

		//get field value type
		//if its a pointer field, dereference it:
		fieldValueType := structTypeField.Type
		log.Debugf("  %s.Field[%d]=%s kind=%s type=%s", structType.Name(), i, structTypeField.Name, structTypeField.Type.Kind(), structTypeField.Type.Name())
		if structTypeField.Type.Kind() == reflect.Ptr {
			fieldValueType = fieldValueType.Elem()
			log.Debugf("Dereferenced %s!", fieldValueType.Name())
		}

		parsable, ok := mem.NewX(fieldValueType).(IParsable)
		if !ok {
			log.Debugf("  %s.Field[%d]=%s type=%v is not parsable - ignored", structType.Name(), i, structTypeField.Name, fieldValueType.Name())
			continue
		}

		log.Debugf("  %s.Field[%d]=%s type=%v parsing ...", structType.Name(), i, structTypeField.Name, fieldValueType)
		var err error
		remain, err = parsable.Parse(log, remain, parsable)
		if err != nil {
			//did not parse
			if structTypeField.Type.Kind() == reflect.Ptr {
				//pointer fields are optional
				//failed to parse here, just continue with next field
				reflect.ValueOf(v).Elem().Field(i).Set(reflect.Zero(structTypeField.Type))
				log.Debugf("%s.Field[%d]=%s is absent", structType.Name(), i, structTypeField.Name)
			} else {
				return l, log.Wrapf(err, "Failed to parse %s.Field[%d]=%s from line %d: %.32s ...", structType.Name(), i, structTypeField.Name, remain.LineNr(), remain.Next())
			}
		} else {
			//parsed
			if structTypeField.Type.Kind() == reflect.Ptr {
				//log.Debugf("  Field[%d]=%s type=%v set ptr ...", i, structTypeField.Name, fieldValueType)
				reflect.ValueOf(v).Elem().Field(i).Set(reflect.ValueOf(parsable))
			} else {
				//log.Debugf("  Field[%d]=%s type=%v set value ...", i, structTypeField.Name, fieldValueType)
				reflect.ValueOf(v).Elem().Field(i).Set(reflect.ValueOf(parsable).Elem())
			}
		} //if parsed
	} //for each field
	return remain, nil
} //Seq.Parse()
