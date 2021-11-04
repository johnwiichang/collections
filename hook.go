package collections

import (
	"reflect"
	"sync"
)

var (
	compareHooks = &sync.Map{}
)

func getCompareHook(caller, target reflect.Type) *reflect.Value {
	storage, existed := compareHooks.Load(caller)
	if !existed {
		storage = &sync.Map{}
		compareHooks.Store(caller, storage)
	}
	function, existed := storage.(*sync.Map).Load(target)
	if !existed {
		function = getHook(caller, "EqualsTo*", newFunc(caller, target)(types.Bool)())
		storage.(*sync.Map).Store(target, function)
	}
	return function.(*reflect.Value)
}
