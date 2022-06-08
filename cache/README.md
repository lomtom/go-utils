介绍
---

**1. map类型缓存**
- 提供间隔时间对数据进行过期清理
- 可手动开启/停止清理能力
- 可手动清除全部缓存
- 缓存持久化（计划中）
- ...

接口
---
```go
// IsExpired judge whether the data is expired
IsExpired(key string) (bool, error)
// DeleteExpired delete all expired data
DeleteExpired()

// StartGc start gc
// After the expiration time is set, GC will be started automatically without manual GC
StartGc() error
// StopGc stop gc
StopGc() error

// Get data
// When the data does not exist or expires, it will return nonexistence（false）
Get(key string) (interface{}, bool)
// GetAndDelete get data and delete by key
GetAndDelete(key string) (interface{}, bool)
// GetAndExpired  get data and expire by key
// It will be deleted at the next clearing. If the clearing capability is not enabled, it will never be deleted
GetAndExpired(key string) (interface{}, bool)

// Delete delete data by key
Delete(key string) (interface{}, bool)


// Set  data by key，it will overwrite the data if the key exists
Set(key string, value interface{})
// Add data，Cannot add existing data
// To override the addition, use the set method
Add(key string, value interface{}) error
// Clear remove all data
Clear()
// Keys get all keys
Keys() []string
```

使用
---
**引用**
```go
import "github.com/lomtom/go-utils/cache"
```

**示例**
```go
func Test(t *testing.T) {{
    c := cache.NewMapCache()
    c.Set("1", 1)
    fmt.Println(c.Get("1"))
}
```
**输出**
```shell
1 true
```

