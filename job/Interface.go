package job

type TimerJob interface {
	GetParam() map[string]interface{}
	SetParam(params map[string]interface{}) error
}

type job interface {
	getName() string
	validate() error
}

type TimerJobInterface interface {
	Start() error
	Stop() error
	job
}

type PoolInterface interface {
	StartAll() error
	StopAll() error
	StopJob(j TimerJobInterface) error
	StartJob(j TimerJobInterface) error
	Add(j TimerJobInterface) error
	Remove(j TimerJobInterface) error
}
