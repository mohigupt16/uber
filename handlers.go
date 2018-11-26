package main

/* 
 * HTTP Handler functions for endpoints
 */

import (
    "fmt"
    "log"
    "net/http"
    "encoding/json"
    //"strconv"
)

/* Sets the received error string in http response writer as Json
 * Inputs :
 *      w - writer for response
 *      err - error message as string
 * Returns :
 *      None
 */      
func setHttpErrorWithJson(w http.ResponseWriter, err string, errCode int) {
    /* convert the error message to JSON for reply */
    msg, e := json.Marshal(&makeError{Mesg : []string{err}})
    if e != nil {
        log.Printf("Error: %s", e)
    }
    http.Error(w, string(msg), errCode)
}


/* Http Handler for 'GET /drivers' 
 * Inputs :
 *      w - writer for response
 *      r - HTTP request object
 * Returns :
 *      None
 */
func getDrivers(w http.ResponseWriter, r *http.Request) {

    /* Validate query parameters as per given requirement */
    vs, errStr, errCode := validateParams(r, GetDrivers)
    if len(errStr) > 0 {
        setHttpErrorWithJson(w, errStr, errCode)
        return
    }

    /* READ Latency : Since reads are comparatively faster than
     * writes, we are directly quering the DB for this request
     * Reads are easy to optimize with latest DBs(sql/nosql) 
     * as they offer parallel node reads using sharding
     * and replicas
     *
     * Still, if unacceptable latency is observed during
     * parallel requests, we can built a cache which can
     * be updated, say 10sec, after updating DBs as
     * all Driver Update requests should get complete
     * by this time. 
     */
    results, errStr, errCode := getNearestDrivers(vs)
    if len(errStr) > 0 {
        setHttpErrorWithJson(w, errStr, errCode)
        return
    }

    /* convert results into JSON before sending in response */
    var resp []string 
    for _, r := range results {
        val, err := json.Marshal(r)
        if err != nil {
            log.Printf("Error: %s", err)
            setHttpErrorWithJson(w, "Internal error", 500)
            return;
        }
        resp = append(resp, string(val))
    }

    fmt.Fprint(w, resp) 
}


/* Http Handler for 'PUT /drivers/{id}/location' requests 
 * Inputs :
 *      w - writer for response
 *      r - HTTP request object
 * Returns :
 *      None
 */
func putDriver(w http.ResponseWriter, r *http.Request) {

    /* Validate params */
    vs, errStr, errCode := validateParams(r, PutDriver)
    if len(errStr) > 0 {
        setHttpErrorWithJson(w, errStr, errCode)
        return
    }

    /* Send request to dispatcher on channel that dispatcher is listening to
     */
    payload := DriverStore{Id: vs["id"], Latitude: vs["lat"], Longitude: vs["lon"], AccOrDist: vs["acc"]}
    work := Job{Payload: payload} 
    select {
        case JobQueue <- work :
            //fmt.Printf("Successfully dispatched %v\n", work)
        default :
            /* this case ensures that this thread does not block on JobQueue in case its full to its capacity
             */
            setHttpErrorWithJson(w, "Server overloaded. Try after sometime.", 513)
            return
    }
}


