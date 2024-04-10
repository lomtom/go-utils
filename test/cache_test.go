package test

import (
	"github.com/lomtom/go-utils/cache"

	"fmt"
	"runtime"
	"testing"
	"time"
)

func TestBackup(t *testing.T) {
	c, err := cache.NewMapCache[int](cache.SetExpirationTime(time.Minute), cache.SetGcInterval(time.Second*10), cache.SetEnablePersistence("test"), cache.SetPersistencePath("/tmp/cache/persistence"))
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	c.SetDefault("3", 3, time.Hour*2)
	c.Set("4", 4)
	go func() {
		for {
			time.Sleep(time.Second * 1)
			fmt.Println(c.GetWithExpiration("3"))
			fmt.Println(c.GetWithExpiration("4"))
		}
	}()
	fmt.Println(c.Get("3"))
	time.Sleep(time.Second * 3)
	fmt.Println(c.Get("3"))
}

func TestGc(t *testing.T) {
	go func() {
		for {
			number := runtime.NumGoroutine()
			fmt.Println("number:", number)
			time.Sleep(time.Second)
			runtime.GC()
		}
	}()
	index := 0
	for {
		_, _ = cache.NewMapCache[any](cache.SetExpirationTime(time.Millisecond), cache.SetGcInterval(time.Millisecond*10))
		time.Sleep(time.Second)
		index += 1
	}
}
