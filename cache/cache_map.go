package cache

import (
	"fmt"
	"sync"
	"time"
)

/******************************************* 数据 *******************************************/

type item struct {
	object     interface{} // 数据
	expiration int64       // 过期时间
}

// Expired 判断数据项是否已经过期
func (item *item) expired() bool {
	return time.Now().UnixMicro() > item.expiration
}

// 设置过期，将会在下个清理缓存的周期进行清除
func (item *item) setExpired() {
	item.expiration = time.Now().UnixMicro()
}

/******************************************* 过期策略 *******************************************/

const (
	// DefaultExpiration 默认过期时间标志 永不过期
	DefaultExpiration time.Duration = -1

	// DefaultInterval 默认过期间隔 一分钟
	DefaultInterval = time.Minute
)

type expirationPolicy struct {
	expiration time.Duration // 过期时间
	gcInterval time.Duration // 过期数据项清理周期
}

func newExpirationPolicy() expirationPolicy {
	return expirationPolicy{
		expiration: DefaultExpiration,
		gcInterval: DefaultInterval,
	}
}

// CreateOptionFunc 初始化可选参数（定义过期策略）
type CreateOptionFunc func(o *expirationPolicy)

// SetExpirationTime  设置过期时间
// expiration 过期时间
func SetExpirationTime(expiration time.Duration) CreateOptionFunc {
	return func(o *expirationPolicy) {
		o.expiration = expiration
	}
}

// SetGcInterval 设置gc间隔
// gcInterval 清理周期 周期为0时，自动改为1分钟
func SetGcInterval(gcInterval time.Duration) CreateOptionFunc {
	if gcInterval == 0 {
		gcInterval = time.Minute
	}
	return func(o *expirationPolicy) {
		o.gcInterval = gcInterval
	}
}

/******************************************* 缓存 *******************************************/

type Cache struct {
	expiration time.Duration   // 过期时间
	items      map[string]item // 缓存数据项存储在 map 中
	mu         sync.RWMutex    // 读写锁
	gcInterval time.Duration   // 过期数据项清理周期
	stopGc     chan bool
}

// NewCache 新建缓存
func NewCache(opts ...CreateOptionFunc) *Cache {
	exp := newExpirationPolicy()
	for _, opt := range opts {
		opt(&exp)
	}
	res := &Cache{
		expiration: exp.expiration,
		gcInterval: exp.gcInterval,
	}
	if exp.expiration != DefaultExpiration {
		// 开启gc
		go res.gcLoop()
	}
	return res
}

// 过期缓存数据项清理
func (c *Cache) gcLoop() {
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
func (c *Cache) StopGc() {
	c.stopGc <- true
}

// 删除缓存数据项
func (c *Cache) del(key string) {
	delete(c.items, key)
}

// 设置缓存数据
func (c *Cache) set(key string, value interface{}, expiration int64) {
	c.items[key] = item{
		value,
		expiration,
	}
}

func (c *Cache) get(key string) (*item, bool) {
	value, ok := c.items[key]
	if !ok || value.expired() {
		return nil, false
	}
	return &value, true
}

// 初始化数据
func (c *Cache) judgeAndInitItem() {
	if c.items == nil {
		c.items = make(map[string]item)
	}
}

// IsExpired 判断是否过期
func (c *Cache) IsExpired(key string) (bool, error) {
	value, ok := c.items[key]
	if !ok {
		return false, fmt.Errorf("该数据不存在")
	}
	return value.expired(), nil
}

// DeleteExpired 删除过期数据项
func (c *Cache) DeleteExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for k, v := range c.items {
		if v.expired() {
			c.del(k)
		}
	}
}

func (c *Cache) Delete(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	value, ok := c.get(key)
	if ok {
		c.del(key)
		return value.object, ok
	}
	return nil, ok
}

func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.judgeAndInitItem()
	c.set(key, value, time.Now().Add(c.expiration).UnixMicro())
}

// Add 如需覆盖添加，请使用Set方法
func (c *Cache) Add(key string, value interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.judgeAndInitItem()
	if _, ok := c.items[key]; ok {
		return fmt.Errorf("数据 %s 已存在", key)
	}
	c.set(key, value, time.Now().Add(c.expiration).UnixMicro())
	return nil
}

// Get 不存在或过期都会返回不存在
// 返回数据、是否存在
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	value, ok := c.items[key]
	if !ok || value.expired() {
		return nil, false
	}
	return value.object, true
}

func (c *Cache) GetAndDelete(key string) (interface{}, bool) {
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

func (c *Cache) GetAndExpired(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	value, ok := c.items[key]
	if !ok || value.expired() {
		return nil, false
	}
	// 过期
	c.set(key, value, time.Now().UnixMicro())
	return value.object, true
}
