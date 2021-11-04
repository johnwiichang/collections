package collections_test

import (
	"testing"

	"github.com/johnwiichang/collections"
)

func TestFrom(t *testing.T) {
	var array = [2]int{1, 2}
	collections.From(array).List().Slice()
	collections.From(&array).List().Slice()
	collections.From(map[string]interface{}{}).Dictionary()
	EstimateFail(t, func(*testing.T) {
		collections.From([]int{}).Dictionary()
	})
	EstimateFail(t, func(*testing.T) {
		collections.From(map[string]interface{}{}).List()
	})
}

//EstimateFail 本函数执行的方法预期会抛出 panic 错误
func EstimateFail(t *testing.T, function func(*testing.T)) {
	defer func() {
		if recover() == nil {
			t.Fail()
		}
	}()
	function(t)
}
