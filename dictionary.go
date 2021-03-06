package collections

import (
	"reflect"
)

type (
	Dictionary interface {
		Map(m ...interface{}) interface{}
		Keys() List
		Values() List
		Where(f interface{}) Dictionary
		Count(f ...interface{}) int
		ForEach(f interface{}) Dictionary
		Select(f interface{}) Dictionary
		Merge(d Dictionary, onConflict ...interface{}) Dictionary

		Type() reflect.Type
	}

	dictionary struct {
		t     reflect.Type
		value *reflect.Value
	}
)

func newDictionary(t reflect.Type, cap ...int) *dictionary {
	var value reflect.Value
	if len(cap) == 0 {
		value = reflect.MakeMap(t)
	} else {
		value = reflect.MakeMapWithSize(t, cap[0])
	}
	return &dictionary{t: t, value: &value}
}

//Map 获取当前映射集合的值
//可以传入具名 map 以确保符合预期，或者使用返回值进行类型断言。
func (dict *dictionary) Map(m ...interface{}) interface{} {
	if len(m) == 0 {
		return dict.value.Interface()
	}
	val, dst := dict.value, reflect.Indirect(reflect.ValueOf(m[0]))
	if err := typeRequired(dst.Type(), dict.t); err != nil {
		panic(err)
	}
	for _, key := range val.MapKeys() {
		dst.SetMapIndex(key, val.MapIndex(key))
	}
	return dst.Interface()
}

//Keys 获取当前映射集合的键集
//直接返回 List 类型用于持续计算。
func (dict *dictionary) Keys() List {
	val := dict.value
	var keys = newList(reflect.SliceOf(dict.t.Key()), val.Len())
	for index, key := range val.MapKeys() {
		keys.value.Index(index).Set(key)
	}
	return keys
}

//Values 获取当前映射集合的值集
//直接返回 List 类型用于持续计算。
func (dict *dictionary) Values() List {
	val := dict.value
	var values = newList(reflect.SliceOf(dict.t.Elem()), val.Len())
	for index, key := range val.MapKeys() {
		values.value.Index(index).Set(val.MapIndex(key))
	}
	return values
}

func (dict *dictionary) Where(f interface{}) Dictionary {
	val, function := dict.value, reflect.ValueOf(f)
	if err := typeRequired(function.Type(),
		newFunc(dict.t.Key())(types.Bool)(),
		newFunc(dict.t.Key(), dict.t.Elem())(types.Bool)(),
	); err != nil {
		panic(err)
	}
	var newmap, numin = newDictionary(dict.t), function.Type().NumIn()
	for _, key := range val.MapKeys() {
		var args = []reflect.Value{key, val.MapIndex(key)}
		if call(function, args[:numin]...)[0].Bool() {
			newmap.value.SetMapIndex(key, args[1])
		}
	}
	return newmap
}

func (dict *dictionary) Count(f ...interface{}) int {
	if len(f) > 0 {
		return dict.Where(f[0]).Count()
	}
	return dict.value.Len()
}

func (dict *dictionary) ForEach(f interface{}) Dictionary {
	val, function := dict.value, reflect.ValueOf(f)
	if err := typeRequired(function.Type(),
		newFunc(dict.t.Key())(types.AnyTypes)(),
		newFunc(dict.t.Key(), dict.t.Elem())(types.AnyTypes)(),
	); err != nil {
		panic(err)
	}
	var numin = function.Type().NumIn()
	for _, key := range val.MapKeys() {
		var args = []reflect.Value{key, val.MapIndex(key)}
		var back = call(function, args[:numin]...)
		if len(back) > 0 {
			var first, last = back[0], back[len(back)-1].Interface()
			var stop = first.Kind() == reflect.Bool && !first.Bool()
			if !stop && last != nil {
				_, stop = last.(error)
			}
			if stop {
				break
			}
		}
	}
	return dict
}

func (dict *dictionary) Select(f interface{}) Dictionary {
	val, function := dict.value, reflect.ValueOf(f)
	var funct = function.Type()
	if err := typeRequired(funct,
		newFunc(dict.t.Key())(types.AnyTypes)(),
		newFunc(dict.t.Key(), dict.t.Elem())(types.AnyTypes)(),
	); err != nil {
		panic(err)
	}
	var kt, vt = dict.t.Key(), funct.Out(0)
	if funct.NumOut() > 1 {
		kt, vt = funct.Out(0), funct.Out(1)
	}
	var newmap, numin = newDictionary(reflect.MapOf(kt, vt), val.Len()), function.Type().NumIn()
	dict.ForEach(func(k, v interface{}) {
		var args = []reflect.Value{reflect.ValueOf(k), reflect.ValueOf(v)}
		var back, key = call(function, args[:numin]...), args[0]
		if len(back) > 1 {
			key, back[0] = back[0], back[1]
		}
		newmap.value.SetMapIndex(key, back[0])
	})
	return newmap
}

func (dict *dictionary) copy() *dictionary {
	var newdict = newDictionary(dict.t)
	dict.ForEach(func(k, v interface{}) {
		newdict.value.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(v))
	})
	return newdict
}

func (dict *dictionary) Type() reflect.Type {
	return dict.t
}

func (dict *dictionary) Merge(d Dictionary, onConflict ...interface{}) Dictionary {
	if err := typeRequired(d.Type(), dict.t); err != nil {
		panic(err)
	}
	var function reflect.Value
	if len(onConflict) == 0 {
		function = makeConflictHandler(d.Type().Elem(), 1)
	} else {
		function = reflect.ValueOf(onConflict[0])
	}
	var funct = function.Type()
	if err := typeRequired(funct,
		newFunc(dict.t.Elem(), types.AnyTypes)(dict.t.Elem())(),
	); err != nil {
		panic(err)
	}
	var newmap = dict.copy()
	d.ForEach(func(k, v interface{}) {
		var key = reflect.ValueOf(k)
		if !newmap.value.MapIndex(key).IsValid() {
			newmap.value.SetMapIndex(key, reflect.ValueOf(v))
		} else {
			var old = newmap.value.MapIndex(key)
			newmap.value.SetMapIndex(key, call(function, old, reflect.ValueOf(v))[0])
		}
	})
	return newmap
}

func makeConflictHandler(t reflect.Type, startIndex int) reflect.Value {
	var funct = reflect.FuncOf([]reflect.Type{t, t}, []reflect.Type{t}, false)
	var funcbody = func(args []reflect.Value) []reflect.Value {
		return args[startIndex : startIndex+1]
	}
	var function = reflect.MakeFunc(funct, funcbody)
	return function
}
