package bytesdispatcher

import (
	"github.com/inverse-inc/packetfence/go/bytearraypool"
	"sync"
)

// BytesHandler an interface for handling bytes
type BytesHandler interface {
	HandleBytes([]byte)
}

// The BytesHandlerFunc type is an adapter to allow the use of
// ordinary functions as a []byte handlers. If f is a function
// with the appropriate signature, BytesHandlerFunc(bytes) is a
// Handler that calls f.
type BytesHandlerFunc func([]byte)

// HandleBytes calls f(bytes)
func (f BytesHandlerFunc) HandleBytes(bytes []byte) {
	f(bytes)
}

type worker struct {
	workerPool    chan<- chan []byte
	jobChannel    chan []byte
	bytesHandler  BytesHandler
	byteArrayPool *bytearraypool.ByteArrayPool
	stopChannel   chan struct{}
	waitGroup     *sync.WaitGroup
}

func initWorker(w *worker, workerPool chan<- chan []byte, bytesHandler BytesHandler, byteArrayPool *bytearraypool.ByteArrayPool, waitGroup *sync.WaitGroup) {
	w.workerPool = workerPool
	w.bytesHandler = bytesHandler
	w.byteArrayPool = byteArrayPool
	w.jobChannel = make(chan []byte)
	w.stopChannel = make(chan struct{})
	w.waitGroup = waitGroup
}

func (w *worker) handleBytes(bytes []byte) {
	defer w.byteArrayPool.Put(bytes)
	w.bytesHandler.HandleBytes(bytes)
}

func (w *worker) start() {
	go func() {
		defer w.waitGroup.Done()
	LOOP:
		for {
			// register the current worker into the worker queue.
			w.workerPool <- w.jobChannel

			select {
			case job := <-w.jobChannel:
				w.handleBytes(job)

			case <-w.stopChannel:
				// we have received a signal to stop
				break LOOP
			}
		}

		// Handle any leftover jobs
		select {
		case job := <-w.jobChannel:
			w.handleBytes(job)
		default:
			return
		}
	}()
}

func (w *worker) stop() {
	go func() {
		w.stopChannel <- struct{}{}
	}()
}

// Dispatcher dispatches work to a set of workers
type Dispatcher struct {
	// A pool of workers channels that are registered with the dispatcher
	maxWorkers    int
	byteArrayPool *bytearraypool.ByteArrayPool
	bytesHandler  BytesHandler
	jobQueue      chan []byte
	workerPool    chan chan []byte
	workers       []worker
	waitGroup     sync.WaitGroup
}

// NewDispatcher create a new Dispatcher
func NewDispatcher(maxWorkers, jobQueueSize int, bytesHandler BytesHandler, byteArrayPool *bytearraypool.ByteArrayPool) *Dispatcher {
	return &Dispatcher{
		maxWorkers:    maxWorkers,
		bytesHandler:  bytesHandler,
		byteArrayPool: byteArrayPool,
		jobQueue:      make(chan []byte, jobQueueSize),
		workerPool:    make(chan chan []byte, maxWorkers),
	}
}

// SubmitJob submit a byte array to be processed
func (d *Dispatcher) SubmitJob(job []byte) {
	d.jobQueue <- job
}

// Run the dispatcher
func (d *Dispatcher) Run() {
	// starting n number of workers
	d.workers = make([]worker, d.maxWorkers)
	d.waitGroup.Add(d.maxWorkers)
	for i := 0; i < d.maxWorkers; i++ {
		initWorker(&d.workers[i], d.workerPool, d.bytesHandler, d.byteArrayPool, &d.waitGroup)
		d.workers[i].start()
	}

	go d.dispatch(d.jobQueue, d.workerPool)
}

// Stop the dispatcher
func (d *Dispatcher) Stop() {
	for i := range d.workers {
		d.workers[i].stop()
	}

	d.Wait()
}

// Wait for all the workers to be finished
func (d *Dispatcher) Wait() {
	d.waitGroup.Wait()
}

func (d *Dispatcher) dispatch(jobQueue <-chan []byte, workerPool <-chan chan []byte) {
	for {
        // Find a worker
        jobChannel := <-workerPool
        //Get a Job
        job := <-jobQueue
        // Send it to the worker queue
        jobChannel <- job
	}
}
