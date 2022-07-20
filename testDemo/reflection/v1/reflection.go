package main

import "reflect"

// func walk(x interface{}, fn func(input string)) {
// 	val := reflect.ValueOf(x)
// 	for i := 0; i < val.NumField(); i++ {
// 		field := val.Field(i)
// 		if field.Kind() == reflect.String {
// 			fn(field.String())
// 		}
// 		if field.Kind() == reflect.Struct {
// 			walk(field.Interface(), fn)
// 		}

// 	}
// }

func walk(x interface{}, fn func(input string)) {
	val := reflect.ValueOf(x)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)

		switch field.Kind() {
		case reflect.String:
			fn(field.String())
		case reflect.Struct:
			walk(field.Interface(), fn)
		}
	}
}