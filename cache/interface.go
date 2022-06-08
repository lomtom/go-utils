package cache

type Interface interface {
	// IsExpired judge whether the data is expired
	IsExpired(key string) (bool, error)
	// DeleteExpired delete all expired data
	DeleteExpired()

	// StartGc start gc
	// After the expiration time is set, GC will be started automatically without manual GC
	StartGc() error
	// StopGc stop gc
	StopGc() error

	// Get  data
	// When the data does not exist or expires, it will return nonexistence（false）
	Get(key string) (interface{}, bool)
	// GetAndDelete get data and delete by key
	GetAndDelete(key string) (interface{}, bool)
	// GetAndExpired  get data and expire by key
	// It will be deleted at the next clearing. If the clearing capability is not enabled, it will never be deleted
	GetAndExpired(key string) (interface{}, bool)

	// Delete delete data by key
	Delete(key string) (interface{}, bool)
}

type MapInterface interface {
	Interface

	// Set  data by key，it will overwrite the data if the key exists
	Set(key string, value interface{})
	// Add data，Cannot add existing data
	// To override the addition, use the set method
	Add(key string, value interface{}) error
	// Clear remove all data
	Clear()
	// Keys get all keys
	Keys() []string
}
