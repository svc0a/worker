package worker

import (
	"errors"
	"github.com/sirupsen/logrus"
	"sync"
)

type Worker[T any] interface {
	start(callback func(data T) error)
	submit(data T)
	stop()
}

type config struct {
	numWorkers int
	chanSize   int
	errHandler func(in error)
}

type worker[T any] struct {
	config
	taskChan chan T
	wg       sync.WaitGroup
	once     sync.Once
}

type Option func(c *config)

func WithErrHandler(f func(err error)) Option {
	return func(c *config) {
		c.errHandler = f
	}
}

func WithWorkerNumber(workerNumber int) Option {
	return func(c *config) {
		c.numWorkers = workerNumber
	}
}

func WithChanSize(size int) Option {
	return func(c *config) {
		c.chanSize = size
	}
}

func Walk[T any](list []T, f func(T) error, options ...Option) {
	wp := define[T](options...)
	defer wp.stop()
	wp.start(f)
	for _, t := range list {
		wp.submit(t)
	}
}

// define 创建新的 workerPool 实例
func define[T any](options ...Option) Worker[T] {
	c := &config{
		numWorkers: 10,
		chanSize:   10,
		errHandler: func(in error) {
			if in != nil {
				logrus.Error(in)
			}
		},
	}
	for _, o := range options {
		o(c)
	}
	var wp = &worker[T]{
		config:   *c,
		taskChan: make(chan T, c.numWorkers),
	}
	return wp
}

// Start 启动协程池
func (wp *worker[T]) start(callback func(data T) error) {
	for i := 0; i < wp.numWorkers; i++ {
		wp.wg.Add(1)
		go func() {
			defer wp.wg.Done()
			for data := range wp.taskChan {
				err := callback(data)
				if errors.Is(err, nil) {
					continue
				}
				wp.errHandler(err)
			}
		}()
	}
}

// Submit 提交任务到协程池
func (wp *worker[T]) submit(data T) {
	wp.taskChan <- data
}

// Stop 关闭协程池并等待所有协程完成
func (wp *worker[T]) stop() {
	wp.once.Do(func() {
		close(wp.taskChan) // 关闭任务通道，通知协程停止接收新任务
		wp.wg.Wait()       // 等待所有协程完成任务
	})
}

// PendingTasks 返回未处理的任务数量
func (wp *worker[T]) PendingTasks() int {
	return len(wp.taskChan)
}
