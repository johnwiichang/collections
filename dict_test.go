package collections_test

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/johnwiichang/collections"
)

var dicts = struct {
	NumberWithTrue collections.Dictionary
}{
	collections.From(map[int]bool{1: true, 2: true, 3: true, 4: true, 5: true}).Dictionary(),
}

func TestDictionaryBasicCopy(t *testing.T) {
	var dict, dst = map[int]bool{1: true, 2: true, 3: true, 4: true, 5: true}, make(map[int]bool)
	if reflect.DeepEqual(dict, dicts.NumberWithTrue.Map(&dst)) {
		if reflect.DeepEqual(dict, dst) {
			EstimateFail(t, func(*testing.T) {
				dicts.NumberWithTrue.Map(new(map[int]int))
			})
		}
	}
}

func TestDictionarySelect(t *testing.T) {
	var dst = slices.Number.ToDictionary(func(n int) string { return strconv.Itoa(n) }).Map()
	var test = dicts.NumberWithTrue.Select(func(k int, v bool) (int, string) {
		return k, strconv.Itoa(k)
	})
	if !reflect.DeepEqual(test.Map(), dst) {
		t.Fail()
	}
}

func TestDictionaryKeyAndValueList(t *testing.T) {
	var keys = dicts.NumberWithTrue.Keys().Sort()
	var values = dicts.NumberWithTrue.Values()
	if !reflect.DeepEqual(keys.Slice(), slices.Number.Slice()) {
		t.Fail()
	}
	if !values.Distinct().Slice().([]bool)[0] {
		t.Fail()
	}
}

func TestDictionaryQuery(t *testing.T) {
	var f = func(k int, v bool) bool { return v }
	if dicts.NumberWithTrue.Count(f) != dicts.NumberWithTrue.Where(f).Count() {
		t.Fail()
	}
}

func TestDictionaryCombine(t *testing.T) {
	var dict2 = dicts.NumberWithTrue.Select(func(k int, v bool) (int, bool) {
		return k + 1, !v
	})
	var dict = dicts.NumberWithTrue.Merge(dict2, func(old bool) bool { return old })
	if dict.Count(func(_ int, v bool) bool { return !v }) != 1 {
		t.Fatalf("there must be only 1 false value.")
	}
	dict = dicts.NumberWithTrue.Merge(dict2)
	if dict.Count(func(_ int, v bool) bool { return v }) != 1 {
		t.Fatalf("there must be only 1 true value.")
	}
	dict = dicts.NumberWithTrue.Merge(collections.From(map[int]bool{}).Dictionary())
	if !reflect.DeepEqual(dict.Map(), dicts.NumberWithTrue.Map()) {
		t.Fail()
	}
}

func TestDictionaryFunctionSignature(t *testing.T) {
	EstimateFail(t, func(*testing.T) {
		dicts.NumberWithTrue.Where(func() {})
	})
	EstimateFail(t, func(*testing.T) {
		dicts.NumberWithTrue.Select(func() {})
	})
	EstimateFail(t, func(*testing.T) {
		dicts.NumberWithTrue.ForEach(func(string) {})
	})
	EstimateFail(t, func(*testing.T) {
		dicts.NumberWithTrue.Merge(collections.From(map[bool]int{true: 1, false: 0}).Dictionary())
	})
	EstimateFail(t, func(*testing.T) {
		dicts.NumberWithTrue.Merge(collections.From(map[int]bool{}).Dictionary(), func(a, b int) int { return a })
	})
}
