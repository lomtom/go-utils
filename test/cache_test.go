package test

import (
	"fmt"
	"github.com/lomtom/go-utils/cache"
	"testing"
)

func TestMapCache(t *testing.T) {
	c := cache.NewMapCache()
	c.Set("1", 1)
	fmt.Println(c.Get("1"))
	fmt.Println(c.Keys())
}
