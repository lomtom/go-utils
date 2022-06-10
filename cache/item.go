package cache

import (
	"time"
)

type Item struct {
	Object     interface{} // data
	Expiration int64       // expiration time
}

// judge whether data is expired
func (item *Item) expired() bool {
	if item.Expiration == 0 {
		return false
	}
	return time.Now().UnixNano()/1e3 > item.Expiration
}

// Set the expiration time, and the data will be cleared in the next cache cleaning cycle
func (item *Item) setExpired() {
	item.Expiration = time.Now().UnixNano() / 1e3
}
