package job

import (
	"errors"
	"log"
	"sync"
)

var notExist = errors.New("任务不在该任务池内")
var exist = errors.New("该任务池内存在相同名称的任务")

type pool struct {
	// 用于存储所有的任务，用于开启和停止
	jobs map[string]TimerJobInterface
	lock sync.Mutex
}

// NewPool 放入的任务将不会自动开启需手动开启
func NewPool(jobs ...TimerJobInterface) (PoolInterface, error) {
	var l sync.Mutex
	jobsMap := make(map[string]TimerJobInterface)
	if len(jobs) != 0 {
		for _, j := range jobs {
			if _, ok := jobsMap[j.getName()]; ok {
				return nil, exist
			}
			jobsMap[j.getName()] = j
		}
	}
	return &pool{
		jobs: jobsMap,
		lock: l,
	}, nil
}

func (p *pool) get(name string) (TimerJobInterface, error) {
	j, ok := p.jobs[name]
	if ok {
		// 如果已经停止，跳过
		return j, nil
	}
	return nil, notExist
}

// 保证线程安全访问map
func (p *pool) add(name string, j TimerJobInterface) {
	p.jobs[name] = j
	log.Printf("%v 任务注册成功 ", j.getName())
}

func (p *pool) StartAll() error {
	p.lock.Lock()
	defer p.lock.Unlock()
	for _, job := range p.jobs {
		err := job.validate()
		if err != nil {
			return err
		}
	}
	for _, job := range p.jobs {
		_ = job.Start()
	}
	return nil
}

func (p *pool) StopAll() error {
	p.lock.Lock()
	defer p.lock.Unlock()
	for _, job := range p.jobs {
		_ = job.Stop()
	}
	return nil
}

func (p *pool) StopJob(j TimerJobInterface) error {
	p.lock.Lock()
	defer p.lock.Unlock()
	j, err := p.get(j.getName())
	if err != nil {
		return err
	}
	return j.Stop()
}

func (p *pool) StartJob(j TimerJobInterface) error {
	p.lock.Lock()
	defer p.lock.Unlock()
	j, err := p.get(j.getName())
	if err != nil {
		return err
	}
	return j.Start()
}

// Put 放入任务(会立即启动)
// 如果名字一样，将会返回错误
func (p *pool) Add(j TimerJobInterface) error {
	p.lock.Lock()
	defer p.lock.Unlock()
	_, err := p.get(j.getName())
	// 不存在则放入
	if err != nil {
		p.add(j.getName(), j)
		return j.Start()
	}
	return exist
}

func (p *pool) Remove(j TimerJobInterface) error {
	p.lock.Lock()
	defer p.lock.Unlock()
	j, err := p.get(j.getName())
	if err != nil {
		return err
	}
	err = j.Stop()
	if err != nil {
		return err
	}
	defer delete(p.jobs, j.getName())
	return nil
}
