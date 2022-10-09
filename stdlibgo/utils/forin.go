package utils

import (
	"fmt"
	"reflect"
)

var (
	ErrForInDestNil = fmt.Errorf("dest cannot be nil")
)

func ForIn(dest interface{}, iteratee func(key interface{}, value interface{}) error) (err error) {
	strct := reflect.ValueOf(dest)
	strctKind := strct.Kind()
	if strctKind == reflect.Pointer && strct.IsNil() {
		err = ErrForInDestNil
		return
	}

	elemKind := reflect.Indirect(strct).Kind()
	switch elemKind {
	case reflect.Map:
		err = ForInMap(dest, iteratee)
	default:
		err = ForInStruct(dest, iteratee)
	}

	return
}

func ForInMap(dest interface{}, iteratee func(key interface{}, value interface{}) error) (err error) {
	strct := reflect.ValueOf(dest)
	if strct.Kind() == reflect.Pointer && strct.IsNil() {
		err = ErrForInDestNil
		return
	}

	values := reflect.Indirect(strct).MapKeys()
	for _, k := range values {
		v := strct.MapIndex(k)
		if err = iteratee(k.Interface(), v.Interface()); err != nil {
			return
		}
	}

	return
}

func ForInStruct(dest interface{}, iteratee func(key interface{}, value interface{}) error) (err error) {
	strct   := reflect.ValueOf(dest)
	if strct.Kind() == reflect.Pointer && strct.IsNil() {
		err = ErrForInDestNil
		return
	}

	typeOfT := reflect.Indirect(strct).Type()
	for i := 0; i < reflect.Indirect(strct).NumField(); i++ {
		v := reflect.Indirect(strct).Field(i).Interface()
		k := typeOfT.Field(i).Name
		if err = iteratee(k, v); err != nil {
			return
		}
	}

	return
}