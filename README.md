
## Clinic Search Service

##### Submission by Ademola Oyewale (saopayne@gmail.com)

###### Submission Date: 03-06-2021

#### Project Overview

Provide a RESTful API to allow search in multiple clinic providers and display results from all the available clinics;

- search for clinics across dental and vet based on specified params;
- get all the clinics;

#### Tools Used in the project

- GO
- Postman Collections (API Documentation)

#### Project Architecture & Design

I split the package structure into three components majorly (cmd, pkg & internal)

#### Concurrency Approach

I used Goroutines and WaitGroup to fetch data from the two endpoints concurrently. If we do have additional endpoints
to fetch data from, we can further abstract the methods to take in arbitrary number of URLs with specific parsers
which will enable the generic parser to collect the data into a uniform interface.

#### Running up the application

###### A. Via Docker

``` 
$ make build
$ make run
$ curl -X GET http://0.0.0.0:8000/v1/clinics/search
```

###### B. Via Makefile

``` 
$ make build 
$ ./bin/linux.amd64/main
```

###### C. Locally

```
$ cd scratchpay_ademola 
$ go run ./cmd/*  
```

Application is now available `http://localhost:8000/`

#### Running Tests

Running tests `$ make test`

#### Endpoints

I created two endpoints:

- `GET: /v1/clinics/`: This returns all the clinics from both endpoints
```json
$ curl -X GET http://0.0.0.0:8000/v1/clinics
>> 
[
    {
        "name":"Good Health Home",
        "state":"FL",
        "availability":{
            "from":"15:00",
            "to":"20:00"
        }
    },
    {
        "name":"National Veterinary Clinic",
        "state":"CA",
        "availability":{
            "from":"15:00",
            "to":"22:30"
        }
    },
    {
        "name":"German Pets Clinics",
        "state":"KS",
        "availability":{
            "from":"08:00",
            "to":"20:00"
        }
    }
]

```

- `POST: /v1/clinics/search`: This enables searching of the clinics based on the specified params
Using:
```json
$ curl -d '{"name":"German"}' -H "Content-Type: application/json" -X POST http://0.0.0.0:8000/v1/clinics/search
>>
[
  {
        "name":"German Pets Clinics",
        "state":"KS",
        "availability":{
            "from":"08:00",
            "to":"20:00"
        }
    }
]
```

#### Documentation

I have included two files in the base directory of the project;

- (open-api.yaml) Open-api which can be rendered here: [Online Swagger Editor](https://editor.swagger.io/)

#### Challenges

- NA

#### Further Improvements

- Include pagination for the Get all clinics endpoint.
- Abstract the URL passing mode for fetching data to accommodate additional URLs that might pop up in the future.
- Implement retry when calling external APIs along with exponential backoff.
- Authentication and Authorization for adequate security.
- Improve monitoring of the API to encourage pro-active instead of reactive behavior.
