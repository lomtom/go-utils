package test

import (
	"fmt"
	"github.com/lomtom/go-utils/v2/job"
	"testing"
	"time"
)

func Test(t *testing.T) {
	j := job.NewTimerJob(func(j job.TimerJob) {
		fmt.Println("这是一个定时任务")
	},
		// 设置间隔时间（默认一分钟）
		job.SetDuration(time.Second))
	err := j.Start()
	if err != nil {
		fmt.Println(err)
		return
	}
	time.Sleep(time.Second * 4)
}

func Scan(j job.TimerJob) {
	m := j.GetParam()
	v, ok := m["1"]
	if ok {
		m["1"] = fmt.Sprintf("%v + 1", v)
	} else {
		m["1"] = "1"
	}
	err := j.SetParam(m)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(m["1"])
}

func TestTimer(t *testing.T) {
	m := make(map[string]interface{})
	j := job.NewTimerJob(Scan, job.SetDuration(time.Second*2), job.SetParam(m))
	err := j.Start()
	if err != nil {
		fmt.Println(err)
		return
	}
	time.Sleep(time.Second * 4)
	err = j.Stop()
	if err != nil {
		fmt.Println(err)
		return
	}
	time.Sleep(time.Second * 2)
	err = j.Start()
	if err != nil {
		fmt.Println(err)
		return
	}
	time.Sleep(time.Minute)
}

func Job1(job.TimerJob) {
	fmt.Println("job1")
}

func Job2(job.TimerJob) {
	fmt.Println("job2")
}

func TestPool(t *testing.T) {
	job1 := job.NewTimerJob(Job1, job.SetDuration(time.Second*2))
	job2 := job.NewTimerJob(Job2, job.SetDuration(time.Second*2))
	pool, err := job.NewPool(job1, job2)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = pool.StartAll()
	if err != nil {
		fmt.Println(err)
		return
	}
	time.Sleep(time.Second * 4)
	err = pool.StopAll()
	if err != nil {
		fmt.Println(err)
		return
	}
	time.Sleep(time.Second * 2)
	err = pool.StartAll()
	if err != nil {
		fmt.Println(err)
		return
	}
	time.Sleep(time.Minute)
}

func TestPool1(t *testing.T) {
	job1 := job.NewTimerJob(Job1, job.SetDuration(time.Second*2), job.SetName("job1"))
	job2 := job.NewTimerJob(Job2, job.SetDuration(time.Second*2), job.SetName("job2"))
	pool, err := job.NewPool(job1)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = pool.StartAll()
	if err != nil {
		fmt.Println(err)
		return
	}
	time.Sleep(time.Second * 4)
	err = pool.Add(job2)
	if err != nil {
		fmt.Println(err)
		return
	}
	time.Sleep(time.Second * 2)
	err = pool.Remove(job1)
	if err != nil {
		fmt.Println(err)
		return
	}
	time.Sleep(time.Minute)
}
