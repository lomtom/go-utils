package cache

import "time"

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
