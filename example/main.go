package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/ixugo/crontab"
)

func main() {
	// 注册业务函数
	crontab.Add("function1", func(crontab.Params) error {
		fmt.Println("function1")
		return nil
	})
	crontab.Add("function2", func(crontab.Params) error {
		fmt.Println("function2")
		return nil
	})
	if err := crontab.Run(); err != nil {
		panic(err)
	}

	// 注册路由
	g := gin.Default()
	api := g.Group("/api")
	crontab.Register(api)
	_ = g.Run(":8081")
}
