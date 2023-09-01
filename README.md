<h1 align=center>go工具包 ⚗️</h1>

<p align=center>这是一个Go语言工具包，旨在提供缓存、定时器和切片工具，为开发者提供方便和高效的工具，支持Go1.18的范型特性。</p>

<p align=center>Made with ♥ by <a href="https://lomtom.cn">lomtom</a></p>

<p align=center>如果你觉得这个项目有用，请给它一个⭐以表示你的支持</p>

## 特性 💡

1. **[Map缓存](cache/README.md)**: 提供了一个基于Go语言`map`的缓存工具，允许开发者轻松存储和检索键值对数据，加速数据访问。

2. **[定时器](job/README.md)**: 内置的定时器工具允许您创建定时任务，以便在指定的时间间隔内执行特定的操作，适用于需要定时执行代码的场景。

3. **[Slice工具](slice/README.md)**: 提供了一组操作切片的工具函数，使得在处理切片数据时更加便捷和高效。

## 安装  ⚒️

使用Go模块的方式，您可以轻松地将该工具包集成到您的项目中：
```shell
go get github.com/lomtom/go-utils
```

## 使用

以下是工具的简单使用，更多使用请点击标题查看详细介绍。

### 1. [Map缓存](cache/README.md)

使用Map缓存工具，您可以轻松存储和检索键值对数据，加速数据访问。
```go
package main

import (
    "fmt"
    "github.com/lomtom/go-utils/cache"
)


func main(){
    c,err := cache.NewMapCache[int]()
    if err != nil {
        fmt.Println("err:", err)
        return
    }
    c.Set("1", 1)
    fmt.Println(c.Get("1"))
}
```

### 2. **[定时器](job/README.md)**

使用定时器工具，您可以创建定时任务，以便在指定的时间间隔内执行特定的操作。
```go
package main

import (
	"fmt"
	"github.com/lomtom/go-utils/job"
)


func main(){
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

```go
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

### 3. **[Slice工具](slice/README.md)**

提供了一组操作切片的工具函数，使得在处理切片数据时更加便捷和高效。

```go
type people struct {
	name string
	age  int
}

func TestSize(t *testing.T) {
	a := assert.NewAssert(t)
	a.Equal(0, Size([]int{}))
	a.Equal(3, Size([]int{1, 2, 3}))
}

func TestIsEmpty(t *testing.T) {
	a := assert.NewAssert(t)
	a.Equal(true, IsEmpty([]int{}))
	a.Equal(false, IsEmpty([]int{1, 2, 3}))
}

func TestContains(t *testing.T) {
	a := assert.NewAssert(t)
	a.Equal(true, Contains([]int{1, 2, 3}, 1))
	a.Equal(true, Contains([]people{}, people{name: "lomtom"}))
}

func TestContainsAll(t *testing.T) {
	a := assert.NewAssert(t)
	a.Equal(true, ContainsAll([]int{1, 2, 3}))
	a.Equal(true, ContainsAll([]int{1, 2, 3}, 1, 2, 3))
	a.Equal(false, ContainsAll([]int{1, 2, 3}, 1, 2, 3, 4))
	a.Equal(true, ContainsAll([]people{{name: "lomtom"}}, people{name: "lomtom"}))
}

func TestAddAll(t *testing.T) {
	a := assert.NewAssert(t)
	a.Equal([]int{1, 2, 3, 1, 2, 3}, AddAll([]int{1, 2, 3}, []int{1, 2, 3}...))
	a.Equal([]people{{name: "lomtom"}, {name: "lomtom"}}, AddAll([]people{{name: "lomtom"}}, people{name: "lomtom"}))
}

func TestReplaceAll(t *testing.T) {
	a := assert.NewAssert(t)
	a.Equal([]int{4, 4, 3, 1, 4, 3}, ReplaceAll([]int{1, 2, 3, 1, 2, 3}, func(value int) int {
		if value == 2 {
			return 4
		}
		return value
	}))
}

func TestRemoveAll(t *testing.T) {
	a := assert.NewAssert(t)
	a.Equal([]int{4}, RemoveAll([]int{1, 2, 3, 4}, []int{1, 2, 3}...))
	a.Equal([]people{}, RemoveAll([]people{{name: "lomtom"}}, people{name: "lomtom"}))
}
```

注意 ⚠️
---
1. 由于v2使用Go 1.18 开发的工具包，低版本将不兼容
2. 如果是低版本请查看[v0.1.4](https://github.com/lomtom/go-utils/tree/v0.1.4)

## 贡献
欢迎贡献代码、报告问题或提出建议。请通过GitHub Issues和Pull Requests参与贡献。

## 许可证
该工具包基于 [MIT许可证](LICENSE) 发布。