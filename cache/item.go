package cache

import "time"

type item struct {
	object     interface{} // data
	expiration int64       // expiration time
}

// judge whether data is expired
func (item *item) expired() bool {
	if item.expiration == 0 {
		return false
	}
	return time.Now().UnixNano()/1e3 > item.expiration
}

// Set the expiration time, and the data will be cleared in the next cache cleaning cycle
func (item *item) setExpired() {
	item.expiration = time.Now().UnixNano() / 1e3
}
