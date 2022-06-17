package slice

import (
	"fmt"
	"reflect"
	"sort"
)

type Collation int

const (
	Asc Collation = iota
	Desc
)

type Predicate[E any] func(value E) bool

type UnaryOperator[E any] func(value E) E

type Comparator[E any] func(value1 E, value2 E) bool

// AddAll add all elements of values to slice
func AddAll[E any](slice []E, values ...E) []E {
	return append(slice, values...)
}

// Contains judge whether slice contain value
func Contains[E any](slice []E, value E) bool {
	for _, v := range slice {
		if reflect.DeepEqual(v, value) {
			return true
		}
	}
	return false
}

// ContainsAll judge whether slice contain all of values
// If one of values does not contain, false will be returned
func ContainsAll[E any](slice []E, values ...E) bool {
	for _, value := range values {
		if !Contains(slice, value) {
			return false
		}
	}
	return true
}

// DeleteAt delete elements of slice from from index to to - 1 index
func DeleteAt[E any](slice []E, from int, to ...int) []E {
	size := len(slice)
	if from < 0 || from >= size {
		return slice
	}
	// delete elements from from index to to[0]
	if len(to) > 0 {
		end := to[0]
		if end <= from {
			return slice
		}
		if end > size {
			end = size
		}
		slice = append(slice[:from], slice[end:]...)
		return slice
	}
	return slice[:from]
}

// Distinct return the distinct elements of slice
func Distinct[E any](slice []E) []E {
	if len(slice) == 0 {
		return []E{}
	}
	var result []E
	// bubble sort
	for i := 0; i < len(slice); i++ {
		value := slice[i]
		skip := true
		for j := range result {
			if reflect.DeepEqual(value, result[j]) {
				skip = false
				break
			}
		}
		if skip {
			result = append(result, value)
		}
	}
	return result
}

// Filter it will filter elements of slice by predicate func
func Filter[E any](slice []E, predicate Predicate[E]) []E {
	if predicate == nil {
		return slice
	}
	result := make([]E, 0, 0)
	for _, value := range slice {
		b := predicate(value)
		if b {
			result = append(result, value)
		}
	}
	return result
}

// InsertAt it will insert elem from values in index
func InsertAt[E any](slice []E, index int, values ...E) []E {
	size := len(slice)
	if index < 0 || index > size {
		return slice
	}
	slice = append(slice[:index], append(values, slice[index:]...)...)
	return slice
}

// IsEmpty judge whether the slice is empty
func IsEmpty[E any](slice []E) bool {
	return len(slice) == 0
}

// RemoveAll remove elements from slice
// 1,2,3  2,3 -> 1
func RemoveAll[E any](slice []E, values ...E) []E {
	result := make([]E, 0)

	for _, value := range slice {
		if !Contains(values, value) {
			result = append(result, value)
		}
	}
	return result
}

// ReplaceAll replace elements from slice by operator func
func ReplaceAll[E any](slice []E, operator UnaryOperator[E]) []E {
	for index, value := range slice {
		slice[index] = operator(value)
	}
	return slice
}

// RetainAll
// 1,2,3   2,3,4  -> 2,3
func RetainAll[E any](slice []E, values ...E) []E {
	result := make([]E, 0)

	for _, value := range slice {
		if Contains(values, value) {
			result = append(result, value)
		}
	}
	return result
}

func Size[E any](slice []E) int {
	return len(slice)
}

// Sort It will be sorted by comparator
func Sort[E any](slice []E, comparator Comparator[E]) {
	sort.Slice(slice, func(index1, index2 int) bool {
		return comparator(slice[index1], slice[index2])
	})
}

// SortByField It will be sorted by field,default asc
func SortByField[E any](slice []E, field string, sortType ...Collation) error {
	sv := reflect.ValueOf(slice)
	t := sv.Type().Elem()
	// Find the field.
	sf, ok := t.FieldByName(field)
	if !ok {
		return fmt.Errorf("field name %s not found", field)
	}
	var compare func(a, b reflect.Value) bool
	switch sf.Type.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if len(sortType) > 0 && sortType[0] == Desc {
			compare = func(a, b reflect.Value) bool { return a.Int() > b.Int() }
		} else {
			compare = func(a, b reflect.Value) bool { return a.Int() < b.Int() }
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if len(sortType) > 0 && sortType[0] == Desc {
			compare = func(a, b reflect.Value) bool { return a.Uint() > b.Uint() }
		} else {
			compare = func(a, b reflect.Value) bool { return a.Uint() < b.Uint() }
		}
	case reflect.Float32, reflect.Float64:
		if len(sortType) > 0 && sortType[0] == Desc {
			compare = func(a, b reflect.Value) bool { return a.Float() > b.Float() }
		} else {
			compare = func(a, b reflect.Value) bool { return a.Float() < b.Float() }
		}
	case reflect.String:
		if len(sortType) > 0 && sortType[0] == Desc {
			compare = func(a, b reflect.Value) bool { return a.String() > b.String() }
		} else {
			compare = func(a, b reflect.Value) bool { return a.String() < b.String() }
		}
	case reflect.Bool:
		if len(sortType) > 0 && sortType[0] == Desc {
			compare = func(a, b reflect.Value) bool { return a.Bool() && !b.Bool() }
		} else {
			compare = func(a, b reflect.Value) bool { return !a.Bool() && b.Bool() }
		}
	default:
		return fmt.Errorf("field type %s not supported", sf.Type)
	}
	sort.Slice(slice, func(i, j int) bool {
		a := sv.Index(i)
		b := sv.Index(j)
		if t.Kind() == reflect.Ptr {
			a = a.Elem()
			b = b.Elem()
		}
		a = a.FieldByIndex(sf.Index)
		b = b.FieldByIndex(sf.Index)
		return compare(a, b)
	})

	return nil
}
