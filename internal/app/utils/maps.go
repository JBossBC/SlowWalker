package utils

import (
	"errors"
	"fmt"
	"reflect"
)

func Format(model any) map[string]any {
	maps := make(map[string]any)
	typ := reflect.TypeOf(model)
	valueOf := reflect.ValueOf(model)
	for i := 0; i < typ.NumField(); i++ {
		key := typ.Field(i).Name
		value := valueOf.Field(i).Interface()
		maps[key] = value
	}
	return maps
}

func Parse(data map[string]any, template any) (err error) {
	defer func() {
		if panicError := recover(); panicError != nil {
			err = errors.New(fmt.Sprintf("%v", panicError))
		}
	}()
	typeof := reflect.TypeOf(template)
	value := reflect.ValueOf(template)
	for i := 0; i < typeof.NumField(); i++ {
		typ := typeof.Field(i)
		val := value.Field(i)
		if !val.CanSet() {
			continue
		}
		mapVal, ok := data[typ.Name]
		if !ok {
			continue
		}
		val.Set(reflect.ValueOf(mapVal))
	}
	return nil
}
