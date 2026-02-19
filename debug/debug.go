package debug

import (
	"fmt"
	"reflect"
)

// InspectStruct takes any object and prints its fields and values.
func InspectStruct(s any) {
	v := reflect.ValueOf(s)

	// If it's a pointer, get the object it points to
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Safety check: Ensure we are looking at a struct
	if v.Kind() != reflect.Struct {
		fmt.Printf("InspectStruct error: expected struct, got %s\n", v.Kind())
		return
	}

	t := v.Type()
	fmt.Printf("--- Inspecting Type: %s ---\n", t.Name())

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		val := v.Field(i)

		// Handle nested pointers to values (like *string)
		var actualValue any
		if val.Kind() == reflect.Ptr && !val.IsNil() {
			actualValue = val.Elem().Interface()
		} else {
			actualValue = val.Interface()
		}

		// Get the JSON tag if it exists
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" {
			jsonTag = "n/a"
		}

		fmt.Printf("%-15s | Tag: %-15s | Value: %v\n", field.Name, jsonTag, actualValue)
	}
}
