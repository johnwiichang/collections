# Collections

A Linq-like Extension for Collections. `Collectes` provides a collection of intuitive functions that is very similar to the LINQ extension method, which can avoid boring redundant code through `Collections`.

## List

List is suitable for Slice and Array types. But it is worth noting that it will be operated in  slice.

### Declare

```go
type Int struct {
	Value int
	Many  []string
}

var slices = struct {
	Number collections.List
	Struct collections.List
}{
	collections.From([]int{1, 2, 3, 4, 5}).List(),
	collections.From([]int{1, 2, 3, 4, 5}).List().Select(func(n int) *Int { return &Int{n, []string{strconv.Itoa(n)}} }),
}
```

*The data type within the List has no limitations.*

### Actions
Some mature methods are listed below for collection operations.

**Slice(slice ...interface{}) interface{}**

Get the internal slice as `interface{}`, user must know the specific type of internal slices or use assertions.

**Select(f interface{}) List**

Create another collection of elements from a collection of elements and output a new collection.

```go
slices.Numbers.Select(func(i, n int) int { return i + n })
```

> `i` is index of the `List` and `n` is the element of the `List`. **Ignore `i` directly if the serial number is not required.**

**SelectMany(f interface{}) List**

Creating another element collection and outputting a new collection from a collection of elements in a collection.

```go
slices.Struct.SelectMany(func(i *Int) []string { return i.Many })
```

> `i` is index of the `List` and `n` is the element of the `List`. **Ignore `i` directly if the serial number is not required.**

**ForEach(f interface{}) List**

Traverse the collection and then invoke a custom function (which will not change the element), support the first parameter to use the `bool` value `false` or the last parameter `error` as not `nil` to terminate traversal.

```go
slices.Struct.ForEach(func(i int, item *Int) (bool) {
	return len(item.Many) != 0
})
```

> The function will not change the elements within any collection!

**ToDictionary(f ...interface{}) Dictionary**

Map the elements in the List collection to Dictionary.

```go
slices.Number.ToDictionary(func(n int) (string, int) { return strconv.Itoa(n), n })
```

> When there is only one input, the default is the elements of the List collection. If indexing is required, specify the first input as int and the second input as the element type of the List.

> When there is only one return value, the key of the dictionary corresponds to the element in the List collection, while when there are two return values, the first value returned will be used as the key and the second value as the value.

**Sort(less ...interface{}) List**

For sorting. A comparator is supported to return whether the `i`-th element is smaller than the `j`-th element.

```go
slices.Number.Select(func(n int) string { return strconv.Itoa(n) }).Sort(func(i, j string) bool {
	return i < j
})
```

> If not specified, then the system's comparison function is used by default. You can refer to the use of the `sort` package.

**Reverse() List**

Invert the collection.

**Distinct() List**

Removes duplicate elements from a List collection.

**Where(f interface{}) List**

Query a List collection.

```go
var f = func(n int) bool { return n > 2 }
slices.Number.Where(f)
```

**Count(f ...interface{}) int**

Counting a List collection can also be conditionally equivalent to `.Where(f).Count()`.

```go
var f = func(n int) bool { return n > 2 }
slices.Number.Count(f)
// slices.Number.Where(f).Count()
```

**Contains(elements ...interface{}) bool**

Determines if the set contains all elements.

```go
slices.Number.Contains(1, 2, 3)
```

**Any(elements ...interface{}) bool**

Determines if the set contains any elements.

```go
slices.Number.Any(5, 6)
```

**Concat(l List) List**

Joins two List collections and returns a new List collection.

```go
slices.Number.Concat(slices.Number)
```

**First(obj interface{}) int**

Returns the position of the first occurrence of the target element.

```go
slices.Number.First(1)
```

**Last(obj interface{}) int**

Returns the position of the last occurrence of the target element.

```go
slices.Number.Last(1)
```

**Intersect(l List) List**

Gets the intersection of a List set with another List set.

```go
slices.Number.Intersect(collections.From([]int{0, 2, 3})
```

**Except(l List) List**

Exclude another List collection from a List collection.

```go
slices.Number.Except(collections.From([]int{0, 2, 3})
```

**Union(l List) List**

Similar to `Concat`, but with final de-duplication.

> Equivalent to `.Concat(l).Distinct()`

**Skip(length int) List**

Sets the position of the cursor inside the List.

**Take(num int) List**

Get the next number of elements from the internal cursor and return it as a list.

> Automatically stops when there are not enough elements.

**Resize(length ...int) List**

Similar to `Take`, but `Resize` performs an intercept operation on the current list.

## Dictionary

Dictionary is suitable for `map`.

### Declare


```go
var dicts = struct {
	NumberWithTrue collections.Dictionary
}{
	collections.From(map[int]bool{1: true, 2: true, 3: true, 4: true, 5: true}).Dictionary(),
}
```

*The data type within the Dictionary must be a valid `map`.*

### Actions
Some mature methods are listed below for collection operations.

**Map(m ...interface{}) interface{}**

Get the internal map as `interface{}`, user must know the specific type of internal maps or use assertions.

**Keys() List**

Get keys of the Dictionary collection as a list collection.

**Values() List**

Get values of the Dictionary collection as a list collection.

**Where(f interface{}) Dictionary**

Query a Dictionary collection.

```go
var f = func(n int) bool { return n > 2 }
slices.Number.Where(f)
```

**Count(f ...interface{}) int**

Counting a Dictionary collection can also be conditionally equivalent to `.Where(f).Count()`.

**ForEach(f interface{}) Dictionary**

Traverse the collection and then invoke a custom function (which will not change the element), support the first parameter to use the `bool` value `false` or the last parameter `error` as not `nil` to terminate traversal.

```go
slices.Struct.ForEach(func(i int, item *Int) (bool) {
	return len(item.Many) != 0
})
```

> The function will not change the elements within any collection!

**Select(f interface{}) Dictionary**

Create another collection of elements from a collection of elements and output a new collection.

```go
dicts.NumberWithTrue.Select(func(k int, v bool) (int, string) {
	return k, strconv.Itoa(k)
})
```

> `k` is key of the `Dictionary` and `v` is the value of the `Dictionary`. **Ignore `v` directly if the value is not required.**

**Merge(d Dictionary, onConflict ...interface{}) Dictionary**

Combine two dictionaries into a `Dictionary`.

```go
var dict2 = dicts.NumberWithTrue.Select(func(k int, v bool) (int, bool) {
	return k + 1, !v
})
var dict = dicts.NumberWithTrue.Merge(dict2, func(old bool) bool { return old })
```

You can make your decisions when conflicting keys are encountered. The conflicting keys are listed in *'old' - 'new'* order and will be overwritten by default using the merged target dictionary values.

> If a new value is not required, it can be ignored directly in the parameters as in the example code.