package cache

import "time"

const (
	// DefaultExpiration Default expiration time flag， never expires
	DefaultExpiration time.Duration = -1

	// DefaultInterval Default expiration interval is one minute
	DefaultInterval = time.Minute

	// DefaultPersistencePath default persistence path
	DefaultPersistencePath = "/val/cache/persistence"
)

// expiration policy
type expirationOption struct {
	expiration time.Duration // Expiration time
	gcInterval time.Duration // Overdue data Item cleaning cycle
}

// persistencePolicy policy
type persistenceOption struct {
	persistenceName   string      // persistence name
	enablePersistence bool        // enable persistencePolicy
	persistencePolicy Persistence // persistencePolicy policy
	persistencePath   string      // persistencePath
}

type options struct {
	expirationOption
	persistenceOption
}

func newOption() options {
	return options{
		expirationOption{
			expiration: DefaultExpiration,
			gcInterval: DefaultInterval,
		},
		persistenceOption{
			enablePersistence: false,
			persistencePolicy: FFB,
			persistencePath:   DefaultPersistencePath,
		},
	}
}

// CreateOptionFunc Initialize optional parameters
type CreateOptionFunc func(o *options)

// SetExpirationTime  set expiration time
// expiration time
func SetExpirationTime(expiration time.Duration) CreateOptionFunc {
	return func(o *options) {
		o.expiration = expiration
	}
}

// SetGcInterval  set gc interval
// When the cleaning cycle is 0, it is automatically adjusted to 1 minute
func SetGcInterval(gcInterval time.Duration) CreateOptionFunc {
	if gcInterval == 0 {
		gcInterval = time.Minute
	}
	return func(o *options) {
		o.gcInterval = gcInterval
	}
}

// SetEnablePersistence Set whether to enable persistencePolicy
func SetEnablePersistence(name string) CreateOptionFunc {
	return func(o *options) {
		o.enablePersistence = true
		o.persistenceName = name
	}
}

// SetPersistencePolicy  set persistencePolicy policy,default persistencePolicy is FFB
func SetPersistencePolicy(policy Persistence) CreateOptionFunc {
	return func(o *options) {
		o.persistencePolicy = policy
	}
}

// SetPersistencePath  set persistence path,default persistence path is DefaultPersistencePath
func SetPersistencePath(path string) CreateOptionFunc {
	return func(o *options) {
		o.persistencePath = path
	}
}
