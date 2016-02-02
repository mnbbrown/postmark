package postmark

import (
	log "github.com/Sirupsen/logrus"
)

var TaskQueue = make(chan Task, 100)

type Dispatcher struct {
	WorkerPool chan chan Task
}

func NewDispatcher(maxWorkers int) *Dispatcher {
	pool := make(chan chan Task, maxWorkers)
	return &Dispatcher{pool}
}

func (d *Dispatcher) Run() {
	for i := 0; i < cap(d.WorkerPool); i++ {
		worker := NewWorker(i, d.WorkerPool)
		worker.Start()
	}
	go d.dispatch()
}

func (d *Dispatcher) dispatch() {
	for {
		select {
		case task := <-TaskQueue:
			go func(task Task) {
				taskChannel := <-d.WorkerPool
				taskChannel <- task
			}(task)
		}
	}
}

type Task struct {
	Message *Message
}

type Worker struct {
	ID         int
	Client     Client
	WorkerPool chan chan Task
	quit       chan bool
}

func NewWorker(id int, wq chan chan Task) Worker {
	return Worker{
		ID:         id,
		WorkerPool: wq,
		quit:       make(chan bool),
	}
}

func (w Worker) Start() {
	log.Printf("Starting email worker %v", w.ID)
	go func() {
		for {
			w.WorkerPool <- TaskQueue

			select {
			case work := <-TaskQueue:
				log.Printf("Load and clear!: %s", work)
			case <-w.quit:
				return
			}
		}
	}()
}

func (w Worker) Stop() {
	log.Println("Stopping email worker %v", w.ID)
	go func() {
		w.quit <- true
	}()
}
