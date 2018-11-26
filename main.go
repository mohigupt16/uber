package main

/*
 * Driver applicaion assignment for Go-JEK
 */

import (
    "log"
    "os"
    "net/http"
)

func main() {

   /* create dispatcher and workers framework to process received 
    * requests asynchronously
    * We should aim to respond to HTTP requests as soon
    * as possible. I/O intensve operations should be 
    * delegated to worker threads for asynchronous 
    * processing
    * 
    * Caveat : Any failure in I/O update represents server's 
    * ------
    * internal failure and client, if conifgured to receive server
    * action reports, should be informed about it through
    * seperate request. This functionality is currently
    * not implemented in scope of current activity
    */
    dispatcher := NewDispatcher(MAX_WORKERS)

    /* call Run() on dispacher to start listening for 
     * incoming requests. All delayed processings should
     * be sent to dispatcher by route handlers for
     * delegating to one of the workers
     */
    dispatcher.Run()

    /* initialise DB based on configuration in models.go
     */
    if ok := initDB(); !ok {
        os.Exit(1)
    }

    /* Register the endpoints to be supported.
     * LoggingMiddleware is a middleware that encloses
     * every request handler method
     */
    routeHandler := http.HandlerFunc(route)
    http.Handle("/", LoggingMiddleware(routeHandler))

    /* should be the last line to start the http server */
    log.Println("Starting listening on http://localhost:8080 ...")
    http.ListenAndServe(":8080", nil)
}


