package main

/*
 * This middleware will ensure that every request will
 * get logged.
 */

import (
    "log"
    "time"
    "net/http"
)

func LoggingMiddleware(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

    /* First log will ensure that transactions causing server crash, if any,
     * gets logged
     */
    t := time.Now()
    log.Printf("%v Received\t- \"%v %v %v\"", r.Host, r.Method, r.URL, r.Proto)

    /* Call the next handler function - route() in our case */
    next.ServeHTTP(w, r)

    /* Final log giving latency
     */
    latency := time.Since(t)
    log.Printf("%v Responded\t- \"%v %v %v\", with RespTime = %vµs", r.Host, r.Method, r.URL, r.Proto, int((latency.Nanoseconds())/10000))
  })
}


 
