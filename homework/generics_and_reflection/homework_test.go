package main

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Person struct {
	Name    string `properties:"name"`
	Address string `properties:"address,omitempty"`
	Age     int    `properties:"age"`
	Married bool   `properties:"married"`
}

func Serialize(person Person) string {
	value := reflect.ValueOf(person)
	typeValue := value.Type()
	fields := make([]string, 0, typeValue.NumField())

	for i := 0; i < typeValue.NumField(); i++ {
		field := typeValue.Field(i)
		key, option := parseTag(field)
		if key == "" {
			continue
		}

		fieldValue := value.Field(i)
		if needSkipField(option, fieldValue) {
			continue
		}

		valueString := formatValue(fieldValue)
		fields = append(fields, key+"="+valueString)
	}

	return strings.Join(fields, "\n")
}

func parseTag(field reflect.StructField) (key string, opts []string) {
	tag := field.Tag.Get("properties")
	if tag == "" {
		return "", nil
	}

	parts := strings.Split(tag, ",")
	return parts[0], parts[1:]
}

func needSkipField(opts []string, val reflect.Value) bool {
	for _, opt := range opts {
		if opt == "omitempty" && val.IsZero() {
			return true
		}
	}

	return false
}

func formatValue(value reflect.Value) string {
	switch value.Kind() {
	case reflect.String:
		return value.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%d", value.Int())
	case reflect.Bool:
		return fmt.Sprintf("%t", value.Bool())
	default:
		return fmt.Sprintf("%v", value.Interface())
	}
}

func TestSerialization(t *testing.T) {
	tests := map[string]struct {
		person Person
		result string
	}{
		"test case with empty fields": {
			result: "name=\nage=0\nmarried=false",
		},
		"test case with fields": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
			},
			result: "name=John Doe\nage=30\nmarried=true",
		},
		"test case with omitempty field": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
				Address: "Paris",
			},
			result: "name=John Doe\naddress=Paris\nage=30\nmarried=true",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := Serialize(test.person)
			assert.Equal(t, test.result, result)
		})
	}
}
