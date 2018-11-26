package main

/*
 * Validators for various HTTP requests defined
 */

import (
    "net/http"
    "strconv"
    "strings"
    "encoding/json"
    "log"
)


/* Wrapper for calling exact Validator for the HTTP request based on called request.
 * Inputs :
 *      r - Htpt request object
 *      api - represents called request
 * Returns : 
 *      url.Values - map[string][]string containing query params as keys
 *                   and their corresponding values in a string array
 *      string - contains the error message in case of validation failure
 *      int - HTTP error code
 *
 * Caveat - we only validate the intended params and then the first value 
 * associated with that param. If more params are present or the intended
 * params carry multiple values, we ignore them in this implementation
 */
func validateParams(r *http.Request, api DrivApis) (Values, string, int)  {
    switch api {
        case GetDrivers:
            return validateGetDriverParams(r)
            
        case PutDriver:
            return validatePutDriverParams(r)

        default:
            return nil, "api not implemented", 404
    }
}


/* Validator for 'PUT /driver'
 * Details as specified in validateParams
 */
func validatePutDriverParams(r *http.Request) (Values, string, int)  {

    /* allow only PUT method for this resource */
    if r.Method != "PUT" {
        return nil, "Method not allowed for requested page", 405
    }

    /* Get parameters for extracting Driver ID */
    uriSegments := strings.Split(r.URL.Path, "/")

    driverId, err := strconv.ParseUint(uriSegments[2], 10, 64)       
    if err != nil {
        return nil, "Invalid driverId type", 400
    }

    /* Decode the provided fields for Driver */
    decoder := json.NewDecoder(r.Body)
    var t DriverUpdates   
    err = decoder.Decode(&t)
    if err != nil {
        log.Println(err)
        return nil, "Request Body format not valid", 422
    }
    defer r.Body.Close()

    if driverId < MIN_DRIVER_ID || driverId > MAX_DRIVER_ID {
        return nil, "DriverID is invalid", 404
    } 
    if t.Latitude > 90 || t.Latitude < -90 {
        return nil, "Latitude should be between +/- 90", 422
    }
    if t.Longitude > 180 || t.Longitude < -180 {
        return nil, "Longitude should be between +/- 180", 422
    }
    if t.Accuracy > 1.0 || t.Accuracy < 0.0 {
        return nil, "Accuracy should be between 0 to 1.0", 422
    }

    vs := make(Values)
    vs.Add("id", float64(driverId))
    vs.Add("lat", t.Latitude)
    vs.Add("lon", t.Longitude)
    vs.Add("acc", t.Accuracy)

    return vs, "", 200      //200 is just a placeholder for our function signatures
}



/* Validator for 'GET /drivers'
 * Details as specified in validateParams
 */
func validateGetDriverParams(r *http.Request) (Values, string, int) {

    /* allow only GET method for this resource */
    if r.Method != "GET" {
        return nil, "Method not allowed for requested page", 405
    }

    vs := r.URL.Query()     //Query() always returns non nil

    var lat,lon,ra,l float64

    if v := vs.Get("latitude"); v == "" {
        return nil, "Mandatory param latitude not specified", 400
    } else {
        f, err := strconv.ParseFloat(v, 64)
        if err != nil {
            return nil, "Invalid latitude type", 400
        } else if f < -90 || f > 90 {
            return nil, "Latitude should be between +/- 90", 400
        }
        lat = f
    }

    if v := vs.Get("longitude"); v == "" {
        return nil, "Mandatory param longitude not specified", 400
    } else {
        f, err := strconv.ParseFloat(v, 64)
        if err != nil {
            return nil, "Invalid longitude type", 400
        } else if f < -180 || f > 180 {
            return nil, "Longitude should be between +/- 180", 400
        }
        lon = f
    }

    if v := vs.Get("radius"); v != "" {
        r, err := strconv.ParseFloat(v, 64)
        if err != nil {
            return nil, "Invalid radius type", 400
        }
        rad := r/(111000)      //convert into radians assuming specified radius is in meters
        if rad <= 0 || (lat-rad) < -90 || (lat+rad) > 90 {
            return nil, "Invalid radius value for specified latitude", 400
        } else if rad <= 0 || (lon-rad) < -180 || (lon+rad) > 180 {
            return nil, "Invalid radius value for specified longitude", 400
        }
        ra = r
    } else {
        ra = RADIUS
    }

    if v := vs.Get("limit"); v != "" {
        lim, err := strconv.ParseUint(v, 10, 64)
        if err != nil {
            return nil, "Invalid limit type", 400
        } else if lim < MIN_DRIVER_ID || lim > MAX_DRIVER_ID {
            s := "Invalid limit value, min " + strconv.FormatUint(MIN_DRIVER_ID, 10) + ", max " +
                    strconv.FormatUint(MAX_DRIVER_ID, 10)
            return nil, s, 400
        }
        l = float64(lim)
    } else {
        l = LIMIT
    } 

    /* All float64s */
    vv := make(Values)
    vv.Add("lat", lat)
    vv.Add("lon", lon)
    vv.Add("rad", ra)
    vv.Add("lim", l)

    return vv, "", 200
}
