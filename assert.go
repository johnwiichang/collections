package collections

import (
	"reflect"
	"strings"
)

type (
	anyType  struct{}
	anyTypes struct{}
)

var (
	types = struct {
		AnyType, AnyTypes, Bool, Int, Slice reflect.Type
	}{
		reflect.TypeOf(anyType{}), reflect.TypeOf(anyTypes{}), reflect.TypeOf(true), reflect.TypeOf(0), reflect.SliceOf(reflect.TypeOf(anyType{})),
	}

	compareResults = struct{ NotMatch, Match, MatchAndStop int }{-1, 0, 1}
)

//typeRequired 类型要求断言
func typeRequired(target reflect.Type, estimate ...reflect.Type) error {
	for _, t := range estimate {
		if result := typeCompare(target, t); result != compareResults.NotMatch {
			return nil
		}
	}
	types := From(estimate).List().Select(func(x reflect.Type) string { return x.String() }).Slice().([]string)
	if len(estimate) > 0 {
		return throwTypeNotCompatiable(strings.Join(types, "', '"), target)
	}
	return nil
}

//typeCompare 类型匹配
func typeCompare(target reflect.Type, estimate reflect.Type) int {
	//对于任意数量的匹配，那么返回匹配且终止后续匹配
	if estimate == types.AnyTypes {
		return compareResults.MatchAndStop
	} else if estimate == types.AnyType || target.Kind() == reflect.Interface {
		//如果任意类型匹配，则返回匹配
		return compareResults.Match
	} else if target.Kind() != estimate.Kind() {
		//类型不同不用匹配
		return compareResults.NotMatch
	}
	switch target.Kind() {
	case reflect.Func:
		var ein, tin = estimate.NumIn(), target.NumIn()
		var eout, tout = estimate.NumOut(), target.NumOut()
		var cursor int
		for ; cursor < ein; cursor++ {
			//如果是往后任意类型，那么直接跳过后续输入匹配
			if estimate.In(cursor) == types.AnyTypes {
				break
			}
			//如果目标函数输入项缺失，则不匹配
			if cursor >= tin {
				return compareResults.NotMatch
			}
			//执行平凡类型匹配
			var result = typeCompare(target.In(cursor), estimate.In(cursor))
			if result == compareResults.NotMatch {
				//不匹配立即返回
				return compareResults.NotMatch
			} else if result == compareResults.MatchAndStop {
				//匹配且终止则跳出
				break
			}
		}
		//输出匹配
		for cursor = 0; cursor < eout; cursor++ {
			//如果是往后任意类型，那么直接跳过后续输出匹配
			if estimate.Out(cursor) == types.AnyTypes {
				break
			}
			//如果目标函数输出项缺失，则不匹配
			if cursor >= tout {
				return compareResults.NotMatch
			}
			//执行平凡类型匹配
			var result = typeCompare(target.Out(cursor), estimate.Out(cursor))
			if result != compareResults.Match {
				//如果匹配且终止、不匹配那么直接返回
				return result
			}
		}
		return compareResults.Match
	case reflect.Slice, reflect.Array:
		return typeCompare(target.Elem(), estimate.Elem())
	default:
		if estimate == target || target.ConvertibleTo(estimate) {
			return compareResults.Match
		}
		return compareResults.NotMatch
	}
}

func valueCompare(v1, v2 reflect.Value) bool {
	t1, t2 := v1.Type(), v2.Type()
	if function := getCompareHook(t1, t2); function != nil {
		if call(*function, v1, v2)[0].Bool() {
			return true
		}
	} else if function = getCompareHook(t2, t1); function != nil {
		if call(*function, v2, v1)[0].Bool() {
			return true
		}
	}
	return v1.Interface() == v2.Interface()
}

//newFunc 柯里化函数类型声明
func newFunc(input ...reflect.Type) func(...reflect.Type) func(...bool) reflect.Type {
	return func(output ...reflect.Type) func(...bool) reflect.Type {
		return func(variadic ...bool) reflect.Type {
			return reflect.FuncOf(input, output, append(variadic, false)[0])
		}
	}
}

//getHook 获取对象钩子，附带类型检查。
//会使用 getBlurMatchFunction 以支持带星号（*）的模糊匹配。
func getHook(t reflect.Type, name string, estimate ...reflect.Type) *reflect.Value {
	if function := getBlurMatchFunction(name); function == nil {
		if method, existed := t.MethodByName(name); existed {
			if typeRequired(method.Type, estimate...) == nil {
				return &method.Func
			}
		}
	} else {
		length := t.NumMethod()
		for i := 0; i < length; i++ {
			var method = t.Method(i)
			if function(method.Name) && typeRequired(method.Type, estimate...) == nil {
				return &method.Func
			}
		}
	}
	return nil
}

//getBlurMatchFunction 获取模糊匹配函数（如果不属于模糊匹配，则会返回空函数）
// `*key ` → 后缀匹配
// ` key*` → 前缀匹配
// `*key*` → 包含匹配
func getBlurMatchFunction(name string) func(string) bool {
	if strings.HasPrefix(name, "*") {
		name = name[1:]
		if strings.HasSuffix(name, "*") {
			name = name[:len(name)-1]
			return func(s string) bool {
				return strings.Contains(s, name)
			}
		} else {
			return func(s string) bool {
				return strings.HasSuffix(s, name)
			}
		}
	} else if strings.HasSuffix(name, "*") {
		name = name[:len(name)-1]
		return func(s string) bool {
			return strings.HasPrefix(s, name)
		}
	}
	return nil
}

//call 转换类型并调用函数
func call(f reflect.Value, in ...reflect.Value) []reflect.Value {
	t := f.Type()
	var args = make([]reflect.Value, t.NumIn())
	for index := range in {
		args[index] = in[index].Convert(t.In(index))
	}
	return f.Call(in)
}
