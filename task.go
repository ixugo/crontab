package crontab

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
	"gopkg.in/yaml.v3"
)

var e *Engine
var once sync.Once

// 错误
var (
	ErrNoExistTask = fmt.Errorf("任务不存在")
	ErrNoExistFunc = fmt.Errorf("函数不存在")
)

// Handler 业务函数
type Handler func(Params) error

// Default 提供默认实例
func Default() *Engine {
	once.Do(func() {
		e = &Engine{
			data: make(map[string]Handler),
			cron: cron.New(cron.WithSeconds()),
		}
	})
	return e
}

// Add 添加任务函数
func Add(name string, task func(Params) error) {
	Default().Add(name, task)
}

// Run 从可执行程序目录下加载任务
// crontab.yaml 或 configs/crontab.yaml，如果找到配置文件，则加载
func Run() error {
	return Default().init().Run()
}

// Engine ...
type Engine struct {
	data  map[string]Handler
	tasks []*Task
	cron  *cron.Cron
}

// Add 添加任务函数
func (e *Engine) Add(name string, task Handler) {
	e.data[name] = task
}

// Run 运行任务
func (e *Engine) Run(tasks ...*Task) error {
	if len(e.tasks) == 0 {
		e.tasks = tasks
	}
	for i := range e.tasks {
		t := e.tasks[i]
		f, exist := e.data[t.Func.Name]
		if !exist {
			slog.Error("func not found", slog.String("key", t.Key))
			continue
		}
		id, err := e.cron.AddFunc(t.Cron, func() {
			t.LastTimeAt = time.Now().Format(time.DateTime)
			t.Count++
			if err := f(t.Func.Params); err != nil {
				t.Result = err.Error()
			} else {
				t.Result = "OK"
			}
		})
		if err != nil {
			return err
		}
		t.id = id
	}
	e.cron.Start()
	return nil
}

// Stop 停止任务
func (e *Engine) Stop(key string) error {
	for _, v := range e.tasks {
		if key == v.Key {
			if v.id <= 0 {
				return nil
			}
			e.cron.Remove(v.id)
			v.id = -1
			return nil
		}
	}
	return ErrNoExistTask
}

// Start 开始任务
func (e *Engine) Start(key string) error {
	for _, t := range e.tasks {
		if key == t.Key {
			f, exist := e.data[t.Func.Name]
			if !exist {
				return ErrNoExistFunc
			}
			t.id, _ = e.cron.AddFunc(t.Cron, func() {
				t.LastTimeAt = time.Now().Format(time.DateTime)
				t.Count++
				if err := f(t.Func.Params); err != nil {
					t.Result = err.Error()
				} else {
					t.Result = "OK"
				}
			})
			return nil
		}
	}
	return ErrNoExistTask
}

// Exec 立即执行任务
func (e *Engine) Exec(key string) error {
	for _, t := range e.tasks {
		if key == t.Key {
			f, exist := e.data[t.Func.Name]
			if !exist {
				return ErrNoExistFunc
			}
			return f(t.Func.Params)
		}
	}
	return ErrNoExistTask
}

func (e *Engine) init() *Engine {
	dir := filepath.Dir(os.Args[0])
	{
		f := filepath.Join(dir, "crontab.yaml")
		if b, err := os.ReadFile(f); err == nil { // nolint
			var out Model
			if err := yaml.Unmarshal(b, &out); err == nil {
				e.tasks = out.Tasks
				return e
			}
		}
	}
	{
		f := filepath.Join(dir, "configs/crontab.yaml")
		if b, err := os.ReadFile(f); err == nil { // nolint
			var out Model
			if err := yaml.Unmarshal(b, &out); err == nil {
				e.tasks = out.Tasks
				return e
			}
		}
	}
	return e
}
