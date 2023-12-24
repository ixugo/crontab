package crontab

import "github.com/gin-gonic/gin"

// Register 注册路由
func Register(g *gin.RouterGroup) {
	g.GET("/crontab", FindTasks)
	g.POST("/crontab/:key/exec", ExecTask)
	g.DELETE("/crontab/:key", StopTask)
	g.POST("/crontab/:key", StartTask)
	g.POST("/crontab/reload", ReloadTasks)
}

// FindTasks 查询全部任务列表
func FindTasks(c *gin.Context) {
	c.IndentedJSON(200, gin.H{
		"items": Default().Tasks(),
	})
}

// StopTask 停止指定任务
func StopTask(c *gin.Context) {
	key := c.Param("key")
	if err := Default().Stop(key); err != nil {
		c.JSON(400, gin.H{"msg": err})
		return
	}
	c.JSON(200, gin.H{"key": key})
}

// StartTask 开始指定任务
func StartTask(c *gin.Context) {
	key := c.Param("key")
	if err := Default().Start(key); err != nil {
		c.JSON(400, gin.H{"msg": err})
		return
	}
	c.JSON(200, gin.H{"key": key})
}

// ExecTask 立即执行任务
func ExecTask(c *gin.Context) {
	key := c.Param("key")
	if err := Default().Exec(key); err != nil {
		c.JSON(400, gin.H{"msg": err})
		return
	}
	c.JSON(200, gin.H{"key": key})
}

// ReloadTasks 重载任务，用于配置文件更新
func ReloadTasks(c *gin.Context) {
	if err := Default().reload(); err != nil {
		c.JSON(400, gin.H{"msg": err})
		return
	}
	c.JSON(200, gin.H{"msg": "ok"})
}
