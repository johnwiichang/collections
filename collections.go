package collections

import (
	"reflect"
)

type (
	collections reflect.Value

	Collections interface {
		List() List
		Dictionary() Dictionary
	}
)

func From(obj interface{}) Collections {
	var val = collections(reflect.Indirect(reflect.ValueOf(obj)))
	return &val
}

func (collections *collections) Value() *reflect.Value {
	return (*reflect.Value)(collections)
}

//List 获取 List 集合
//如果类型不为 List，那么会抛出 panic 异常。如果给出数组，将会自动转换为 Slice（如果不可求址则拷贝）。
func (collections *collections) List() List {
	var kind = collections.Value().Kind()
	if kind != reflect.Slice && kind != reflect.Array {
		panic(throwTypeNotCompatiable("List", collections.Value().Type()).Error())
	}
	var value = collections.Value()
	if kind == reflect.Array {
		if value.CanAddr() {
			slice := value.Slice(0, value.Len())
			value = &slice
		} else {
			slice := reflect.New(reflect.SliceOf(value.Type().Elem())).Elem()
			reflect.Copy(slice, *value)
			value = &slice
		}
	}
	return &list{t: value.Type(), value: value}
}

func (collections *collections) Dictionary() Dictionary {
	if collections.Value().Kind() != reflect.Map {
		panic(throwTypeNotCompatiable("Dictionary", collections.Value().Type()).Error())
	}
	return &dictionary{t: collections.Value().Type(), value: collections.Value()}
}
