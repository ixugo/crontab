用于定时任务可视化

![image](http://img.golang.space/img-1703425476111.png)

## 快速开始

1. GET Package
```bash
go get -u github.com/ixugo/crontab
```

2. 编辑配置
在可执行文件同目录下，或 configs 目录下，写入 crontab.yaml 文件。
```yaml
version: 1
tasks:
  - key: task1
    title: 定时任务1
    description: 这是定时任务1的描述，每 2 秒执行一次
    cron: "*/2 * * * * *"
    func:
      name: logic
      params:
        expired: 1h

  - key: task2
    title: 定时任务2
    description: 这是定时任务2的描述，每 5 秒执行一次
    cron: "*/5 * * * * *"
    func:
      name: function2
      params:
        expired: 1h
```

3. 执行定时任务
```go

func logic(crontab.Params) error {
	fmt.Println("function1")
	return nil
}
func main() {
	// 注册业务函数
	crontab.Register(logic)
	// 不允许匿名函数
	// crontab.Register(func(crontab.Params) error {
	// 	fmt.Println("function2")
	// 	return nil
	// })
	if err := crontab.Run(); err != nil {
		panic(err)
	}

	// 注册路由
	g := gin.Default()
	api := g.Group("/api")
	crontab.RegisterAPI(api)
	_ = g.Run(":8081")
}
```

提供了默认的 5 个 API

未提供前端页面

+ `GET /crontab`  查询所有任务
+ `POST /crontab/:key/exec` 立即执行指定任务
+ `DELETE /crontab/:key` 停止指定任务
+ `POST /crontab/:key` 启动指定任务
+ `POST /crontab/reload` 重新加载配置文件


常用的设置，注意从秒开始。

每间隔 5 分钟执行一次 `0 */5 * * * *`

每天凌晨 1 点执行一次 `0 0 1 * * *`

每 3 天的凌晨 2 点执行一次 `0 0 2 */3 * *`

每个小时的执行一次 `0 0 * * * *`

每秒执行 `@every 1s`
