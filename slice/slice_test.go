package slice

import (
	"github.com/lomtom/go-utils/v2/assert"
	"testing"
)

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

func TestRetainAll(t *testing.T) {
	a := assert.NewAssert(t)
	a.Equal([]int{1, 2, 3}, RetainAll([]int{1, 2, 3, 4}, []int{1, 2, 3}...))
	a.Equal([]int{2, 3}, RetainAll([]int{1, 2, 3}, []int{2, 3, 4}...))
	a.Equal([]people{}, RetainAll([]people{}, people{name: "lomtom"}))
}

func TestFilter(t *testing.T) {
	a := assert.NewAssert(t)
	a.Equal([]int{3, 4}, Filter([]int{1, 2, 3, 4}, func(value int) bool {
		return value > 2
	}))
}

func TestDeleteAt(t *testing.T) {
	a := assert.NewAssert(t)

	a.Equal([]string{"a", "b", "c"}, DeleteAt([]string{"a", "b", "c"}, -1))
	a.Equal([]string{"a", "b", "c"}, DeleteAt([]string{"a", "b", "c"}, 3))
	a.Equal([]string{}, DeleteAt([]string{"a", "b", "c"}, 0))
	a.Equal([]string{"a"}, DeleteAt([]string{"a", "b", "c"}, 1))
	a.Equal([]string{"a", "b"}, DeleteAt([]string{"a", "b", "c"}, 2))

	a.Equal([]string{"b", "c"}, DeleteAt([]string{"a", "b", "c"}, 0, 1))
	a.Equal([]string{"c"}, DeleteAt([]string{"a", "b", "c"}, 0, 2))
	a.Equal([]string{}, DeleteAt([]string{"a", "b", "c"}, 0, 3))
	a.Equal([]string{}, DeleteAt([]string{"a", "b", "c"}, 0, 4))
	a.Equal([]string{"a"}, DeleteAt([]string{"a", "b", "c"}, 1, 3))
	a.Equal([]string{"a"}, DeleteAt([]string{"a", "b", "c"}, 1, 4))

}

func TestInsertAt(t *testing.T) {
	a := assert.NewAssert(t)

	a.Equal([]string{"a", "b", "c"}, InsertAt([]string{"a", "b", "c"}, -1, "1"))
	a.Equal([]string{"a", "b", "c"}, InsertAt([]string{"a", "b", "c"}, 4, "1"))
	a.Equal([]string{"1", "a", "b", "c"}, InsertAt([]string{"a", "b", "c"}, 0, "1"))
	a.Equal([]string{"a", "1", "b", "c"}, InsertAt([]string{"a", "b", "c"}, 1, "1"))
	a.Equal([]string{"a", "b", "1", "c"}, InsertAt([]string{"a", "b", "c"}, 2, "1"))
	a.Equal([]string{"a", "b", "c", "1"}, InsertAt([]string{"a", "b", "c"}, 3, "1"))
	a.Equal([]string{"1", "2", "3", "a", "b", "c"}, InsertAt([]string{"a", "b", "c"}, 0, "1", "2", "3"))
	a.Equal([]string{"a", "1", "2", "3", "b", "c"}, InsertAt([]string{"a", "b", "c"}, 1, "1", "2", "3"))
	a.Equal([]string{"a", "b", "c", "1", "2", "3"}, InsertAt([]string{"a", "b", "c"}, 3, "1", "2", "3"))
}

func TestDistinct(t *testing.T) {
	a := assert.NewAssert(t)

	a.Equal([]string{}, Distinct([]string{}))
	a.Equal([]string{"a", "b", "c"}, Distinct([]string{"a", "b", "c"}))
	a.Equal([]string{"a", "b", "c"}, Distinct([]string{"a", "b", "b", "c"}))
	a.Equal([]string{"a", "b", "c"}, Distinct([]string{"a", "b", "b", "c", "c"}))
	a.Equal([]int{1, 2, 3}, Distinct([]int{1, 2, 3, 2, 3, 2, 3}))
}

func TestSort(t *testing.T) {
	a := assert.NewAssert(t)
	ints := []int{2, 3, 4, 1}
	Sort(ints, func(value1 int, value2 int) bool {
		return value1 < value2
	})
	a.Equal([]int{1, 2, 3, 4}, ints)

	peoples := []people{{"lomtom", 20}, {"lomtom", 18}, {"lomtom", 19}}
	Sort(peoples, func(value1, value2 people) bool {
		return value1.age < value2.age
	})
	a.Equal([]people{{"lomtom", 18}, {"lomtom", 19}, {"lomtom", 20}}, peoples)
}

func TestSortByField(t *testing.T) {
	a := assert.NewAssert(t)
	peoples := []people{{"lomtom", 20}, {"lomtom", 18}, {"lomtom", 19}}
	a.Equal(nil, SortByField(peoples, "age"))
	a.Equal([]people{{"lomtom", 18}, {"lomtom", 19}, {"lomtom", 20}}, peoples)
}
