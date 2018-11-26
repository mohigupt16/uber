package main

/* 
 * Defines test cases for unit testing
 * Though the tests should be created at function level 
 * to evalute every function for both positive and negative
 * values for arguments, but we are here defining minimal
 * tests that simulates creating the data structures for testing
 * out 3 functionalities of the main code which are reading/writing
 * in DB, validation func for given Get requests and validation for
 * given Put requests.
 *
 * Test in go can be automated by running - go test
 * which will execute all functions in this file that start with
 * TestXxxx or Test_xxxx
 *
 * To automate tests, execute - $ go test
 */

import (
    "testing"
    "net/http"
    "net/url"
    "io/ioutil"
    "strings"
)


/* Tests our code for reading and writing the driver data in DB 
 */
func Test_write_and_read_db(t *testing.T) {
    /* initialize DB */
    initDB()
    defer closeDB("test_metadata.txt")      //dumps the Driver data from memory into this file after executing
                                            //the test; defer is GO's way of executing this func() after
                                            //current func() has exited

    /* put some drivers in DB */
    var drivers = []DriverStore{
                        {Id: 1234, Latitude: 12.97161923, Longitude: 77.59463452, AccOrDist: 0.7},
                        {Id: 6547, Latitude: 12.96161923, Longitude: 77.58463452, AccOrDist: 0.8},
                        {Id: 1234, Latitude: 10.97161923, Longitude: 75.59463452, AccOrDist: 0.9},
                     }
    for _, work := range drivers {
        job := Job{Payload: work}
        if err := job.WriteToDB(); err != nil {
            t.Error("Expected nil, got ", err)
        }
    }
    
    /* try finding nearest driver */
    var clients = []Values{
                        {"lat":12,"lon":77,"rad":200000,"lim":1},       //rad is in meters, 1 match
                        {"lat":12,"lon":77,"rad":200000,"lim":2},       //2 matches
                        {"lat":12,"lon":77,"rad":200000,"lim":4},       //2 matches
                        {"lat":12,"lon":77,"rad":1000,"lim":1},         //0 match
                    }
    for i, c := range clients {
        results, errStr, _ := getNearestDrivers(c)
        if len(errStr) > 0 {
            t.Error("Expected nil, got error ", errStr)
            continue
        }
        if i==0 && len(results)!=1 {
            t.Error("Expected count 1, got ", results)
        }
        if i==1 && len(results)!=2 {
            t.Error("Expected count 2, got ", results)
        }
        if i==2 && len(results)!=2 {
            t.Error("Expected count 2, got ", results)
        }
        if i==3 && len(results)!=0 {
            t.Error("Expected count 0, got ", results)
        }
    }
}


/* Tests the validation function for GET requests
 */
func Test_validate_get_request(t *testing.T) {
    var r http.Request
    var u url.URL
    r.Method = "GET"
    r.URL = &u

    queries := []string {
                    "latitude=12&longitude=77&radius=200000&limit=2",       //valid
                    "latitude=220&longitude=77&radius=200000&limit=2",      //invalid
                    "latitude=12&longitude=-777&radius=200000&limit=2",     //invalid
                    "latitude=12&longitude=77&radius=-200&limit=2",         //invalid
                    "latitude=12&longitude=77&radius=18000000&limit=2",     //invalid
                    "latitude=12&longitude=77&radius=200000&limit=51000",   //invalid
                }

    for i, q := range queries {
        r.URL.RawQuery = q 
        _, errStr, _ := validateGetDriverParams(&r)
        if i==0 && len(errStr) > 0 {
            t.Error("Expected nil for GetRequest, got error - ", errStr)
        } 
        if i > 0 && len(errStr) == 0 {
            t.Error("Expected error for GetRequest number -  ", i)
        }
    }
}


/* Tests the validation function for PUT requests
 */
func Test_validate_put_request(t *testing.T) {

    var r http.Request
    var u url.URL
    r.Method = "PUT"
    r.URL = &u
    r.URL.Path = "/drivers/123/location"

    putData := []string {
                    `{ "latitude": 12.97161923, "longitude":  77.59463452, "accuracy": 0.7 }`,     //valid                   
                    `{ "latitude":212.97161923, "longitude":  77.59463452, "accuracy": 0.7 }`,     //invalid
                    `{ "latitude": 12.97161923, "longitude":-777.59463452, "accuracy": 0.7 }`,     //invalid
                    `{ "latitude": 12.97161923, "longitude":  77.59463452, "accuracy": 1.7 }`,     //invalid
                }

    for i, d := range putData {
        r.Body = ioutil.NopCloser(strings.NewReader(d))

        _, errStr, _ := validatePutDriverParams(&r)
        if i==0 && len(errStr) > 0 {
            t.Error("Expected nil for PutRequest, got error - ", errStr)
        } 
        if i > 0 && len(errStr) == 0 {
            t.Error("Expected error for PutRequest number -  ", i)
        }
    }
}


