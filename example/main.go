package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/ixugo/crontab"
)

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
