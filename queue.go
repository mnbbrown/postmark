package postmark

type task struct {
	Message *Message
}

var taskQueue = make(chan task, 100)

// Dispatcher is a worker queue
type Dispatcher struct {
	workerPool chan chan task
	Client     *Client
}

// NewDispatcher creates a new dispater with workers
func NewDispatcher(maxworkers int, client *Client) *Dispatcher {
	pool := make(chan chan task, maxworkers)
	return &Dispatcher{pool, client}
}

// Add adds a new task to the queue
func (d *Dispatcher) Add(m *Message) {
	taskQueue <- task{m}
}

// Run starts the dispatcher
func (d *Dispatcher) Run() {
	for i := 0; i < cap(d.workerPool); i++ {
		worker := newWorker(i, d.workerPool, d.Client)
		worker.Start()
	}
	go d.dispatch()
}

func (d *Dispatcher) dispatch() {
	for {
		select {
		case ts := <-taskQueue:
			go func(t task) {
				taskChannel := <-d.workerPool
				taskChannel <- t
			}(ts)
		}
	}
}

type worker struct {
	ID         int
	Client     *Client
	workerPool chan chan task
	quit       chan bool
}

func newWorker(id int, wq chan chan task, client *Client) worker {
	return worker{
		ID:         id,
		workerPool: wq,
		Client:     client,
		quit:       make(chan bool),
	}
}

func (w worker) Start() {
	go func() {
		for {
			w.workerPool <- taskQueue

			select {
			case task := <-taskQueue:
				w.Client.SendMessage(task.Message)
			case <-w.quit:
				return
			}
		}
	}()
}

func (w worker) Stop() {
	go func() {
		w.quit <- true
	}()
}
