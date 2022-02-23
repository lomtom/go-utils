package cache

import (
	"log"
	"testing"
)

func TestMapCache(t *testing.T) {
	c := NewMapCache()
	c.Set("1", 1)
	log.Println(c.Get("1"))
}
