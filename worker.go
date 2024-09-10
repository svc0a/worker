package worker

import (
	"github.com/sirupsen/logrus"
	"sync"
)

type Pool[T any] interface {
	Start(callback func(data T) error)
	Submit(data T)
	Stop()
}

// workerPool 管理协程池
type workerPool[T any] struct {
	numWorkers int
	taskChan   chan T
	wg         sync.WaitGroup
	once       sync.Once
	errHandler func(err error)
}

type Option[T any] func(pool *workerPool[T])

func WithErrHandler[T any](in func(err error)) Option[T] {
	return func(wp *workerPool[T]) {
		wp.errHandler = in
	}
}

// New 创建新的 workerPool 实例
func New[T any](numWorkers int, options ...Option[T]) Pool[T] {
	wp := &workerPool[T]{
		numWorkers: numWorkers,
		taskChan:   make(chan T),
		errHandler: func(err error) {
			logrus.Error(err)
		},
	}
	for _, o := range options {
		o(wp)
	}
	return wp
}

// Start 启动协程池
func (wp *workerPool[T]) Start(callback func(data T) error) {
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
func (wp *workerPool[T]) Submit(data T) {
	wp.taskChan <- data
}

// Stop 关闭协程池并等待所有协程完成
func (wp *workerPool[T]) Stop() {
	wp.once.Do(func() {
		close(wp.taskChan) // 关闭任务通道，通知协程停止接收新任务
		wp.wg.Wait()       // 等待所有协程完成任务
	})
}
