package helper

import (
	"fmt"
	"reflect"
	"strings"
)

type Fields struct {
	Key      string
	Value    string
	Identity string
	DataType string
	FType    string
}

func GetFields(i interface{}) (flds []Fields) {
	f := new(Fields)

	for _, fieldname := range FieldNames(i) {
		field, ok := reflect.TypeOf(i).Elem().FieldByName(fieldname)
		if !ok {
			panic("Field not found")
		}
		f.Key = fieldname
		tagdata := strings.Split(GetStructTag(field, "stt"), ",")
		f.FType = fmt.Sprintf("%v", field.Type.Kind())
		f.Value = tagdata[0]
		f.Identity = ""
		f.DataType = ""

		if len(tagdata) > 1 {

			f.DataType = tagdata[1]
		}
		if len(tagdata) > 2 {
			f.Identity = strings.Join(tagdata[2:], " ")
		}

		flds = append(flds, *f)
	}

	return flds

}
func FieldNames(i interface{}) (names []string) {

	val := reflect.ValueOf(i).Elem()
	for i := 0; i < val.NumField(); i++ {
		val.Field(i).Interface()
		names = append(names, val.Type().Field(i).Name)
	}
	return names
}
func GetValues(i interface{}) (values []string) {
	val := reflect.ValueOf(i).Elem()
	for i := 0; i < val.NumField(); i++ {
		if len(fmt.Sprintf("%v", val.Field(i).Interface())) > 0 {
			values = append(values, GetSqlString(val.Field(i).Interface()))
		}
	}
	return values
}
func GetStructTag(f reflect.StructField, tagName string) string {
	return string(f.Tag.Get(tagName))
}
func StructName(_struct interface{}) string {
	if t := reflect.TypeOf(_struct); t.Kind() == reflect.Ptr {
		return strings.ToLower(t.Elem().Name())
	}
	return ""
}
func GetOffSet(i int) string {
	if i > 0 {
		return fmt.Sprint(" OFFSET ", i)
	}
	return ""
}
func GetLimit(i int) string {
	if i > 0 {
		return fmt.Sprint(" LIMIT ", i)
	}
	return ""
}
func GetSqlString(i interface{}) string {
	switch i.(type) {
	case string:
		return fmt.Sprintf(`'%v'`, i)
	default:
		return fmt.Sprintf(`%v`, i)
	}
}
