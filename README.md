# DDD Error
**_DDD Error_** is a _generic domain-driven exception wrapper_ made for Go.

**_DDD Error_** aims to _ease the lack of exception handling_ in The Go Programming Language by defining a _wide selection of common exceptions_ 
which happen inside the _domain and/or infrastructure_ layer(s).

**_DDD Error_** is dependency-free, it _uses built-in packages_ such as errors.


## Use cases
- Implement retry strategy and circuit breaker resiliency patterns by adding Network exception to the whitelist.
- Not Acknowledge messages from an event bus if got a Network or Infrastructure generic exception.
- Get an HTTP/gRPC/OpenCensus status code from an error.
- Implement multiple strategies when an specific (or generic) type of error was thrown in.
- More.

## Usability

**HTTP status codes**
```go
err := ddderr.NewNotFound("user")
log.Print(ddderr.GetDescription(err)) // prints: "user not found"

if errors.Is(err, ddderr.NotFound) {
  log.Print(http.StatusNotFound) // prints: 404
  return
}

log.Print(http.StatusInternalServerError) // prints: 500
```

**Domain generic exceptions**
```go
err := ddderr.NewDomain("something happened inside domain")
log.Print(ddderr.GetDescription(err)) // prints: "something happened inside domain"

if errors.Is(err, ddderr.GenericDomain) {
  log.Print(http.StatusBadRequest) // prints: 400
  return
}

log.Print(http.StatusInternalServerError) // prints: 500
```

**Infrastructure generic exceptions**
```go
dbErr := errors.New("specific db error")
err := ddderr.NewInfrastructure(dbErr, "custom error description")
log.Print(ddderr.GetDescription(err)) // prints: "custom error description"
log.Print(ddderr.GetParentDescription(err)) // prints: "specific db error"
```

**Implement multiple strategies depending on exception kind**
```go
dbErr := errors.New("failed to connect to host 127.0.0.1")
err := ddderr.NewFailedRemoteCall(dbErr, "apache cassandra")
log.Print(ddderr.GetDescription(err)) // prints: "remote call to apache cassandra has failed"
log.Print(ddderr.GetParentDescription(err)) // prints: "failed to connect to host 127.0.0.1"

if ddderr.IsDomain(err) {
  // handle domain errors
  log.Print("error comes from domain")
} else if ddderr.IsInfrastructure(err) {
  // handle infrastructure errors
  log.Print("error comes from infrastructure")
}
```
