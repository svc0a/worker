package worker

import (
	"errors"
	"github.com/sirupsen/logrus"
	"sync"
)

type Worker[T any, Err error] interface {
	start(callback func(data T) Err)
	submit(data T)
	stop()
}

type worker[T any, Error error] struct {
	numWorkers int
	taskChan   chan T
	wg         sync.WaitGroup
	once       sync.Once
	errHandler func(rsp Error)
}

type Option[T any, Error error] func(pool *worker[T, Error])

func WithErrHandler[T any, Error error](in func(err Error)) Option[T, Error] {
	return func(wp *worker[T, Error]) {
		wp.errHandler = in
	}
}
func WithChanSize[T any, Error error](in int64) Option[T, Error] {
	return func(wp *worker[T, Error]) {
		wp.taskChan = make(chan T, in)
	}
}

func Walk[T any, Error error](list []T, f func(T) Error, options ...Option[T, Error]) {
	wp := define[T, Error](1000, options...)
	defer wp.stop()
	for _, t := range list {
		wp.submit(t)
	}
	wp.start(f)
}

// define 创建新的 workerPool 实例
func define[T any, Error error](numWorkers int, options ...Option[T, Error]) Worker[T, Error] {
	wp := &worker[T, Error]{
		numWorkers: numWorkers,
		taskChan:   make(chan T, numWorkers),
		errHandler: func(err Error) {
			logrus.Error(err)
		},
	}
	for _, o := range options {
		o(wp)
	}
	return wp
}

// Start 启动协程池
func (wp *worker[T, Error]) start(callback func(data T) Error) {
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
func (wp *worker[T, Error]) submit(data T) {
	wp.taskChan <- data
}

// Stop 关闭协程池并等待所有协程完成
func (wp *worker[T, Error]) stop() {
	wp.once.Do(func() {
		close(wp.taskChan) // 关闭任务通道，通知协程停止接收新任务
		wp.wg.Wait()       // 等待所有协程完成任务
	})
}

// PendingTasks 返回未处理的任务数量
func (wp *worker[T, Error]) PendingTasks() int {
	return len(wp.taskChan)
}
