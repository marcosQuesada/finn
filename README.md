# Finn Http Client 

Http client library kata exercise.

## Development notes
- Implemented as a http Client library, so, no application project structure, and some default values are hardcoded, as BaseUrl that points to "production" (account api server). 
Alternative constructors has been created to override those parameters, as NewClientWithUrl
- I was thinking in adding some validation, using json annotations, but taking in mind the time constraints, I preferred to invest it on testing deeply unmarshall protocol, as is one of the typical bug points, that consumes time but ensures results
- Simplicity has been key in the whole development
    - from api provided responses, and to spend the less possible time, I used an auto-generator that converts from json to struct (https://mholt.github.io/json-to-go/), and then clean out all the relative entities
    - library presented in 2 layers:
        -http helper that wraps original golang http client, that forwards context, which enables context cancellation or scoped timeouts, I preferred this approach over adding configured timeout on http.Client, basically to achieve the same timeout scoping we need to set it in multiple places (connection, read, write timeout...), while using context timeout is a unique solution that covers all the scenarios. Once said that, helper layer creates http requests, including encode/decode, executes http request and validates basic http response status codes translating them to errors. Http responses included too, allowing fine-grained validations (an example of that is the assertion of 201 status code on account created)  
        -api accessors are using the http handler layer, so they build the entry point of the library, in the current challenge just applied on the account scenario, but valid to other api endpoints.
        
- TDD philosophy has been followed to design the library, everything covered using unit-test, and integration tests can be found in test folder, those integration tests can serve as implementation examples too
- integration test fixtures implemented in list test, ideally on CI environment those fixtures would be created/destroyed from sql

## Unit test         
 ```
    go test -v --race ./...
```

## Integration tests
```
    docker-compose up
```

## Improvements
- validation layer can be easily achieved using json annotations (gopkg.in/go-playground/validator.v9)        

    
