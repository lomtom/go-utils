package cache

type Interface interface {
	// IsExpired 判断是否过期
	IsExpired(key string) (bool, error)
	// DeleteExpired 删除过期数据项
	DeleteExpired()

	// StartGc 重新gc
	// 设置过期时间后，会自动开启gc，无需手动gc
	StartGc() error
	StopGc() error

	// Get 获取数据
	// 不存在或过期都会返回不存在
	// 返回数据、是否存在
	Get(key string) (interface{}, bool)
	// GetAndDelete 获取数据并删除
	GetAndDelete(key string) (interface{}, bool)
	// GetAndExpired  获取数据并过期
	// 将在下一次清除时删除，若未开启清除能力，将永远不会删除
	GetAndExpired(key string) (interface{}, bool)

	// Delete 删除数据
	Delete(key string) (interface{}, bool)
}

type MapInterface interface {
	Interface

	// Set 添加/修改数据，将会覆盖
	Set(key string, value interface{})
	// Add 添加数据，若有相同
	// 如需覆盖添加，请使用Set方法
	Add(key string, value interface{}) error
	// Clear 清除所有数据
	Clear()
}
