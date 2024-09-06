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
}

// New 创建新的 workerPool 实例
func New[T any](numWorkers int) Pool[T] {
	return &workerPool[T]{
		numWorkers: numWorkers,
		taskChan:   make(chan T),
	}
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
					logrus.Fatalf("error: %v", err)
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
