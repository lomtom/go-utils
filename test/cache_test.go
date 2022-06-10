package test

import (
	"fmt"
	"github.com/lomtom/go-utils/cache"
	"testing"
	"time"
)

func TestBackup(t *testing.T) {
	c, err := cache.NewMapCache(cache.SetExpirationTime(time.Minute), cache.SetGcInterval(time.Second*10), cache.SetEnablePersistence("test"), cache.SetPersistencePath("/tmp/cache/persistence"))
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	c.Set("3", 1)
	time.Sleep(time.Second * 5)
	c.Set("2", 4)
	time.Sleep(time.Second * 5)
	fmt.Println(c.Get("1"))
	time.Sleep(time.Second * 5)
}
