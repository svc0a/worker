package worker

import (
	"github.com/sirupsen/logrus"
	"sync"
)

type Pool[T any, Err any] interface {
	Start(callback func(data T) Err)
	Submit(data T)
	Stop()
}

// workerPool 管理协程池
type workerPool[T any, Err any] struct {
	numWorkers int
	taskChan   chan T
	wg         sync.WaitGroup
	once       sync.Once
	errHandler func(rsp Err)
}

type Option[T any, Err any] func(pool *workerPool[T, Err])

func WithErrHandler[T any, Err any](in func(err Err)) Option[T, Err] {
	return func(wp *workerPool[T, Err]) {
		wp.errHandler = in
	}
}
func WithChanSize[T any, Err any](in int64) Option[T, Err] {
	return func(wp *workerPool[T, Err]) {
		wp.taskChan = make(chan T, in)
	}
}

// New 创建新的 workerPool 实例
func New[T any, Err any](numWorkers int, options ...Option[T, Err]) Pool[T, Err] {
	wp := &workerPool[T, Err]{
		numWorkers: numWorkers,
		taskChan:   make(chan T, numWorkers),
		errHandler: func(err Err) {
			logrus.Error(err)
		},
	}
	for _, o := range options {
		o(wp)
	}
	return wp
}

// Start 启动协程池
func (wp *workerPool[T, Err]) Start(callback func(data T) Err) {
	for i := 0; i < wp.numWorkers; i++ {
		wp.wg.Add(1)
		go func() {
			defer wp.wg.Done()
			for data := range wp.taskChan {
				err := callback(data)
				if err != nil {
					wp.errHandler(err)
				}
			}
		}()
	}
}

// Submit 提交任务到协程池
func (wp *workerPool[T, Err]) Submit(data T) {
	wp.taskChan <- data
}

// Stop 关闭协程池并等待所有协程完成
func (wp *workerPool[T, Err]) Stop() {
	wp.once.Do(func() {
		close(wp.taskChan) // 关闭任务通道，通知协程停止接收新任务
		wp.wg.Wait()       // 等待所有协程完成任务
	})
}

// PendingTasks 返回未处理的任务数量
func (wp *workerPool[T, Err]) PendingTasks() int {
	return len(wp.taskChan)
}
