package validate

import (
	"fmt"
	"reflect"
	"regexp"
)

const (
	IgnoreFields = "-"
	emailPattern = `^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$`
	urlPattern   = `/(http|https):\/\/([\w.]+\/?)\S*/　`
)

var (
	EMAIL_REG = regexp.MustCompile(emailPattern)
	URL_REG   = regexp.MustCompile(emailPattern)
)

// Validate field validate
// if validate ,return "", true
// if not, return the first invalid field name, and false
func Validate(i interface{}) error {

	refValue := reflect.ValueOf(i)
	refType := reflect.TypeOf(i)

	// 传入的是指针情况，需要使用Elem()获取元素
	if refValue.Kind() == reflect.Ptr {
		refValue = reflect.ValueOf(i).Elem()
		refType = reflect.TypeOf(i).Elem()
	}


	for i := 0; i < refType.NumField(); i++ {
		field := refValue.Field(i)
		types := refType.Field(i)

		//if field.Kind() == reflect.Ptr {
		//	field = refValue.Field(i).Elem()
		//	types = refType.Field(i)
		//}

		tag := types.Tag.Get("validate")
		if tag == IgnoreFields {
			continue
		}

		switch field.Kind() {
		case reflect.Float32, reflect.Float64,
			reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if err := verification(field, tag); err != nil {
				return fmt.Errorf("\"%s\" %s", types.Name, err.Error())
			}
		case reflect.String:
			if err := verification(field, tag); err != nil {
				return fmt.Errorf("\"%s\" %s", types.Name, err.Error())
			}

		case reflect.Struct:
			if err := Validate(field.Interface()); err != nil {
				return fmt.Errorf("\"%s\".%s", types.Name, err.Error())
			}
		case reflect.Array, reflect.Slice:
			for j := 0; j < field.Len(); j++ {
				//Validate(field.Index(j).Interface())
				if err := Validate(field.Index(j).Interface()); err != nil {
					return fmt.Errorf("\"%s[%d]\".%s", types.Name, j, err.Error())
				}
			}

		}
	}
	return nil
}

func verification(v interface{}, tags string) error {

	opt := &Tag{}
	buildTags(tags, opt)

	optType := reflect.TypeOf(opt).Elem()
	optValue := reflect.ValueOf(opt).Elem()
	for i := 0; i < optType.NumField(); i++ {
		filed := optType.Field(i)
		value := optValue.Field(i)
		if value.IsZero() {
			continue
		}
		if methodName, ok := filed.Tag.Lookup("method"); ok {
			method := optValue.MethodByName(methodName)

			param := []reflect.Value{
				reflect.ValueOf(v),
			}
			results := method.Call(param)

			if len(results) == 1 && !results[0].Bool()  {
				return fmt.Errorf("does not satisfy the condition of %s ", filed.Name)
			}

			if len(results) == 2 {
				if !results[0].Bool() {
					return fmt.Errorf("does not satisfy the condition of %s ( %s )", filed.Name, results[1].String())
				}

			}
		}
	}

	return nil
}
