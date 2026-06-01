package luamapper

import (
	"reflect"
)

// Maps struct values into slices that can be mapped to lua table through custom `lua` struct tags
func MapStructToSlice(s any) *map[string]any {
	t := reflect.TypeOf(s)
	val := reflect.ValueOf(s)

	if t.Kind() == reflect.Pointer {
		t = t.Elem()
		val = val.Elem()
	}

	fieldAmount := t.NumField()
	mappped := make(map[string]any, fieldAmount)

	for i := range fieldAmount {
		field := t.Field(i)
		key, exists := field.Tag.Lookup("lua")
		if exists {
			mappped[key] = val.Field(i).Interface()
		}
	}

	return &mappped
}
