package crontab

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/robfig/cron/v3"
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

// Register 注册任务函数
func Register(name string, task func(Params) error) {
	// name := runtime.FuncForPC(reflect.ValueOf(task).Pointer()).Name()
	// name = name[strings.LastIndex(name, ".")+1:]
	// if strings.HasPrefix(name, "func") {
	// panic("不允许注册匿名函数，请更换函数名")
	// }
	Default().Register(name, task)
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

// Register 注册任务函数
func (e *Engine) Register(name string, task Handler) {
	_, exist := e.data[name]
	if exist {
		panic("定时任务不允许重名函数")
	}
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
		id, err := e.cron.AddFunc(t.Cron, e.wrap(t, f))
		if err != nil {
			return err
		}
		t.ID = id
	}
	e.cron.Start()
	return nil
}

// Tasks 获取任务列表
func (e *Engine) Tasks() []Task {
	es := e.cron.Entries()

	out := make([]Task, 0, len(e.tasks))
	for _, t := range e.tasks {
		if t.ID > 0 && int(t.ID) <= len(es) {
			v := getEntry(es, t.ID)
			if !v.Prev.IsZero() {
				t.LastTimeAt = v.Prev.Format(time.DateTime)
			}
			t.NextTimeAt = v.Next.Format(time.DateTime)
		}
		t := *t
		out = append(out, t)
	}

	return out
}

func getEntry(es []cron.Entry, id cron.EntryID) cron.Entry {
	for _, v := range es {
		if v.ID == id {
			return v
		}
	}
	return cron.Entry{}
}

// Stop 停止任务
func (e *Engine) Stop(key string) error {
	for _, t := range e.tasks {
		if key == t.Key {
			if t.ID <= 0 {
				return nil
			}
			e.cron.Remove(t.ID)
			t.ID = -1
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
			t.ID, _ = e.cron.AddFunc(t.Cron, e.wrap(t, f))
			return nil
		}
	}
	return ErrNoExistTask
}
func (e *Engine) wrap(t *Task, f Handler) func() {
	return func() {
		t.Count++
		if err := f(t.Func); err != nil {
			t.Result = err.Error()
		} else {
			t.Result = "OK"
		}
	}
}

// Exec 立即执行任务
func (e *Engine) Exec(key string) error {
	for _, t := range e.tasks {
		if key == t.Key {
			e.cron.Entry(t.ID).WrappedJob.Run()
			return nil
		}
	}
	return ErrNoExistTask
}

// Reload 重载任务，此重载从会指定路径读取配置文件
func (e *Engine) Reload() error {
	for _, t := range e.tasks {
		_ = e.Stop(t.Key)
	}
	e.cron.Stop()
	e.cron = cron.New(cron.WithSeconds())
	return e.init().Run()
}

func (e *Engine) init() *Engine {
	dir := filepath.Dir(os.Args[0])
	{
		f := filepath.Join(dir, "crontab.toml")
		if b, err := os.ReadFile(f); err == nil { // nolint
			if out, err := tasks(b); err == nil {
				e.tasks = out.Tasks
			}
		}
	}
	{
		f := filepath.Join(dir, "configs/crontab.toml")
		if b, err := os.ReadFile(f); err == nil { // nolint
			if out, err := tasks(b); err == nil {
				e.tasks = out.Tasks
			}
		}
	}
	return e
}

func tasks(b []byte) (Model, error) {
	var out Model
	_, err := toml.NewDecoder(bytes.NewReader(b)).Decode(&out)
	if err == nil {
		for _, v := range out.Tasks {
			v.ID = -1
		}
	}
	return out, err
}
