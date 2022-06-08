package cache

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type Map struct {
	items  map[string]item // Cache data items are stored in the map
	mu     sync.RWMutex    // Read write lock
	stopGc chan bool
	isGc   bool
	options
}

// NewMapCache create a cache with Map
func NewMapCache(opts ...CreateOptionFunc) MapInterface {
	exp := newOption()
	for _, opt := range opts {
		opt(&exp)
	}
	res := &Map{
		options: exp,
	}
	if exp.expiration != DefaultExpiration {
		// 开启gc
		_ = res.StartGc()
	}
	return res
}

// Expired cache data item cleanup
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

// StopGc stop gc
func (c *Map) StopGc() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.isGc {
		return errors.New("GC is closed")
	}
	c.isGc = false
	c.stopGc <- true
	return nil
}

// StartGc start gc
// After the expiration time is set, GC will be started automatically without manual GC
func (c *Map) StartGc() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.isGc {
		return errors.New("GC has been started")
	}
	c.isGc = true
	go c.gcLoop()
	return nil
}

// delete data by key
func (c *Map) del(key string) {
	delete(c.items, key)
}

// set cache data by key
func (c *Map) set(key string, value interface{}, expiration int64) {
	c.items[key] = item{
		value,
		expiration,
	}
}

// get data by key
func (c *Map) get(key string) (*item, bool) {
	value, ok := c.items[key]
	if !ok || value.expired() {
		return nil, false
	}
	return &value, true
}

// generate expiration time
func (c *Map) generateExpiration() int64 {
	if c.expiration == DefaultExpiration {
		return 0
	}
	return time.Now().Add(c.expiration).UnixNano() / 1e3
}

// init data
func (c *Map) judgeAndInitItem() {
	if c.items == nil {
		c.items = make(map[string]item)
	}
}

// IsExpired judge whether the data is expired
func (c *Map) IsExpired(key string) (bool, error) {
	value, ok := c.items[key]
	if !ok {
		return false, fmt.Errorf("the data %s does not exist", key)
	}
	return value.expired(), nil
}

// DeleteExpired delete all expired data
func (c *Map) DeleteExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for k, v := range c.items {
		if v.expired() {
			c.del(k)
		}
	}
}

// Delete delete data by key
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

// Set  data by key，it will overwrite the data if the key exists
func (c *Map) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.judgeAndInitItem()

	c.set(key, value, c.generateExpiration())
}

// Add data，Cannot add existing data
// To override the addition, use the set method
func (c *Map) Add(key string, value interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.judgeAndInitItem()
	if _, ok := c.items[key]; ok {
		return fmt.Errorf("data %s already exists", key)
	}

	c.set(key, value, c.generateExpiration())
	return nil
}

// Get  data
// When the data does not exist or expires, it will return nonexistence（false）
func (c *Map) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	value, ok := c.items[key]
	if !ok || value.expired() {
		return nil, false
	}
	return value.object, true
}

// GetAndDelete get data and delete by key
func (c *Map) GetAndDelete(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	value, ok := c.items[key]
	if !ok || value.expired() {
		return nil, false
	}
	// delete
	c.del(key)
	return value.object, true
}

// GetAndExpired  get data and expire by key
// It will be deleted at the next clearing. If the clearing capability is not enabled, it will never be deleted
func (c *Map) GetAndExpired(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	value, ok := c.items[key]
	if !ok || value.expired() {
		return nil, false
	}
	// Set expiration
	c.set(key, value, time.Now().UnixNano()/1e3)
	return value.object, true
}

// Clear remove all data
func (c *Map) Clear() {
	c.items = make(map[string]item)
}

// Keys get all keys
func (c *Map) Keys() []string {
	res := make([]string, 0)
	for k, _ := range c.items {
		res = append(res, k)
	}
	return res
}
