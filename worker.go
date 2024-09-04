package worker

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"sync"
)

// WorkerPool 管理协程池
type WorkerPool[T any] struct {
	numWorkers int
	taskChan   chan T
	wg         sync.WaitGroup
}

// NewWorkerPool 创建新的 WorkerPool 实例
func NewWorkerPool[T any](numWorkers int) *WorkerPool[T] {
	return &WorkerPool[T]{
		numWorkers: numWorkers,
		taskChan:   make(chan T),
	}
}

// Start 启动协程池
func (wp *WorkerPool[T]) Start(uploadFunc func(data T) (string, error)) {
	for i := 0; i < wp.numWorkers; i++ {
		wp.wg.Add(1)
		go func() {
			defer wp.wg.Done()
			for data := range wp.taskChan {
				_, err := uploadFunc(data)
				if err != nil {
					logrus.Fatalf("error uploading file: %v", err)
				} else {
					fmt.Printf("Successfully uploaded %+v\n", data)
				}
			}
		}()
	}
}

// Submit 提交任务到协程池
func (wp *WorkerPool[T]) Submit(data T) {
	wp.taskChan <- data
}

// Stop 关闭协程池并等待所有协程完成
func (wp *WorkerPool[T]) Stop() {
	close(wp.taskChan)
	wp.wg.Wait()
}
