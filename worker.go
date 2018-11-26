package main

/* 
 * Worker represents threads that actually picks 
 * asynchronous tasks like DB updates from dispatcher
 */

import (
    "log"
)

func NewWorker(workerPool chan chan Job) Worker {
    return Worker{
        WorkerPool: workerPool,
        JobChannel: make(chan Job),
        quit:       make(chan bool)}
}

/* Start method starts the run loop for the worker, listening for a 
 * quit channel in case we need to stop it
 */
func (w Worker) Start() {
    go func() {
        for {
            // register the current worker into the worker queue.
            w.WorkerPool <- w.JobChannel

            select {
            case job := <-w.JobChannel:
                // we have received a work request.
                if err := job.WriteToDB(); err != nil {
                    log.Printf("Error updating in DB: %s", err.Error())
                }

            case <-w.quit:
                // we have received a signal to stop
                return
            }
        }
    }()
}


/* Stop signals the worker to stop listening for work requests.
 */
func (w Worker) Stop() {
    go func() {
        w.quit <- true
    }()
}


