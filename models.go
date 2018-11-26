package main

/*
 * Data Models for this application
 */


/* Schema for receiving driver updates from 'Put /drivers/{id}/location' requests */
type DriverUpdates struct {
    Latitude  float64   `json:"latitude"`
    Longitude float64   `json:"longitude"`
    Accuracy  float64   `json:"accuracy"`
}

/* Schema for creating responses to 'GET /drivers'
 */
type DriverResp struct {
    Id          int     `json:"id"`
    Latitude    string  `json:"latitude"`
    Longitude   string  `json:"longitude"`
    Distance    string  `json:"dist"`
}
        
/* Schema for receving request params in 'GET /drivers' 
 */
type DriverGet struct {
    Latitude    float64  `json:"latitude"`
    Longitude   float64  `json:"longitude"`
    Radius      float64  `json:"radius"`
    Limit       float64  `json:"limit"`
}
        

/* Schema for storing driver details 
 */
type DriverStore struct {
    Id          float64  `json:"id"`
    Latitude    float64  `json:"latitude"`
    Longitude   float64  `json:"longitude"`
    AccOrDist   float64  `json:"distance"`  //using this field to store accurracy while writing in DB and
                                            //as distance from provided coordinates while responding to 
                                            //nearestDriver request
}

/* Helper struct for converting a string error message into a json 
 */ 
type makeError struct { 
    Mesg    []string  `json:"errors"`
}

/* Enum Simulation as enums are not available in Go 
 */
type DrivApis int
const (
    GetDrivers = 1  
    PutDriver  = 2
)

/* Enum Simulation as enums are not available in Go 
 */
type DBStores int
const (
    STORE_IN_MEMORY= 10 
    STORE_MYSQL     = 11
)

/* Configuration params for this application
 */
type config int
const (
    /* Driver Attribute Defaults */
    MIN_DRIVER_ID = 1
    MAX_DRIVER_ID = 50000
    RADIUS = 500
    LIMIT  = 10

    /* Worker/Dispatcher Defaults */
    MAX_WORKERS = 4
    MAX_QUEUE = 50

    /* Selected DB type */
    CURRENT_DB = STORE_IN_MEMORY
)



/* Job represents the job to be done like update Driver details in DB 
 * or extract results from DB to update cache  */
type Job struct {
    Payload DriverStore
}


/* Worker represents the thread that executes the job
 * We should keep worker count to max number of cores
 * available across machines */
type Worker struct {
    WorkerPool  chan chan Job
    JobChannel  chan Job
    quit        chan bool
}

/* Dispatcher is responsible for asynchronously delegating received requests processing
 * as a Job to one of the workers
 */
type Dispatcher struct {
    // A pool of workers channels that are registered with the dispatcher
    WorkerPool chan chan Job
    maxWorkers int
}

/* A buffered channel that we can send work requests on 
 * Capacity is set to twice the support required but can
 * be kept at some other smaller multiple if MAX_DRIVER_ID
 * is too high
 */
var JobQueue = make(chan Job, 2*MAX_DRIVER_ID)

var inMemDb map[float64]DriverStore


/* 
 * TODO: We should rather make key as an int constant and value as interface{}
 * to allow it to take up any type of values.
 * This will alleviate from keeping same structure for return types of
 * validateParams() while still keeping the same design in place.
 *
 * We will know during performance testing if this change is required
 */
type Values map[string]float64


/* creates a map entry
 */
func (v Values) Add(key string, value float64) {
    v[key] = value
}

