package collections_test

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/johnwiichang/collections"
)

var slices = struct {
	Number collections.List
	Struct collections.List
}{
	collections.From([]int{1, 2, 3, 4, 5}).List(),
	collections.From([]int{1, 2, 3, 4, 5}).List().Select(func(n int) *Int { return &Int{n, []string{strconv.Itoa(n)}} }),
}

type Int struct {
	Value int
	Many  []string
}

func (i *Int) EqualsTo(n int) bool {
	return i.Value == n
}

func TestListBasicCopy(t *testing.T) {
	var slice, dstSlice = []int{1, 2, 3, 4, 5}, make([]int, 2)
	if reflect.DeepEqual(slice, slices.Number.Slice(&dstSlice)) {
		if reflect.DeepEqual(slice, dstSlice) {
			EstimateFail(t, func(*testing.T) {
				slices.Number.Slice(new([]string))
			})
		}
	}
}

func TestSliceSelect(t *testing.T) {
	var dst = []int{1, 3, 5, 7, 9}
	if !reflect.DeepEqual(dst, slices.Number.Select(func(i, n int) int { return i + n }).Slice()) {
		t.Fail()
	}
	var dst2 = slices.Number.Select(func(n int) string { return strconv.Itoa(n) }).Slice().([]string)
	if !reflect.DeepEqual(dst2, slices.Struct.SelectMany(func(i *Int) []string { return i.Many }).Slice()) {
		t.Fail()
	}
}

func TestSliceToMap(t *testing.T) {
	var map1 = map[string]int{"1": 1, "2": 2, "3": 3, "4": 4, "5": 5}
	var map2 = map[int]string{1: "1", 2: "2", 3: "3", 4: "4", 5: "5"}
	var test = slices.Number.ToDictionary(func(n int) (string, int) { return strconv.Itoa(n), n })
	if !reflect.DeepEqual(test.Map(), map1) {
		t.Fail()
	}
	test = slices.Number.ToDictionary(func(n int) string { return strconv.Itoa(n) })
	if !reflect.DeepEqual(test.Map(), map2) {
		t.Fail()
	}
	slices.Number.ToDictionary()
}

func TestSliceDistinct(t *testing.T) {
	var slice = collections.From([]int{1, 2, 3, 4, 5}).List()
	if !reflect.DeepEqual(slices.Number.Slice(), slice.Union(slices.Number).Slice()) {
		t.Fail()
	}
}

func TestSliceCount(t *testing.T) {
	var f = func(n int) bool { return n > 2 }
	if slices.Number.Where(f).Count() != slices.Number.Count(f) {
		t.Fail()
	}
}

func TestSliceContains(t *testing.T) {
	if slices.Number.Contains(collections.From([]string{"1", "2", "3"}).List()) || !slices.Number.Any(1, "2", 3) {
		t.Fail()
	}
	if !slices.Number.Contains(1, 2, 3) || !slices.Number.Contains() {
		t.Fail()
	}
	if !slices.Number.Any(5, 6) || !slices.Number.Any() {
		t.Fail()
	}
}

func TestSliceSort(t *testing.T) {
	var dst = collections.From([]float64{1.1, 2.1, 4.1, 3.1, 5.0}).List().Sort()
	var dstSlice = dst.Select(func(item float64) int { return int(item) }).Slice()
	if !reflect.DeepEqual(dstSlice, slices.Number.Slice()) {
		t.Fail()
	}
	dst = slices.Number.Select(func(n int) string { return strconv.Itoa(n) }).Sort(func(i, j string) bool {
		return i < j
	}).Reverse()
	if !reflect.DeepEqual([]string{"5", "4", "3", "2", "1"}, dst.Slice()) {
		t.Fail()
	}
	EstimateFail(t, func(*testing.T) {
		collections.From([]interface{}{"1", 2, 'r'}).List().Sort()
	})
	EstimateFail(t, func(*testing.T) {
		collections.From([]interface{}{"1", 2, 'r'}).List().Sort(func() {})
	})
}

func TestSliceFind(t *testing.T) {
	if slices.Number.First(1) != 0 {
		t.Fail()
	}
	if slices.Number.First(func(item int) bool {
		return strconv.Itoa(item) == "1"
	}) != 0 {
		t.Fail()
	}
	var slice = slices.Number.Concat(slices.Number)
	if slice.Last(1) != 5 {
		t.Fail()
	}
	if slice.Last(func(item int) bool {
		return strconv.Itoa(item) == "1"
	}) != 5 {
		t.Fail()
	}
	EstimateFail(t, func(*testing.T) {
		slices.Number.Concat(collections.From([]string{"1"}).List())
	})
}

func TestSliceIntersect(t *testing.T) {
	var slice = slices.Number.Intersect(collections.From([]int{0, 2, 3}).List())
	if !reflect.DeepEqual(slice.Slice(), []int{2, 3}) {
		t.Fail()
	}
}

func TestSliceExcept(t *testing.T) {
	var slice = slices.Number.Except(collections.From([]int{0, 2, 3}).List())
	if !reflect.DeepEqual(slice.Slice(), []int{1, 4, 5}) {
		t.Fail()
	}
	if slices.Struct.Except(slices.Number).Count() != 0 || slices.Number.Except(slices.Struct).Count() != 0 {
		t.Fail()
	}
}

func TestSliceFunctionSignature(t *testing.T) {
	EstimateFail(t, func(*testing.T) {
		slices.Number.Where(func() {})
	})
	EstimateFail(t, func(*testing.T) {
		slices.Number.First(func() {})
	})
	EstimateFail(t, func(*testing.T) {
		slices.Number.Select(func() {})
	})
	EstimateFail(t, func(*testing.T) {
		slices.Number.SelectMany(func() {})
	})
	EstimateFail(t, func(*testing.T) {
		slices.Number.ForEach(func(string) {})
	})
}

func TestSliceSkipAndTake(t *testing.T) {
	var slice = slices.Number.Skip(3).Take(1).Slice().([]int)
	if slice[0] != 4 {
		t.Fail()
	}
	if slices.Number.Skip(10).Take(1).Count() != 0 {
		t.Fail()
	}
	if slices.Number.Skip(4).Take(10).Count() != 1 {
		t.Fail()
	}
}

func TestSliceSkipAndResize(t *testing.T) {
	var slice = slices.Number.Skip(0).Take(5)
	if slice.Skip(5).Resize(10).Count() != 0 {
		t.Fail()
	}
	slice = slice.Concat(slices.Number)
	if !reflect.DeepEqual(slice.Skip(1).Resize(2).Slice(), []int{2, 3}) {
		t.Fail()
	}
	if !reflect.DeepEqual(slice.Skip(0).Resize(5).Slice(), []int{2, 3}) {
		t.Fail()
	}
	if !reflect.DeepEqual(slice.Skip(1).Resize().Slice(), []int{3}) {
		t.Fail()
	}
}
