package main

/*
 * DB Wrapper that provides interface between the application
 * and underlying database. This Wrapper should ensure that
 * one can call same api(s) irrespective of the configured DB
 * (like mysql or in memory, etc). The wrapper should translate
 * to exact read and write call for configured DB.
 */

import (
    "os"
    "log"
    "math"
    "encoding/json"
)

/* Initialises the configured database
 * Returns true in case of sucessful init, else returns false
 */
func initDB() bool {
    switch CURRENT_DB {
        case STORE_IN_MEMORY :
            inMemDb = make(map[float64]DriverStore)
            return true
        case STORE_MYSQL :
            //TODO:
            log.Println("Exiting: configured DB not supported!")
            return false
        default :
            log.Println("Exiting: configured DB not supported!")
            return false
    }    
}   


/* For graceful shutdown of DB in case of stopping the application
 * 
 * Caveat : We can put logic in initDB in case of IN_MEM DBs to persist data
 * ------
 * while stopping the application so that it can be re-read
 * during application restart
 * Not doing init() with reading of data from disk as its not in scope of current activity
 *
 * But showing the writing of data for demo.
 */
func closeDB(f string) {
    switch CURRENT_DB {
        case STORE_IN_MEMORY :
            // Write to disk to sustain last populated values
            file, err := os.OpenFile(string(f), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
            if err != nil {
                log.Printf("Error in persisting data to file %v, got error while opening - %v", f, err)
                return
            }
            for _, v := range inMemDb {
                b, err := json.Marshal(v)
                if _, err = file.Write(b); err != nil {
                    log.Printf("Error in persisting Data to file %v, got - ", f, err)
                    file.Close()
                    return
                }
            }
            if err := file.Close(); err != nil {
                log.Printf("Error in persisting Data to file %v, got - %v", f, err)
            }
            return 
        case STORE_MYSQL :
            //TODO:
            log.Println("Exiting: configured DB not supported!")
            return 
        default :
            log.Println("Exiting: configured DB not supported!")
            return
    }    
}

/* Wrapper for extracting nearest drivers from DB for HTTP request coordinates.
 * It calls appropriate DB wrapper based on configured DB type.
 * Inputs :
 *      r - Htpt request object
 *      api - represents called request
 * Returns :
 *      []DriverStore - array containing list of nearest drivers found
 *      string - contains the error message in case of any failure
 *      int - HTTP error code
 */
func getNearestDrivers(v Values) ([]DriverStore, string, int) {

    /* We should avoid making full scan of DB for specified coordinates
     * We can do this by determining min/max lat/lon values based on received coordinates
     * This would return a rectangle of diameter=2*radius around given coordinate 
     * and contains nearest drivers.
     * We can then evaluate each entry in this result set to find out the drivers
     * that fall within this circle of received corrdinates as below
     * For this, we should create two more tables in DB, one for {lat,id} and other for 
     * {lon, id}. The commod ids can then be fetched from main table which contains
     * denormalised data(i.e. contain entire details with id)
     *
     * This implementation is good when data is saved in some external DB(like mySql)
     * Sample Implementation :
     * minLat, maxLat, minLon, maxLon, err := getRangeOfCoordinates(v)
     */

    switch {
        case CURRENT_DB == STORE_IN_MEMORY :
            return getNearestDriversFromInMemStore(v)
        default :
            return nil, "Internal Error!", 500          // we can not mark it 4xx because its our server
                                                        // error to NOT properly set the CURRENT_DB, 
                                                        // it indicates that initDB() messed up
    }
}

/* Wrapper for extracting nearest drivers from STORE_IN_MEMORY DB for HTTP request coordinates.
 * The fucntion iterates through all entries and calculate distance in meters with its own coordinates
 * If returned distance is less than provided radius value, it is put into a sorted set
 * A count is kept for found results and we return if the limit is met even if full scannign is not done.
 *
 * Caveat : Since requirement specifies only max "lim" coordinates within the provided radius, we are
 * not binded to find the top nearest drivers.
 *
 * Input and Return values are same as described for getNearestDrivers()
 */
func getNearestDriversFromInMemStore(v Values) ([]DriverStore, string, int) {

   /* TODO : we should handle the accuracy also by considering two values
    * for each lat - min(=acc*lat) and max(=(1-acc)*lat). Same for lon. 
    * This is because accuracy indicates that driver may be present anywhere
    * within latitude/longitude range {min, max} as calculated above
    */

    la := v["lat"]
    lo := v["lon"]
    ra := v["rad"]
    li := (int)(v["lim"])
    cnt := 0
    dis := 0.0
    var s []DriverStore
    var d DriverStore
    for _, value := range inMemDb {
        d = value
        dis = Distance(la, lo, d.Latitude, d.Longitude)     //see above note-TODO for incorporating accuracy
        if dis <= ra {
            d.AccOrDist = dis        // reusing this field to return distance calculated
            s = append(s, d)
            cnt++
        }
        if cnt == li {
            return s, "", 200
        }
    }
    return s, "", 200
}

func Distance(lat1, lon1, lat2, lon2 float64) float64 {
    // convert to radians
    // must cast radius as float to multiply later
    var la1, lo1, la2, lo2, r float64
    la1 = lat1 * math.Pi / 180
    lo1 = lon1 * math.Pi / 180
    la2 = lat2 * math.Pi / 180
    lo2 = lon2 * math.Pi / 180

    r = 6378100 // Earth radius in METERS

    // calculate
    h := hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*hsin(lo2-lo1)

    return 2 * r * math.Asin(math.Sqrt(h))
}

func hsin(theta float64) float64 {
    return math.Pow(math.Sin(theta/2), 2)
}

/* 
 * It calculates min amd max values for specified coordinates for 
 * shortening the scan in external DBs like mysql.
 * It is not meaningful in case of STORE_IN_MEMORY storage type.
 * 
 * Assumption : radius is already in radians (should have been converted during validation step)
 */
func getRangeOfCoordinates(v Values) (float64, float64, float64, float64){
    lat := v["lat"]
    lon := v["lon"]
    rad := v["rad"]

    min_lat := lat - rad
    max_lat := lat + rad
    min_lon := lon - rad
    max_lon := lon + rad

    return min_lat, max_lat, min_lon, max_lon
}


/* 
 * This function does the writing to configured DB for received
 * record(Job)
 */
func (v Job) WriteToDB() error {
    switch {
        case CURRENT_DB == STORE_IN_MEMORY :
            inMemDb[v.Payload.Id] = v.Payload
        default :
            return  nil   
    }
    return nil
}


