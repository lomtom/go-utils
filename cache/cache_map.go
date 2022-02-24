package cache

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

/******************************************* 缓存 *******************************************/

type Map struct {
	items      map[string]item // 缓存数据项存储在 map 中
	expiration time.Duration   // 过期时间
	mu         sync.RWMutex    // 读写锁
	gcInterval time.Duration   // 过期数据项清理周期
	stopGc     chan bool
	isGc       bool
}

// NewMapCache 新建缓存
func NewMapCache(opts ...CreateOptionFunc) MapInterface {
	exp := newExpirationPolicy()
	for _, opt := range opts {
		opt(&exp)
	}
	res := &Map{
		expiration: exp.expiration,
		gcInterval: exp.gcInterval,
	}
	if exp.expiration != DefaultExpiration {
		// 开启gc
		_ = res.StartGc()
	}
	return res
}

// 过期缓存数据项清理
func (c *Map) gcLoop() {
	ticker := time.NewTicker(c.gcInterval)
	for {
		select {
		case <-ticker.C:
			c.DeleteExpired()
		case <-c.stopGc:
			ticker.Stop()
			return
		}
	}
}

// StopGc 停止gc
func (c *Map) StopGc() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.isGc {
		return errors.New("GC程序已关闭")
	}
	c.isGc = false
	c.stopGc <- true
	return nil
}

// StartGc 重新gc
// 设置过期时间时，会自动开启gc，无需手动gc
func (c *Map) StartGc() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.isGc {
		return errors.New("GC程序已开启")
	}
	c.isGc = true
	go c.gcLoop()
	return nil
}

// 删除缓存数据项
func (c *Map) del(key string) {
	delete(c.items, key)
}

// 设置缓存数据
func (c *Map) set(key string, value interface{}, expiration int64) {
	c.items[key] = item{
		value,
		expiration,
	}
}

func (c *Map) get(key string) (*item, bool) {
	value, ok := c.items[key]
	if !ok || value.expired() {
		return nil, false
	}
	return &value, true
}

// 生成过期时间
func (c *Map) generateExpiration() int64 {
	if c.expiration == DefaultExpiration {
		return 0
	}
	return time.Now().Add(c.expiration).UnixNano() / 1e3
}

// 初始化数据
func (c *Map) judgeAndInitItem() {
	if c.items == nil {
		c.items = make(map[string]item)
	}
}

// IsExpired 判断是否过期
func (c *Map) IsExpired(key string) (bool, error) {
	value, ok := c.items[key]
	if !ok {
		return false, fmt.Errorf("该数据不存在")
	}
	return value.expired(), nil
}

// DeleteExpired 删除过期数据项
func (c *Map) DeleteExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for k, v := range c.items {
		if v.expired() {
			c.del(k)
		}
	}
}

func (c *Map) Delete(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	value, ok := c.get(key)
	if ok {
		c.del(key)
		return value.object, ok
	}
	return nil, ok
}

func (c *Map) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.judgeAndInitItem()

	c.set(key, value, c.generateExpiration())
}

// Add 如需覆盖添加，请使用Set方法
func (c *Map) Add(key string, value interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.judgeAndInitItem()
	if _, ok := c.items[key]; ok {
		return fmt.Errorf("数据 %s 已存在", key)
	}

	c.set(key, value, c.generateExpiration())
	return nil
}

// Get 不存在或过期都会返回不存在
// 返回数据、是否存在
func (c *Map) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	value, ok := c.items[key]
	if !ok || value.expired() {
		return nil, false
	}
	return value.object, true
}

func (c *Map) GetAndDelete(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	value, ok := c.items[key]
	if !ok || value.expired() {
		return nil, false
	}
	// 删除
	c.del(key)
	return value.object, true
}

func (c *Map) GetAndExpired(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	value, ok := c.items[key]
	if !ok || value.expired() {
		return nil, false
	}
	// 过期
	c.set(key, value, time.Now().UnixNano()/1e3)
	return value.object, true
}

func (c *Map) Clear() {
	c.items = make(map[string]item)
}
