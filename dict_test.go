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
}
