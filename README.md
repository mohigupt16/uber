# 'Find My Driver' 



##    TABLE OF CONTENTS
1.    Introduction
2.    Why Go?
3.    Problem Statement
4.    Infrastructure Requirements
5.    Go language Installation 
6.    Application Installation
7.    Compilation and Run 
8.    Automated Tests
9.    Test using curl (You can also test using tools like Postman!)
10.    Design and Approach
11.    NOTES Regarding This Activity


##    Introduction
This application runs a HTTP Server in Gooogle's Go language.
Its a sample applciation to mimic a server that provides UBER like services
to connected drivers and customers.

I have recently started exploring one component written in Go language.
To understand that, I decided to try out a sample web framework. This 
is my first attempt to develop a web application with Restful Interface.
Due to the benefits mentioned below, I quickly went through Go language 
docs for first few days and then started exploring the design for 
this sample application.
Finally implemented the code in Go-lang.



##    Why Go?
Go has developed a good popularity in recent years -
1. It exhibits simple build interface and low build time requirement during application 
maintainence and releases. 
2. It provides quite rich packages for developing web applications. 
Also, third party libraries are being actively contributed by open source community 
and integrating these apps is quite simple and straightforward, much of which is 
taken care of by the go build tool.
3. Almost no to very small makefiles are required due to Go's ability to auto
detect the dependencies from code structure.
4. Unlike other web frameworks like Node.js, Go is capable of multi-threading through
its go routines which are easy to implement and provides capabilities at application 
level for parallel processing, thereby taking advantage of multi-core machine architecture.
5. Go provides easy cross platform compilation without any changes in native api interfaces.

Due to above factors, web application is quite fast and easy to maintain across releases.
Go has been quite popular for developing small to big frameworks and has seen wide 
acceptability by community.


##    Problem Statement
For an UBER like company, create a driving application that can be used to 
keep track of drivers current location and provide an ability to search drivers
in a given area. Assume any number as the max number of drivers say 50,000, who may
be active at any point in time. 
We need to build 2 APIs to achieve this, both APIs should respond within 100ms.


### (1) Driver Location
Drivers should be able to send their current location every 60 seconds. They’ll call following
API to update their location
Expected Request:

```

PUT /drivers/{id}/location
{
  "latitude": 12.97161923,
    "longitude": 77.59463452,
    "accuracy": 0.7
}

```


Expected Respnose:

```

- 200 OK on successful update
Body: {}
- 404 Not Found if the driver ID is invalid (valid driver ids - 1 to 50000)
Body: {}
- 422 Unprocessable Entity - with appropriate message. For example:
{"errors": ["Latitude should be between +/- 90"]}

```


Load Assumption:
50,000 requetss per 60 sec.


### (2) Find Drivers
Customer applications will use following API to find drivers around a given location
Request:

```

GET /drivers
Parameters:
"latitude" - mandatory
"longitude" - mandatory
"radius" - optional defaults to 500 meters
"limit" - optional defaults to 10

```

Response:

```

- 200 OK
[
{id: 42, latitude: 12.97161923, longitude: 77.59463452, distance: 123},
{id: 84, latitude: 12.97161923, longitude: 77.59463452, distance: 123}
]
- 400 Bad Request - If the parameters are wrong
{"errors": ["Latitude should be between +/- 90"]}
Distance in the response is a straight line distance between driver's location and location in
the query

```

Load Assumption : 
20 concurrent requests 



##  Infrastructure Requirements
As stated above, Go application are cross-platform compilable.
System Requirements - 
>   OS - Linux 2.6.23 or later with glibc   
>   Architecture - amd64, 386, arm    
>   Note - CentOS/RHEL 5.x are not supported


##  Go language Installation 
*Note* : 'root' or 'sudo' privledges required
1. Download archive : 
```
$ wget https://storage.googleapis.com/golang/go1.7.4.linux-amd64.tar.gz
```
2. Extract it into /usr/local, creating a Go tree in /usr/local/go
```
$ tar -C /usr/local -xzf go1.7.4.linux-amd64.tar.gz
```
3. Add /usr/local/go/bin to the PATH environment variable (or update path in $HOME/.profile and source)    
```
$ export PATH=$PATH:/usr/local/go/bin
```
4. Create a directory to contain your workspace, $HOME/work for example, and set the GOPATH environment variable 
```
$ mkdir $HOME/work
$ export GOPATH=$HOME/work      #Add it to ~/.bash_profile and source 
```





##    Application Installation
1. Goto GO's workspace directory as set above
```
$ cd $GOPATH
```

2. Create `src` directory and `cd` into it
```
$ mkdir src && cd src
```

3. Clone this project
```
$ go get github.com/mohigupt16/uber
```

4. You should see following dir structure containing .go source code files (This is GO's convention of code layout) -  
```
src/github.com/mohigupt16/uber/
```

    

##    Compilation and Run 
1. Goto src directory created above  
```
$ cd $GOPATH/src/github.com/mohigupt16/uber
```
2. Compile :
```
$ go install
```
   This will create a bin directory in $GOPATH (parallel to src dir) containing the application binary
```
$ ls $GOPATH/bin/uber
```
3. Run :
```
$ $GOPATH/bin/uber
```
  This will start the server at - localhost:8080

Note: Alternatively, you can also build from path $GOPATH/src by explicitly specifying the project to be built
```
$ go install github.com/mohigupt16/uber
```



##    Automated Tests
Go provides an inbuilt functionality to execute the test cases defined within src code directory.
The nomenclature is that src directory may contain the tests in one more multiple \*\_test.go files
which are ignored by compiler for building the code. 
We have defined our test cases in main\_test.go file in the src directory and these can be automated
by running following command -

1. Goto src code directory as created above
```
$ cd $GOPATH/src/github.com/mohigupt16/uber
```

2. Run tests
```
$ go test
```


If all tests get passed, PASS would be dumped on screen with no error
In case of failure, appropriate failure statement would be reported with file's line number
Note that failed cases are reported after excuting all defined test cases in \*\_test.go files.

*Note(1)*: Test automation will create a metadata file $GOPATH/src/github.com/mohigupt16/uber/test\_metadata.txt 
          This contains the dump of the cahced driver data which was inserted in IN_MEMORY DB during the
          tests and gets appended upon every subsequent test.

*Note(2)* : Alternatively, you can also test from path $GOPATH/src by explicitly specifying the project to run
          automated tests. In this case, look for status 'ok' in case of success. Failure case is same.  
```
$ go test github.com/mohigupt16/uber
```


Note(3) : You can change test data values inside the main\_test.go file to some invalid values and re-run 
above 'go test' cmd to see how errors get reported.


##    Test using curl
Two API(s) are supported -
1. GET /drivers
2. PUT /driver/{id}/location

Create some enteries using curl. For ex :
```
$   curl -XPUT "localhost:8080/drivers/12/location" -d '{"latitude":12.97161923,"longitude":77.59463452,"accuracy":0.7}'
$   curl -XPUT "localhost:8080/drivers/12/location" -d '{"latitude":11.97161923,"longitude":76.59463452,"accuracy":0.7}'
$   curl -XPUT "localhost:8080/drivers/12/location" -d '{"latitude":10.97161923,"longitude":75.59463452,"accuracy":0.7}'
```
Get drivers for specified params. For ex :
```
$   curl -XGET "localhost:8080/drivers?latitude=12&longitude=77&radius=200000&limit=2"
```

*Note* : You can also test using tools like Postman!



##    Design and Approach
All src files and their functions carry proper headers to understand their definition and code flow.

In summary, there are following src files containing entire code and the flow is also explained per file -

main.go - Contains the main function for intializing database and middlewares and specfying route handlers for
desired endpoints. It then starts the HTTP server at "localhost:8080"


### [models.go](models.go) -
Defines the structs used for various data models used for extrcating parameters from incoming
requests, storing the data in database store and sending the params in response body.
Also defines the constants and config params that are referanced through out the code.
*NOTE* For doing Configuration changes for this application, please refer to "config" const declarations
We could move these to read from a config file at run time in a production like scenario.


### [loggingMiddleware.go](loggingMiddleware.go) - 
Defines the middleware functions that encloses the route handler functions so that pre-
and post-processing code can be added to every received request. We use it to log every incoming request.
Currently, the log goes to stdout.



### [router.go](router.go) - 
Defines the two desired endpoints and do an exact path binding(except for training '/') with route
handlers using go's regexp package.



### [handlers.go](handlers.go) - 
Defines request handlers for the supported requests.
1. In case of `PUT` requests, we delayed the storing of data by delegating the record(after validation) to 
a dispatcher which is infinitely listening on a queue, JobQueue, for incoming records as events.
The disptacher then blocks on WorkerPool until it gets one worker thread that can execute this event(job).
On finding one such worker, delegates this record to that worker's queue, jobChannel.
Each worker, upon receiving an event in their respective queue, jobChannel, executes it and then blocks on
jobChannel until they receive another event from dispatacher.

2. In case of `GET` requests, we have implemented only IN\_MEMORY store (which is a map) for storing the driver
details. The code has been written in such a way that any new store can be added as a plugin and can be
specified from config. 
This request is synchronously processed where in case of IN\_MEMORY records, all records are iterated till
"limit" count of nearest drivers have been found within specified radius which is assumed in meters 
(we do not attempt to find "top" nearest).


3. *Optimized Search* : In case of external DBs like mysql, we can make a smarter approach to narrow our scan of 
records in DB and so as to decrease the IO. 
We can do this by determining min/max lat/lon values based on received coordinates
This would return a rectangle of diameter=2\*radius around given coordinate
and contains nearest drivers.
We can then evaluate each entry in this result set to find out the drivers
that fall within this circle of received corrdinates as below
For this, we should create two more tables in DB, one for {lat,id} and other for
{lon, id}. The commod ids can then be fetched from main table which contains
denormalised data(i.e. contain entire details with id)


4. *Cache* : Also, we can maintain a cache of all stored drivers and update it every 60 sec (say 10sec
post PUT updates) after all driver updates have been received.



5. [validators.go](validators.go) - 
Defines functions for validating the parameters received with GET/PUT requests and 
extracting them and return to handler functions for processing.



6. [dispatcher.go](dispatcher.go) - 
It aims to provide a framework for asynchronous processing of I/O intensive part of 
the received PUT requests. HTTP response is sent as soon as the validation is passed. Thereafter,
update in DB store is done as explained above. 
*NOTE*  We can not directly use the route handler to delegate the request to worker as all 
workers may be busy at some point and so the request handler will get blocked. 
Dispatcher removes this blocking call.



7. [worker.go](worker.go) - 
It helps in parallelizing the asynchronous processing of incoming requests by spawning
multiple threads, each representing a worker. Worker will do th I/O intensive operations thereby
allowing server to scale to millions of requests per minute.



8. [dbWrapper.go](dbWrapper.go) - 
It is responsible for receiving read or write requests from the handlers/workers and makes
it opaque to the underlying DB store as specified in configuration. Plan was to provide 2 implementations
\- for in memory(RAM) storage and mysql storage - but currently mysql part has not been
implemented. However, code structure has been kept so that it can easily plugged.




##    *NOTES* \- Regarding This Activity
It should be noted that there are many nice open packages(libraries) which are now part of official
GO Docs but have been maintained as seperate github libraries and DONOT come as part of official 
go installation. These packages help a lot in writing even simpler and less code. But I have not
used them because of the requirement NOT to use any external library so resorted to only those that
come with official go installation.
These libraries make it much simpler to log requests, extract params like ResponseCode while logging
the incoming traffic, managing and accesing HTTP Requests and Responses (like go-gin), etc



####    IGNORE THIS PART
what is done?
1. server created
2. Get and Put handling done, params validated and extracted
3. Create framework to process requests asynchronously(for PUT request processing)
4. Add processing of GET requests and handle all error codes
6. Add log middleware for tracking every request
8. Create ReadMe
7. Add test framework and automation script for installing and running application






















