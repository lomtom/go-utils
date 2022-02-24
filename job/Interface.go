package job

// TimerJob 定时任务内可用
type TimerJob interface {
	// GetParam 获取参数
	GetParam() map[string]interface{}
	// SetParam 设置参数
	SetParam(params map[string]interface{}) error
}

type job interface {
	getName() string
	validate() error
}

type TimerJobInterface interface {
	// Start 开启任务
	Start() error
	// Stop 停止任务
	Stop() error
	job
}

type PoolInterface interface {
	// StartAll 开启全部任务
	StartAll() error
	// StopAll 停止全部任务
	StopAll() error
	// StopJob 停止某一个任务
	StopJob(j TimerJobInterface) error
	// StopJobByName 停止某一个任务（通过名字）
	StopJobByName(name string) error
	// StartJob 开启某一任务
	StartJob(j TimerJobInterface) error
	// StartJobByName  开启某一任务（通过名字）
	StartJobByName(name string) error
	// Add 放入任务(会立即启动)
	// 如果名字一样，将会返回错误
	Add(j TimerJobInterface) error
	// Remove 移除任务
	Remove(j TimerJobInterface) error
}
