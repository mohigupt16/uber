package main

/* 
 * Dispatcher represents thread that receives IO intensive
 * requests(that can be processed asynchronously)
 * for offline DB updates, etc.
 * It delegates these requests to one of the workers
 */ 


func NewDispatcher(count int) *Dispatcher {
    pool := make(chan chan Job, count)
    return &Dispatcher{WorkerPool: pool, maxWorkers: count}
}

func (d *Dispatcher) Run() {
    // starting n number of workers
    for i := 0; i < d.maxWorkers; i++ {
        worker := NewWorker(d.WorkerPool)
        worker.Start()
    }

    go d.dispatch()     //spawns a new thread with this routine
}

func (d *Dispatcher) dispatch() {
    for {
        select {
        case job := <-JobQueue:
            // a job request has been received
            go func(job Job) {
                // try to obtain a worker job channel that is available.
                // this will block until a worker is idle
                jobChannel := <-d.WorkerPool

                // dispatch the job to the worker job channel
                jobChannel <- job
            }(job)
        }
    }
}
