package main

/*
 * Routes definitions
 */

import (
    "net/http"
    "regexp"
)

/* Regexes for acceptable endpoints. Optionally allows the trailing '/' */
var rGetDriv = regexp.MustCompile(`^/drivers$(/?)$`)                // GET /drivers
var rPutDriv = regexp.MustCompile(`^/drivers/\d+/location(/?)$`)    // PUT /drivers/{id}/location

/* Routes all acceptable endpoints to their repective handlers
 * Inputs :
 *      w - writer for response
 *      err - error message as string
 * Returns :
 *      None
 */      
func route(w http.ResponseWriter, r *http.Request) {
    switch {
        case rPutDriv.MatchString(r.URL.Path):
                putDriver(w, r)
                return
        case rGetDriv.MatchString(r.URL.Path):
                getDrivers(w, r)
                return
        default:   
                http.Error(w, "Bad Request - Resource Unknown!", 404)
    }
}
