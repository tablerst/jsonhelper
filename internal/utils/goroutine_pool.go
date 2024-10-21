package utils

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Task接口，每个任务都需要实现Execute方法
type Task interface {
	Execute()
}

// ExampleTask是Task接口的一个简单实现
type ExampleTask struct {
	ID int
}

func (t ExampleTask) Execute() {
	fmt.Printf("Executing task %d\n", t.ID)
	time.Sleep(time.Second) // 模拟任务处理耗时
}

// Worker结构体从任务通道接收任务并执行
type Worker struct {
	id int
}

func (w *Worker) Start(ctx context.Context, wg *sync.WaitGroup, tasks <-chan Task) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Worker %d exiting\n", w.id)
			return
		case task, ok := <-tasks:
			if !ok {
				fmt.Printf("Worker %d no more tasks\n", w.id)
				return
			}
			task.Execute()
		}
	}
}

// Pool结构体管理Worker和Task
type Pool struct {
	taskQueue   chan Task
	workerCount int
	ctx         context.Context
	cancel      context.CancelFunc
	wg          sync.WaitGroup
}

func NewPool(workerCount int, queueSize int) *Pool {
	ctx, cancel := context.WithCancel(context.Background())
	return &Pool{
		taskQueue:   make(chan Task, queueSize),
		workerCount: workerCount,
		ctx:         ctx,
		cancel:      cancel,
	}
}

func (p *Pool) Start() {
	for i := 0; i < p.workerCount; i++ {
		worker := Worker{id: i}
		p.wg.Add(1)
		go worker.Start(p.ctx, &p.wg, p.taskQueue)
	}
}

func (p *Pool) Submit(task Task) {
	p.taskQueue <- task
}

func (p *Pool) Shutdown() {
	close(p.taskQueue)
	p.cancel()
	p.wg.Wait()
}

func main() {
	pool := NewPool(5, 10)
	pool.Start()

	for i := 0; i < 20; i++ {
		pool.Submit(ExampleTask{ID: i})
	}

	// 安全关闭池
	pool.Shutdown()
}
