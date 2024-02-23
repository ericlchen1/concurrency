package main

import (
	"fmt"
	"sync"
)

const NUM_THREADS = 10
const NUM_JOBS = 100

type Job func()

type ThreadPool struct {
	workQueue chan Job
	wg        sync.WaitGroup
}

func NewThreadPool(numThreads int) *ThreadPool {
	threadPool := &ThreadPool{
		workQueue: make(chan Job),
	}

	for i := 0; i < numThreads; i++ {
		threadPool.wg.Add(1)
		go func() {
			defer threadPool.wg.Done()
			for job := range threadPool.workQueue {
				job()
			}
		}()
	}
	return threadPool
}

func (t *ThreadPool) AddJob(job Job) {
	t.workQueue <- job
}

func (t *ThreadPool) Wait() {
	close(t.workQueue)
	t.wg.Wait()
}

func main() {
	threadPool := NewThreadPool(NUM_THREADS)

	for i := 0; i < 100; i++ {
		threadPool.AddJob(func() {
			fmt.Printf("Job %d\n", i)
		})
	}

	threadPool.Wait()
}
