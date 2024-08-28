package middlewares

import (
	"reflect"

	"github.com/microcosm-cc/bluemonday"
)

func XSSMiddleware(param interface{}) interface{} {
	policy := bluemonday.UGCPolicy()
	val := reflect.ValueOf(param)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return param
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if field.Kind() == reflect.String {
			sanitizedValue := policy.Sanitize(field.String())
			field.SetString(sanitizedValue)
		} else if field.Kind() == reflect.Struct {
			sanitizedNested := XSSMiddleware(field.Addr().Interface())
			field.Set(reflect.ValueOf(sanitizedNested).Elem())
		}
	}

	return param
}
