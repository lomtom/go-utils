package cache

type Interface interface {
	IsExpired(key string) (bool, error)
	DeleteExpired()

	StartGc() error
	StopGc() error

	Get(key string) (interface{}, bool)
	GetAndDelete(key string) (interface{}, bool)
	GetAndExpired(key string) (interface{}, bool)

	Delete(key string) (interface{}, bool)
}

type MapInterface interface {
	Interface

	Set(key string, value interface{})
	Add(key string, value interface{}) error
	Clear()
}
