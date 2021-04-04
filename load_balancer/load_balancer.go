package load_balancer

import (
	"container/heap"
	"math/rand"
	"time"
)

// Request contains the operation to perform and
// the channel to return the result.
type Request struct {
	fn func() int
	c  chan int
}

func workFn() int {
	return rand.Int()
}

func furtherProcess(int) {
	return
}

// requester
func requester(work chan<- Request) {
	c := make(chan int)
	for {
		// fake load
		time.Sleep(time.Duration(rand.Int63n(2)) * time.Second)
		work <- Request{workFn, c}
		result := <-c //wait for answer
		furtherProcess(result)
	}
}

// Worker contains a channel of requests, plus some load tracking data.
type Worker struct {
	requests chan Request // work to do (buffered channel)
	pending  int          // count of pending tasks
	index    int          // index in the heap
}

// work gets request from balancer and do work. Balancer sends request to
// most lightly loaded worker. The channel of requests (w.requests) delivers
// requests to each worker. The balancer tracks the number of pending requests
// as a measure of load. Each response goes directly to its requester.
func (w *Worker) work(done chan *Worker) {
	for {
		req := <-w.requests
		req.c <- req.fn()
		done <- w
	}
}

// Pool is an implementation of Heap interface.
type Pool []*Worker

func (p Pool) Less(i, j int) bool {
	return p[i].pending < p[j].pending
}

func (p Pool) Len() int {
	return len(p)
}

func (p Pool) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p *Pool) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*p = append(*p, x.(*Worker))
}

func (p *Pool) Pop() interface{} {
	old := *p
	n := len(old)
	x := old[n-1]
	*p = old[0 : n-1]
	return x
}

// Balancer contains a pool of workers and a single channel to which requesters
// can report task completion.
type Balancer struct {
	pool Pool
	done chan *Worker
}

func (b *Balancer) balance(work chan Request) {
	for {
		select {
		case req := <-work: // receive a request
			b.dispatch(req) // send it to a worker
		case w := <-b.done: // a worker has finished
			b.completed(w) // update its info
		}
	}
}

// dispatch
func (b *Balancer) dispatch(req Request) {
	// Grab the least loaded worker...
	w := heap.Pop(&b.pool).(*Worker)
	// ...send it the task.
	w.requests <- req
	// One more in this work queue.
	w.pending++
	// Put it into its place on the heap.
	heap.Push(&b.pool, w)
}

// completed
// Job is complete; update heap
func (b *Balancer) completed(w *Worker) {
	// One fewer in the queue.
	w.pending--
	// Remove it from heap.
	heap.Remove(&b.pool, w.index)
	// Put it into its place on the heap.
	heap.Push(&b.pool, w)
}
