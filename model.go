package crontab

import (
	"time"

	"github.com/robfig/cron/v3"
)

// Model 配置
type Model struct {
	Version string  `yaml:"version"`
	Tasks   []*Task `yaml:"tasks"`
}

// Task 任务
type Task struct {
	Key         string `yaml:"key" json:"key"`
	Title       string `yaml:"title" json:"title"`
	Description string `yaml:"description" json:"description"`
	Cron        string `yaml:"cron" json:"cron"`
	Func        Params `yaml:"func" json:"-"`

	Record `yaml:"-"`
}

// Record 运行时的数据记录
type Record struct {
	LastTimeAt string       `yaml:"-" json:"last_time_at"`
	NextTimeAt string       `yaml:"-" json:"next_time_at"`
	Count      int64        `yaml:"-" json:"count"`
	Result     string       `yaml:"-" json:"result"`
	ID         cron.EntryID `yaml:"-" json:"id"`
}

// Params ...
type Params struct {
	Name    string        `yaml:"name" json:"name"`
	Expired time.Duration `yaml:"expired" json:"expired"`
}
