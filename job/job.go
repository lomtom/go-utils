package job

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"sync"
	"time"
)

type jobFunc func(j TimerJob)

type jobAction func(j *timerJob)

type timerJob struct {
	// 定时任务执行方法
	jf jobAction
	// 定时任务名称
	name string
	// 定时任务参数
	params map[string]interface{}
	// 可通过ch终止
	stopTimer chan int
	// 间隔时间
	interval time.Duration
	// 执行次数
	count int64
	// 编号
	id string
	// 重启次数
	startCount int64
	// 是否启动
	isStart bool
	// mu
	mu sync.Mutex
}

func jobRealAction(jobFunc2 jobFunc) jobAction {
	return func(j *timerJob) {
		j.mu.Lock()
		defer j.mu.Unlock()
		j.countIncrease()
		log.Printf("%v 第%v次  执行时任务 start....", j.getName(), j.getCount())
		jobFunc2(j)
		log.Printf("%v 第%v次  执行时任务 end....", j.getName(), j.getCount())
	}
}

// NewTimerJob
// 默认一分钟执行一次
func NewTimerJob(jf jobFunc, opts ...CreateOptionFunc) TimerJobInterface {
	createOption := newTimerOption()
	for _, opt := range opts {
		opt(&createOption)
	}
	return &timerJob{
		params:    createOption.params,
		stopTimer: make(chan int),
		jf:        jobRealAction(jf),
		name:      createOption.name,
		interval:  createOption.interval,
		id:        time.Now().Format("2006-01-02 15:04:05"),
		mu:        sync.Mutex{},
	}
}

func (j *timerJob) countIncrease() {
	j.count++
}

func (j *timerJob) getName() string {
	return j.name
}

func (j *timerJob) getCount() int64 {
	return j.count
}

func (j *timerJob) validate() error {
	if reflect.ValueOf(j.jf).IsZero() {
		return errors.New("定时任务不能为空")
	}
	if reflect.ValueOf(j.name).IsZero() {
		return errors.New("定时名称不能为空")
	}
	if reflect.ValueOf(j.interval).IsZero() {
		return errors.New("定时间隔不能为空")
	}
	if j.isStart {
		return errors.New(fmt.Sprintf("任务 %s 已经启动", j.name))
	}
	return nil
}

// Start 开启任务
func (j *timerJob) Start() error {
	err := j.validate()
	if err != nil {
		return err
	}
	j.startCount++
	j.isStart = true
	ticker := time.NewTicker(j.interval)
	log.Printf("%v 第%v次  开始执行任务", j.name, j.count+1)
	go func() {
		for {
			select {
			case <-ticker.C:
				j.jf(j)
			case <-j.stopTimer:
				ticker.Stop()
				return
			}
		}
	}()
	// 开启后，立马触发一次任务
	j.jf(j)
	return nil
}

// Stop 停止任务
func (j *timerJob) Stop() error {
	// 如果已经停止，跳过
	if !j.isStart {
		return errors.New(fmt.Sprintf("任务 %s 已经停止", j.name))
	}
	j.isStart = false
	log.Printf("%v 第%v次  停止执行任务", j.name, j.count)
	j.stopTimer <- 1
	return nil
}

// GetParam 获取参数
func (j *timerJob) GetParam() map[string]interface{} {
	return j.params
}

// SetParam 设置参数
func (j *timerJob) SetParam(params map[string]interface{}) error {
	j.params = params
	return nil
}
