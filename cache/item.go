package cache

import (
	"time"
)

type Item[E any] struct {
	Object     E     // data
	Expiration int64 // expiration time
}

// judge whether data is expired
func (item *Item[E]) expired() bool {
	if item.Expiration == 0 {
		return false
	}
	return time.Now().UnixNano()/1e3 > item.Expiration
}

// Set the expiration time, and the data will be cleared in the next cache cleaning cycle
func (item *Item[E]) setExpired() {
	item.Expiration = time.Now().UnixNano() / 1e3
}
