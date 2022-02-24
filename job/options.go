package job

import (
	"fmt"
	"reflect"
	"time"
)

const (
	defaultName     = "job"
	defaultInterval = time.Minute
)

type option struct {
	name   string
	params map[string]interface{}
}

type timerOption struct {
	option
	interval time.Duration
}

func newTimerOption() timerOption {
	return timerOption{
		option{
			fmt.Sprintf("%v_%v", defaultName, time.Now().UnixNano()/1e3),
			nil,
		},
		defaultInterval,
	}
}

type CreateOptionFunc func(o *timerOption)

// SetName 设置名字
func SetName(name string) CreateOptionFunc {
	if reflect.ValueOf(name).IsZero() {
		name = defaultName + time.Now().Format("2006-01-02 15:04:05")
	}
	return func(o *timerOption) {
		o.name = name
	}
}

// SetDuration 设置间隔时间
// 不设置，默认间隔一分钟
func SetDuration(interval time.Duration) CreateOptionFunc {
	if reflect.ValueOf(interval).IsZero() {
		interval = defaultInterval
	}
	return func(o *timerOption) {
		o.interval = interval
	}
}

// SetParam 设置参数
func SetParam(param map[string]interface{}) CreateOptionFunc {
	return func(o *timerOption) {
		o.params = param
	}
}
