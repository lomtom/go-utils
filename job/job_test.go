package job

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func Scan(j TimerJob) {
	m := j.GetParam()
	v, ok := m["1"]
	if ok {
		m["1"] = fmt.Sprintf("%v + 1", v)
	} else {
		m["1"] = "1"
	}
	err := j.SetParam(m)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(m["1"])
}

func TestTimer(t *testing.T) {
	m := make(map[string]interface{})
	job := NewTimerJob(Scan, SetDuration(time.Second*2), SetParam(m))
	err := job.Start()
	if err != nil {
		log.Println(err)
		return
	}
	time.Sleep(time.Second * 4)
	err = job.Stop()
	if err != nil {
		log.Println(err)
		return
	}
	time.Sleep(time.Second * 2)
	err = job.Start()
	if err != nil {
		log.Println(err)
		return
	}
	time.Sleep(time.Minute)
}

func Job1(j TimerJob) {
	log.Println("job1")
}

func Job2(j TimerJob) {
	log.Println("job2")
}

func TestPool(t *testing.T) {
	job1 := NewTimerJob(Job1, SetDuration(time.Second*2))
	job2 := NewTimerJob(Job2, SetDuration(time.Second*2))
	pool, err := NewPool(job1, job2)
	if err != nil {
		log.Println(err)
		return
	}
	err = pool.StartAll()
	if err != nil {
		log.Println(err)
		return
	}
	time.Sleep(time.Second * 4)
	err = pool.StopAll()
	if err != nil {
		log.Println(err)
		return
	}
	time.Sleep(time.Second * 2)
	err = pool.StartAll()
	if err != nil {
		log.Println(err)
		return
	}
	time.Sleep(time.Minute)
}

func TestPool1(t *testing.T) {
	job1 := NewTimerJob(Job1, SetDuration(time.Second*2), SetName("job1"))
	job2 := NewTimerJob(Job2, SetDuration(time.Second*2), SetName("job2"))
	pool, err := NewPool(job1)
	if err != nil {
		log.Println(err)
		return
	}
	err = pool.StartAll()
	if err != nil {
		log.Println(err)
		return
	}
	time.Sleep(time.Second * 4)
	err = pool.Add(job2)
	if err != nil {
		log.Println(err)
		return
	}
	time.Sleep(time.Second * 2)
	err = pool.Remove(job1)
	if err != nil {
		log.Println(err)
		return
	}
	time.Sleep(time.Minute)
}
