package validate

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Tag struct {
	//	Tag    string
	Eq  string `method:"Equals"`
	Gt  string `method:"GreaterThan"`
	Gte string `method:"GreaterThanOrEqual"`
	Lt  string `method:"LessThan"`
	Lte string `method:"LessThanOrEqual"`
	Ne  string `method:"NotEqual"`

	Min    int `method:"MaxValidate"`
	Max    int `method:"MinValidate"`
	Length int `method:"LengthValidate"`

	Email   bool `method:"EmailValidate"`
	Url     bool `method:"UrlValidate"`
	NoSpace bool `method:"NoSpaceValidate"`
}

func buildTags(tag string, opt interface{}) {
	tags := strings.Split(tag, ";")
	optReflect := reflect.ValueOf(opt)
	// 判读那是否为指针类型，或这元素是否可设置值
	if optReflect.Kind() != reflect.Ptr || !optReflect.Elem().CanSet() {
		return
	}

	label := ""
	value := ""
	optReflect = optReflect.Elem()
	for _, item := range tags {
		kv := strings.Split(strings.TrimSpace(item), "=")

		if len(kv) == 1 {
			label = Capitalize(kv[0])
			value = "true"
		}
		if len(kv) == 2 {
			label = Capitalize(kv[0])
			value = kv[1]
		}

		//注意这里只有公开字段才可以设置，不然会报错
		field := optReflect.FieldByName(label)
		if ok := field.IsValid(); !ok {
			continue
		}

		switch field.Kind() {
		case reflect.String:
			field.SetString(value)
		case reflect.Bool:
			v, err := strconv.ParseBool(value)
			if err != nil {
				field.SetBool(false)
			}
			field.SetBool(v)
		case reflect.Float32, reflect.Float64:
			v, err := strconv.ParseFloat(value, 64)
			if err != nil {
				field.SetFloat(0.00)
			}
			field.SetFloat(v)
		case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
			v, err := strconv.Atoi(value)
			if err != nil {
				field.SetInt(0)
			}
			field.SetInt(int64(v))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			v, err := strconv.ParseUint(value, 0, 64)
			if err != nil {
				field.SetUint(0)
			}
			field.SetUint(v)
		default:
		}
	}
}

func (n Tag) MaxValidate(v reflect.Value) (bool, string) {
	length := len(v.String())
	if length <= n.Max {
		return true, ""
	}
	return false, fmt.Sprintf("no more than %d characters, but %d characters were entered", n.Max, length)
}

func (n Tag) NoSpaceValidate(v reflect.Value) (bool, string) {
	chars := []rune(v.String())
	for _, char := range chars {
		if char == ' ' {
			return false, fmt.Sprintf("cannot contain spaces characters")
		}
	}
	return true, ""
}

func (n Tag) MinValidate(v reflect.Value) (bool, string) {
	length := len(v.String())
	if length >= n.Min {
		return true, ""
	}
	return false, fmt.Sprintf("no less than %d characters, but %d characters were entered", n.Min, length)
}

func (n Tag) EmailValidate(v reflect.Value) bool {
	email := v.String()
	return EMAIL_REG.MatchString(email)
}

func (n Tag) UrlValidate(v reflect.Value) (bool, string) {
	url := v.String()
	return URL_REG.MatchString(url), "email address format is incorrect"
}

func (n Tag) LengthValidate(v reflect.Value) (bool, string) {
	length := len(v.String())
	if length != n.Length {
		return false , fmt.Sprintf("required %d characters, but %d characters were entered", n.Length, length)
	}
	return true, ""
}

func (n Tag) NotEqual(v reflect.Value) (bool, string) {
	switch v.Kind() {
	case reflect.Float32, reflect.Float64:
		eq, err := strconv.ParseFloat(n.Ne, 64)
		if err != nil || v.Float() == eq {
			return false, fmt.Sprintf("cannot be equal to %f", eq)
		}
		return true, ""

	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		eq, err := strconv.ParseInt(n.Ne, 10, 64)
		if err != nil || v.Int() == eq {
			return false, fmt.Sprintf("cannot be equal to %d", eq)
		}
		return true, ""

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		eq, err := strconv.ParseUint(n.Ne, 10, 64)
		if err != nil || v.Uint() == eq {
			return false, fmt.Sprintf("cannot be equal to %d", eq)
		}
		return true, ""
	}

	return true, ""
}

func (n Tag) LessThanOrEqual(v reflect.Value) (bool, string) {
	switch v.Kind() {
	case reflect.Float32, reflect.Float64:
		eq, err := strconv.ParseFloat(n.Lte, 64)
		if err != nil || v.Float() > eq {
			return false, fmt.Sprintf("less than or equal %f", eq)
		}
		return true, ""

	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		eq, err := strconv.ParseInt(n.Lte, 10, 64)
		if err != nil || v.Int() > eq {
			return false, fmt.Sprintf("less than or equal to %d", eq)
		}
		return true, ""

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		eq, err := strconv.ParseUint(n.Lte, 10, 64)
		if err != nil || v.Uint() > eq {
			return false, fmt.Sprintf("less than or equal to %d", eq)
		}
		return true, ""
	}

	return true, ""
}

func (n Tag) LessThan(v reflect.Value) (bool, string) {
	switch v.Kind() {
	case reflect.Float32, reflect.Float64:
		eq, err := strconv.ParseFloat(n.Lt, 64)
		if err != nil || v.Float() >= eq {
			return false, fmt.Sprintf("less than %f", eq)
		}
		return true, ""

	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		eq, err := strconv.ParseInt(n.Lt, 10, 64)
		if err != nil || v.Int() >= eq {
			return false, fmt.Sprintf("less than %d", eq)
		}
		return true, ""

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		eq, err := strconv.ParseUint(n.Lt, 10, 64)
		if err != nil || v.Uint() >= eq {
			return false, fmt.Sprintf("less than %d", eq)
		}
		return true, ""
	}

	return true, ""
}

func (n Tag) GreaterThanOrEqual(v reflect.Value) (bool, string) {
	switch v.Kind() {
	case reflect.Float32, reflect.Float64:
		eq, err := strconv.ParseFloat(n.Gte, 64)
		if err != nil || v.Float() < eq {
			return false, fmt.Sprintf("greate than or equal to %f", eq)
		}
		return true, ""

	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		eq, err := strconv.ParseInt(n.Gte, 10, 64)
		if err != nil || v.Int() < eq {
			return false, fmt.Sprintf("greate than or equal to %d", eq)
		}
		return true, ""

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		eq, err := strconv.ParseUint(n.Gte, 10, 64)
		if err != nil || v.Uint() < eq {
			return false, fmt.Sprintf("greate than or equal to %d", eq)
		}
		return true, ""
	}

	return true, ""
}

func (n Tag) Equals(v reflect.Value) (bool, string) {
	switch v.Kind() {
	case reflect.Float32, reflect.Float64:
		eq, err := strconv.ParseFloat(n.Eq, 64)
		if err != nil || v.Float() != eq {
			return false, fmt.Sprintf("equal to %f", eq)
		}
		return true, ""

	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		eq, err := strconv.ParseInt(n.Eq, 10, 64)
		if err != nil || v.Int() != eq {
			return false, fmt.Sprintf("equal to %d", eq)
		}
		return true, ""

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		eq, err := strconv.ParseUint(n.Eq, 10, 64)
		if err != nil || v.Uint() != eq {
			return false, fmt.Sprintf("equal to %d", eq)
		}
		return true, ""
	}

	return true, ""
}

func (n Tag) GreaterThan(v reflect.Value) (bool, string) {
	switch v.Kind() {
	case reflect.Float32, reflect.Float64:
		eq, err := strconv.ParseFloat(n.Gt, 64)
		if err != nil || v.Float() <= eq {
			return false, fmt.Sprintf("greate than %f", eq)
		}
		return true, ""

	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		eq, err := strconv.ParseInt(n.Gt, 10, 64)
		if err != nil || v.Int() <= eq {
			return false, fmt.Sprintf("greate than %d", eq)
		}
		return true, ""

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		eq, err := strconv.ParseUint(n.Gt, 10, 64)
		if err != nil || v.Uint() <= eq {
			return false, fmt.Sprintf("greate than %d", eq)
		}
		return true, ""
	}

	return true, ""
}
