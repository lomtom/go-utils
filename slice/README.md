介绍
---
为切片工具类，提供一些常用方法

接口
---
```go
// AddAll add all elements of values to slice
AddAll[T any](slice []T, values ...T) []T

// Contains judge whether slice contain value
Contains[T any](slice []T, value T) bool

// ContainsAll judge whether slice contain all of values
// If one of values does not contain, false will be returned
ContainsAll[T any](slice []T, values ...T) bool

// DeleteAt delete elements of slice from from index to to - 1 index
DeleteAt[T any](slice []T, from int, to ...int) []T

// Distinct return the distinct elements of slice
Distinct[T any](slice []T) []T

// Filter it will filter elements of slice by predicate func
Filter[T any](slice []T, predicate Predicate[T]) []T

// InsertAt it will insert elem from values in index
InsertAt[T any](slice []T, index int, values ...T) []T

// IsEmpty judge whether the slice is empty
IsEmpty[T any](slice []T) bool

// RemoveAll remove elements from slice
// 1,2,3  2,3 -> 1
RemoveAll[T any](slice []T, values ...T) []T

// ReplaceAll replace elements from slice by operator func
ReplaceAll[T any](slice []T, operator UnaryOperator[T]) []T

// RetainAll
// 1,2,3   2,3,4  -> 2,3
RetainAll[T any](slice []T, values ...T) []T

Size[T any](slice []T) int

// Sort It will be sorted by comparator
Sort[T any](slice []T, comparator Comparator[T])

// SortByField It will be sorted by field,default asc
SortByField[T any](slice []T, field string, sortType ...Collation) error
```


使用
---
**示例**
点击[示例](slice_test.go)查看完整使用
```go
type people struct {
	name string
	age  int
}

func TestSize(t *testing.T) {
	a := assert.NewAssert(t)
	a.Equal(0, Size([]int{}))
	a.Equal(3, Size([]int{1, 2, 3}))
}

func TestIsEmpty(t *testing.T) {
	a := assert.NewAssert(t)
	a.Equal(true, IsEmpty([]int{}))
	a.Equal(false, IsEmpty([]int{1, 2, 3}))
}

func TestContains(t *testing.T) {
	a := assert.NewAssert(t)
	a.Equal(true, Contains([]int{1, 2, 3}, 1))
	a.Equal(true, Contains([]people{}, people{name: "lomtom"}))
}

func TestContainsAll(t *testing.T) {
	a := assert.NewAssert(t)
	a.Equal(true, ContainsAll([]int{1, 2, 3}))
	a.Equal(true, ContainsAll([]int{1, 2, 3}, 1, 2, 3))
	a.Equal(false, ContainsAll([]int{1, 2, 3}, 1, 2, 3, 4))
	a.Equal(true, ContainsAll([]people{{name: "lomtom"}}, people{name: "lomtom"}))
}

func TestAddAll(t *testing.T) {
	a := assert.NewAssert(t)
	a.Equal([]int{1, 2, 3, 1, 2, 3}, AddAll([]int{1, 2, 3}, []int{1, 2, 3}...))
	a.Equal([]people{{name: "lomtom"}, {name: "lomtom"}}, AddAll([]people{{name: "lomtom"}}, people{name: "lomtom"}))
}

func TestReplaceAll(t *testing.T) {
	a := assert.NewAssert(t)
	a.Equal([]int{4, 4, 3, 1, 4, 3}, ReplaceAll([]int{1, 2, 3, 1, 2, 3}, func(value int) int {
		if value == 2 {
			return 4
		}
		return value
	}))
}

func TestRemoveAll(t *testing.T) {
	a := assert.NewAssert(t)
	a.Equal([]int{4}, RemoveAll([]int{1, 2, 3, 4}, []int{1, 2, 3}...))
	a.Equal([]people{}, RemoveAll([]people{{name: "lomtom"}}, people{name: "lomtom"}))
}
```