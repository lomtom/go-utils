package assert

import (
	"fmt"
	"reflect"
	"runtime"
	"testing"
)

const (
	Equal int = iota
	NotEqual
	Less
	More
)

type Interface interface {
	Equal(expected, actual any)
}

type Assert struct {
	testing.TB
}

func NewAssert(t testing.TB) Interface {
	return &Assert{t}
}

func (a *Assert) Equal(expected, actual any) {
	if compare(expected, actual) != Equal {
		printError(a, expected, actual)
	}
}

func compare(elem1, elem2 any) int {
	if reflect.DeepEqual(elem1, elem2) {
		return Equal
	}
	return NotEqual
}

func printError(t testing.TB, expected, actual any) {
	_, file, line, _ := runtime.Caller(2)
	errInfo := fmt.Sprintf("failed : %v, line: %v, expected: %v, actual: %v.", file, line, expected, actual)
	t.Error(errInfo)
	t.FailNow()
}
