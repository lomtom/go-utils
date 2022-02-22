package test

import (
	"github.com/zero028/go-utils/cache"
	"log"
	"testing"
)

func TestMapCache(t *testing.T) {
	c := cache.NewMapCache()
	c.Set("1", 1)
	log.Println(c.Get("1"))
}
