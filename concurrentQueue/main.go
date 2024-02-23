package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

const ELEMENTS = 100_000

type Node struct {
	val  int64
	next *Node
	prev *Node
}

type ConcurrentQueue struct {
	head   *Node
	tail   *Node
	length int64
	mutex  sync.Mutex
}

func NewConcurrentQueue() *ConcurrentQueue {
	head := &Node{}
	tail := &Node{}
	head.next = tail
	tail.prev = head

	return &ConcurrentQueue{
		head:   head,
		tail:   tail,
		length: 0,
		mutex:  sync.Mutex{},
	}
}

func (q *ConcurrentQueue) Push(val int64) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	prevLastElement := q.tail.prev
	prevLastElement.next = &Node{
		val:  val,
		prev: prevLastElement,
		next: q.tail,
	}
	q.tail.prev = prevLastElement.next

	q.length++
}

func (q *ConcurrentQueue) Pop() (int64, error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if q.length == 0 {
		return 0, fmt.Errorf("queue is empty")
	}

	poppedElement := q.head.next
	q.head.next = poppedElement.next
	poppedElement.next.prev = q.head
	q.length--

	return poppedElement.val, nil
}

func main() {
	queue := NewConcurrentQueue()
	var sum int64 = 0

	pushGroup := sync.WaitGroup{}
	popGroup := sync.WaitGroup{}

	// Sweet new behavior for golang 1.22.2
	for i := range ELEMENTS {
		pushGroup.Add(1)
		go func(i int64) {
			defer pushGroup.Done()
			queue.Push(i)
		}(int64(i))

		popGroup.Add(1)
		go func() {
			defer popGroup.Done()
			for {
				val, err := queue.Pop()
				if err == nil {
					atomic.AddInt64(&sum, int64(val))
					break
				}
			}
		}()
	}

	pushGroup.Wait()
	popGroup.Wait()

	correctSum := int64(ELEMENTS * (ELEMENTS - 1) / 2)
	fmt.Println("Answer is correct:", sum == correctSum)
	fmt.Printf("Sum: %d, Correct sum: %d\n", sum, correctSum)
}
