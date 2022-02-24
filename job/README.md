
介绍
---


**1. 定时任务**
- 定时执行任务
- 在单个任务内可传递参数
- ...

**2. 任务管理池**
- 统一管理注入的任务
- ...

接口
---

任务接口：
```go
// GetParam 获取参数
GetParam() map[string]interface{}
// SetParam 设置参数
SetParam(params map[string]interface{}) error


// Start 开启任务
Start() error
// Stop 停止任务
Stop() error
```

任务池接口：
```go
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
```



使用
---
**引用**
```go
import "github.com/zero028/go-utils/job"
```

**示例**

示例一：使用定时任务
```go
func Test(t *testing.T) {
	j := job.NewTimerJob(func(j job.TimerJob) {
		fmt.Println("这是一个定时任务")
	},
	// 设置间隔时间（默认一分钟）
	job.SetDuration(time.Second))
	err := j.Start()
	if err != nil {
		log.Println(err)
		return
	}
	time.Sleep(time.Second * 4)
}
```
输出
```shell
2022/02/24 16:44:30 job_1645692270490900 第1次  开始执行任务
2022/02/24 16:44:31 job_1645692270490900 第1次  执行时任务 start....
这是一个定时任务
2022/02/24 16:44:31 job_1645692270490900 第1次  执行时任务 end....
2022/02/24 16:44:32 job_1645692270490900 第2次  执行时任务 start....
这是一个定时任务
2022/02/24 16:44:32 job_1645692270490900 第2次  执行时任务 end....
2022/02/24 16:44:33 job_1645692270490900 第3次  执行时任务 start....
这是一个定时任务
2022/02/24 16:44:33 job_1645692270490900 第3次  执行时任务 end....
2022/02/24 16:44:34 job_1645692270490900 第4次  执行时任务 start....
这是一个定时任务
2022/02/24 16:44:34 job_1645692270490900 第4次  执行时任务 end....
```

示例二： 使用定时任务设置参数，并且在任务内读取参数，并作修改
```go
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
```

输出
```shell
2022/02/24 16:45:41 job_1645692341457023 第1次  开始执行任务
2022/02/24 16:45:43 job_1645692341457023 第1次  执行时任务 start....
1
2022/02/24 16:45:43 job_1645692341457023 第1次  执行时任务 end....
2022/02/24 16:45:45 job_1645692341457023 第1次  停止执行任务
2022/02/24 16:45:45 job_1645692341457023 第2次  执行时任务 start....
1 + 1
2022/02/24 16:45:45 job_1645692341457023 第2次  执行时任务 end....
2022/02/24 16:45:47 job_1645692341457023 第3次  开始执行任务
2022/02/24 16:45:49 job_1645692341457023 第3次  执行时任务 start....
1 + 1 + 1
2022/02/24 16:45:49 job_1645692341457023 第3次  执行时任务 end....
2022/02/24 16:45:51 job_1645692341457023 第4次  执行时任务 start....
1 + 1 + 1 + 1
2022/02/24 16:45:51 job_1645692341457023 第4次  执行时任务 end....
2022/02/24 16:45:53 job_1645692341457023 第5次  执行时任务 start....
1 + 1 + 1 + 1 + 1
2022/02/24 16:45:53 job_1645692341457023 第5次  执行时任务 end....
2022/02/24 16:45:55 job_1645692341457023 第6次  执行时任务 start....
1 + 1 + 1 + 1 + 1 + 1
2022/02/24 16:45:55 job_1645692341457023 第6次  执行时任务 end....
2022/02/24 16:45:57 job_1645692341457023 第7次  执行时任务 start....
1 + 1 + 1 + 1 + 1 + 1 + 1
2022/02/24 16:45:57 job_1645692341457023 第7次  执行时任务 end....
```

示例三：使用任务池进行管理任务
```go
func Job1(j job.TimerJob) {
	fmt.Println("job1")
}

func Job2(j job.TimerJob) {
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
```