package collections

import (
	"reflect"
	"sort"
)

type (
	List interface {
		Slice(slice ...interface{}) interface{}
		Select(f interface{}) List
		SelectMany(f interface{}) List
		ForEach(f interface{}) List
		ToDictionary(f ...interface{}) Dictionary
		Sort(less ...interface{}) List
		Reverse() List
		Distinct() List
		Where(f interface{}) List
		Count(f ...interface{}) int
		Contains(elements ...interface{}) bool
		Any(elements ...interface{}) bool
		Concat(l List) List
		First(obj interface{}) int
		Last(obj interface{}) int
		Intersect(l List) List
		Except(l List) List
		Union(l List) List
		Skip(length int) List
		Take(num int) List
		Resize(length ...int) List

		Type() reflect.Type
	}

	list struct {
		t     reflect.Type
		value *reflect.Value

		skip int
	}
)

func newList(t reflect.Type, cap ...int) *list {
	var value reflect.Value
	if len(cap) == 0 {
		value = reflect.New(t).Elem()
	} else {
		value = reflect.MakeSlice(t, cap[0], cap[0])
	}
	return &list{t: t, value: &value}
}

func (lst *list) Type() reflect.Type {
	return lst.t
}

func (lst *list) Skip(length int) List {
	lst.skip = length
	return lst
}

//Resize 从跳过后选择一定数量的元素创建新列表（元素不够时不报错但是长度会不足）
func (lst *list) Take(num int) List {
	var newlist, length = newList(lst.t), lst.Count()
	if length > lst.skip && num > 0 {
		if max := length - lst.skip; max < num {
			num = max
		}
		newlist.value.Set(lst.value.Slice(lst.skip, lst.skip+num))
	}
	return newlist
}

//Resize 更改当前列表的范围（就地操作）后还原跳过标识
func (lst *list) Resize(length ...int) List {
	defer func() {
		lst.skip = 0
	}()
	var maxlength, take, skip = lst.Count(), 0, lst.skip
	if lst.skip < maxlength {
		take = maxlength - lst.skip
	} else {
		lst.skip = maxlength
	}
	if len(length) > 0 && length[0] >= 0 && length[0] < take {
		take = length[0]
	}
	lst.value.Set(lst.value.Slice(skip, skip+take))
	return lst
}

func (lst *list) Select(f interface{}) List {
	function := reflect.ValueOf(f)
	if err := typeRequired(function.Type(),
		//支持的函数签名
		newFunc(types.Int, lst.t.Elem())(types.AnyTypes)(),
		newFunc(lst.t.Elem())(types.AnyTypes)(),
	); err != nil {
		panic(err)
	}
	var newlist = newList(reflect.SliceOf(function.Type().Out(0)), lst.value.Len())
	var numin = function.Type().NumIn()
	lst.ForEach(func(i int, val interface{}) {
		var args = []reflect.Value{reflect.ValueOf(i), reflect.ValueOf(val)}
		newlist.value.Index(i).Set(call(function, args[2-numin:2]...)[0])
	})
	return newlist
}

func (lst *list) SelectMany(f interface{}) List {
	function := reflect.ValueOf(f)
	if err := typeRequired(function.Type(),
		//支持的函数签名
		newFunc(types.Int, lst.t.Elem())(types.Slice)(),
		newFunc(lst.t.Elem())(types.Slice)(),
	); err != nil {
		panic(err)
	}
	var newlist = newList(reflect.SliceOf(function.Type().Out(0).Elem()))
	var numin = function.Type().NumIn()
	lst.ForEach(func(i int, val interface{}) {
		var args = []reflect.Value{reflect.ValueOf(i), reflect.ValueOf(val)}
		newlist.value.Set(reflect.AppendSlice(*newlist.value, call(function, args[2-numin:2]...)[0]))
	})
	return newlist
}

func (lst *list) ForEach(f interface{}) List {
	val, function := lst.value, reflect.ValueOf(f)
	if err := typeRequired(function.Type(),
		//支持的函数签名
		newFunc()(types.AnyTypes)(),
		newFunc(types.Int, lst.t.Elem())(types.AnyTypes)(),
		newFunc(lst.t.Elem())(types.AnyTypes)(),
	); err != nil {
		panic(err)
	}
	var numin = function.Type().NumIn()
	for i := 0; i < val.Len(); i++ {
		var args = []reflect.Value{reflect.ValueOf(i), val.Index(i)}
		var back = call(function, args[2-numin:2]...)
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
	return lst
}

func (lst *list) Slice(slice ...interface{}) interface{} {
	if len(slice) == 0 {
		return lst.value.Interface()
	}
	val, dst := lst.value, reflect.Indirect(reflect.ValueOf(slice[0]))
	if err := typeRequired(dst.Type(), lst.t); err != nil {
		panic(err)
	}
	var targetLength = dst.Len()
	for i := 0; i < val.Len(); i++ {
		if i < targetLength {
			dst.Index(i).Set(val.Index(i))
		} else {
			dst.Set(reflect.Append(dst, val.Index(i)))
		}
	}
	return dst.Interface()
}

func (lst *list) ToDictionary(f ...interface{}) Dictionary {
	f = append(f, func() interface{} { return nil })
	val, function := lst.value, reflect.ValueOf(f[0])
	if err := typeRequired(function.Type(),
		//支持的函数签名
		newFunc()(types.AnyType, types.AnyTypes)(),
		newFunc(types.Int, lst.t.Elem())(types.AnyType, types.AnyTypes)(),
		newFunc(lst.t.Elem())(types.AnyType, types.AnyTypes)(),
	); err != nil {
		panic(err)
	}
	var funct = function.Type()
	var kt, vt = lst.t.Elem(), funct.Out(0)
	if funct.NumOut() > 1 {
		kt, vt = funct.Out(0), funct.Out(1)
	}
	var newmap, numin = newDictionary(reflect.MapOf(kt, vt), val.Len()), function.Type().NumIn()
	lst.ForEach(func(i int, val interface{}) {
		var args = []reflect.Value{reflect.ValueOf(i), reflect.ValueOf(val)}
		var back = call(function, args[2-numin:2]...)
		if len(back) > 1 {
			newmap.value.SetMapIndex(back[0], back[1])
		} else {
			newmap.value.SetMapIndex(args[1], back[0])
		}
	})
	return newmap
}

func (lst *list) Distinct() List {
	elem := lst.t.Elem()
	compare := getCompareHook(elem, elem)
	if compare == nil {
		var function = reflect.ValueOf(func(a, b interface{}) bool {
			return a == b
		})
		compare = &function
	}
	var values []reflect.Value
	lst.ForEach(func(obj interface{}) {
		var existed bool
		for _, value := range values {
			if existed = compare.Call([]reflect.Value{reflect.ValueOf(obj), value})[0].Bool(); existed {
				break
			}
		}
		if !existed {
			values = append(values, reflect.ValueOf(obj))
		}
	})
	var newlist = newList(lst.t)
	newlist.value.Set(reflect.Append(*newlist.value, values...))
	return newlist
}

func (lst *list) Count(f ...interface{}) int {
	if len(f) > 0 {
		return lst.Where(f[0]).Count()
	}
	return lst.value.Len()
}

func (lst *list) Where(f interface{}) List {
	compare := reflect.ValueOf(f)
	if err := typeRequired(compare.Type(), newFunc(lst.t.Elem())(types.Bool)()); err != nil {
		panic(err)
	}
	var values []reflect.Value
	lst.ForEach(func(obj interface{}) {
		if compare.Call([]reflect.Value{reflect.ValueOf(obj)})[0].Bool() {
			values = append(values, reflect.ValueOf(obj))
		}
	})
	var newlist = newList(lst.t)
	newlist.value.Set(reflect.Append(*newlist.value, values...))
	return newlist
}

func (lst *list) Contains(elements ...interface{}) bool {
	if length := len(elements); length == 0 {
		return lst.Count() > 0
	}
	for _, element := range elements {
		if !lst.Any(element) {
			return false
		}
	}
	return true
}

func (lst *list) Any(elements ...interface{}) bool {
	if length := len(elements); length == 0 {
		return lst.Count() > 0
	} else if length == 1 {
		if l, ok := elements[0].(List); ok {
			elements = l.Select(func(x interface{}) interface{} { return x }).Slice().([]interface{})
			return lst.Any(elements...)
		}
	}
	for i := 0; i < lst.value.Len(); i++ {
		for _, element := range elements {
			if valueCompare(lst.value.Index(i), reflect.ValueOf(element)) {
				return true
			}
		}
	}
	return false
}

func (lst *list) Sort(less ...interface{}) List {
	var val = lst.value
	var slice sort.Interface
	switch lst.t.Elem().Kind() {
	case reflect.Int:
		slice = sort.IntSlice(val.Interface().([]int))
	case reflect.Float64:
		slice = sort.Float64Slice(val.Interface().([]float64))
	case reflect.String:
		slice = sort.StringSlice(val.Interface().([]string))
	default:
		if len(less) == 0 {
			panic(throwMethodHasNoImplement("less", lst.t).Error())
		}
	}
	if len(less) > 0 {
		var function = reflect.ValueOf(less[0])
		var elem = lst.t.Elem()
		var estimate = newFunc(elem, elem)(types.Bool)()
		if err := typeRequired(function.Type(), estimate); err != nil {
			panic(err)
		}
		sort.Slice(val.Interface(), func(i, j int) bool {
			var args = []reflect.Value{val.Index(i), val.Index(j)}
			return call(function, args...)[0].Bool()
		})
	}
	sort.Sort(slice)
	return lst
}

func (lst *list) Reverse() List {
	var val = lst.value
	var length = lst.value.Len()
	for i := 0; i <= length/2; i++ {
		var temp = reflect.ValueOf(val.Index(i).Interface())
		val.Index(i).Set(val.Index(length - 1 - i))
		val.Index(length - 1 - i).Set(temp)
	}
	return lst
}

func (lst *list) Concat(l List) List {
	if err := typeRequired(l.Type(), lst.t); err != nil {
		panic(err)
	}
	var newlist = newList(lst.t)
	newlist.value.Set(reflect.AppendSlice(*lst.value, reflect.ValueOf(l.Slice())))
	return newlist
}

func (lst *list) First(obj interface{}) (index int) {
	index = -1
	var function = reflect.ValueOf(obj)
	if function.Kind() != reflect.Func {
		function = reflect.ValueOf(func(element interface{}) bool {
			return obj == element
		})
	}
	lst.ForEach(func(i int, item interface{}) (next bool) {
		if next = !call(function, reflect.ValueOf(item))[0].Bool(); !next {
			index = i
		}
		return
	})
	return
}

func (lst *list) Last(obj interface{}) (index int) {
	index = -1
	var function = reflect.ValueOf(obj)
	if function.Kind() != reflect.Func {
		function = reflect.ValueOf(func(element interface{}) bool {
			return obj == element
		})
	}
	defer lst.Reverse()
	lst.Reverse().ForEach(func(i int, item interface{}) (next bool) {
		if next = !call(function, reflect.ValueOf(item))[0].Bool(); !next {
			index = i
		}
		return
	})
	return
}

func (lst *list) Intersect(l List) List {
	var newlist = newList(lst.t)
	lst.ForEach(func(item interface{}) {
		if l.Any(item) {
			newlist.value.Set(reflect.Append(*newlist.value, reflect.ValueOf(item)))
		}
	})
	return newlist
}

func (lst *list) Except(l List) List {
	var newlist = newList(lst.t)
	lst.ForEach(func(item interface{}) {
		if !l.Any(item) {
			newlist.value.Set(reflect.Append(*newlist.value, reflect.ValueOf(item)))
		}
	})
	return newlist
}

func (lst *list) Union(l List) List {
	return lst.Concat(l).Distinct()
}
