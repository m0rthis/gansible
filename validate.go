package gansible

import (
	"fmt"
	"reflect"
	"strings"
)

func validateRequiredFields[T any](args *T, result *Result) {
	validateValue(reflect.ValueOf(args).Elem(), result)
}

func validateValue(v reflect.Value, result *Result) {
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldVal := v.Field(i)

		// Recurse into embedded structs
		if field.Anonymous && field.Type.Kind() == reflect.Struct {
			validateValue(fieldVal, result)
			continue
		}

		if field.Tag.Get("required") == "true" {
			if isDefaultValue(fieldVal.Interface()) {
				result.Msg = fmt.Sprintf("Field '%s' is required", field.Name)
				FailJson(*result)
			}
		}
	}
}

func isDefaultValue(val any) bool {
	if val == nil {
		return true
	}
	v := reflect.ValueOf(val)
	zero := reflect.Zero(v.Type())
	return reflect.DeepEqual(v.Interface(), zero.Interface())
}

func filterArgsByOS[T any](args *T, result *Result) {
	currentOS, _ := GetOS()
	v := reflect.ValueOf(args).Elem()
	t := v.Type()

	var issues []string

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("os")
		if tag != "" && OS(tag) != currentOS {
			val := v.Field(i)
			if !isDefaultValue(val.Interface()) {
				issues = append(issues, fmt.Sprintf("`%s` (requires %s)", field.Name, tag))
			}
		}
	}

	if len(issues) > 0 {
		result.Msg = fmt.Sprintf(
			"Invalid arguments for %s: %s",
			currentOS,
			strings.Join(issues, ", "),
		)
		FailJson(*result)
	}
}
