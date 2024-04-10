package cache

import (
	"errors"
	"fmt"
	"runtime"
	"sync"
	"time"
)

type MapCache[E any] struct {
	*mapCache[E]
}

type mapCache[E any] struct {
	items  map[string]*Item[E] // Cache data items are stored in the map
	mu     sync.RWMutex        // Read write lock
	stopGc chan bool
	isGc   bool
	options
}

// NewMapCache create a cache with mapCache
func NewMapCache[E any](opts ...CreateOptionFunc) (MapInterface[E], error) {
	exp := newOption()
	for _, opt := range opts {
		opt(&exp)
	}
	res := &mapCache[E]{
		options: exp,
		stopGc:  make(chan bool),
	}
	if exp.expiration != DefaultExpiration {
		// start gc
		_ = res.StartGc()
	}
	if exp.enablePersistence {
		res.items = make(map[string]*Item[E])
		err := res.startPersistence(&(res.items))
		if err != nil {
			return nil, err
		}
	}
	c := &MapCache[E]{
		res,
	}
	runtime.SetFinalizer(c, func(m *MapCache[E]) {
		_ = m.StopGc()
	})
	return c, nil
}

// Expired cache data Item cleanup
func (c *mapCache[E]) gcLoop() {
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
func (c *mapCache[E]) StopGc() error {
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
func (c *mapCache[E]) StartGc() error {
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
func (c *mapCache[E]) del(key string) {
	delete(c.items, key)
}

// set cache data by key
func (c *mapCache[E]) set(key string, value E, expiration int64) {
	c.items[key] = &Item[E]{
		Object:     value,
		Expiration: expiration,
	}
}

// get data by key
func (c *mapCache[E]) get(key string) (*Item[E], bool) {
	value, ok := c.items[key]
	if !ok || value.expired() {
		return nil, false
	}
	return value, true
}

// generate expiration time
func (c *mapCache[E]) generateExpiration() int64 {
	if c.expiration == DefaultExpiration {
		return 0
	}
	return time.Now().Add(c.expiration).UnixNano() / 1e3
}

// generate expiration time
func (c *mapCache[E]) generateExpirationForItem(expiration time.Duration) int64 {
	return time.Now().Add(expiration).UnixNano() / 1e3
}

// init data
func (c *mapCache[E]) judgeAndInitItem() {
	if c.items == nil {
		c.items = make(map[string]*Item[E])
	}
}

// IsExpired judge whether the data is expired
func (c *mapCache[E]) IsExpired(key string) (bool, error) {
	value, ok := c.items[key]
	if !ok {
		return false, fmt.Errorf("the data %s does not exist", key)
	}
	return value.expired(), nil
}

// DeleteExpired delete all expired data
func (c *mapCache[E]) DeleteExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for k, v := range c.items {
		if v.expired() {
			c.del(k)
		}
	}
}

// Delete delete data by key
func (c *mapCache[E]) Delete(key string) (E, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	value, ok := c.get(key)
	if ok {
		c.del(key)
		return value.Object, ok
	}
	var zero E
	return zero, ok
}

// Set  data by key，it will overwrite the data if the key exists
func (c *mapCache[E]) Set(key string, value E) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.judgeAndInitItem()

	c.set(key, value, c.generateExpiration())
}

// SetDefault  data by key，it will overwrite the data if the key exists
func (c *mapCache[E]) SetDefault(key string, value E, expiration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.judgeAndInitItem()

	c.set(key, value, c.generateExpirationForItem(expiration))
}

// Add data，Cannot add existing data
// To override the addition, use the set method
func (c *mapCache[E]) Add(key string, value E) error {
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
func (c *mapCache[E]) Get(key string) (E, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	value, ok := c.items[key]
	if !ok || value.expired() {
		var zero E
		return zero, false
	}
	return value.Object, true
}

// GetAndDelete get data and delete by key
func (c *mapCache[E]) GetAndDelete(key string) (E, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	value, ok := c.items[key]
	if !ok || value.expired() {
		var zero E
		return zero, false
	}
	// delete
	c.del(key)
	return value.Object, true
}

// GetAndExpired  get data and expire by key
// It will be deleted at the next clearing. If the clearing capability is not enabled, it will never be deleted
func (c *mapCache[E]) GetAndExpired(key string) (E, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	value, ok := c.items[key]
	if !ok || value.expired() {
		var zero E
		return zero, false
	}
	// SetDefault now as expiration time
	c.set(key, value.Object, time.Now().UnixNano()/1e3)
	return value.Object, true
}

func (c *mapCache[E]) GetWithExpiration(key string) (E, time.Time, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	value, ok := c.items[key]
	if !ok || value.expired() {
		var zero E
		return zero, time.Time{}, false
	}
	return value.Object, time.UnixMicro(value.Expiration), true
}

// Clear remove all data
func (c *mapCache[E]) Clear() {
	c.items = make(map[string]*Item[E])
}

// Keys get all keys
func (c *mapCache[E]) Keys() []string {
	res := make([]string, 0)
	for k := range c.items {
		res = append(res, k)
	}
	return res
}
