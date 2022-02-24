package cache

import "time"

/******************************************* 数据 *******************************************/

type item struct {
	object     interface{} // 数据
	expiration int64       // 过期时间
}

// Expired 判断数据项是否已经过期
func (item *item) expired() bool {
	if item.expiration == 0 {
		return false
	}
	return time.Now().UnixNano()/1e3 > item.expiration
}

// 设置过期，将会在下个清理缓存的周期进行清除
func (item *item) setExpired() {
	item.expiration = time.Now().UnixNano() / 1e3
}
