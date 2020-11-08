# DDD Error

![Go Build](https://github.com/neutrinocorp/ddderr/workflows/Go/badge.svg?branch=master)
[![Go Documentation](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)][godocs]
[![Go Version](https://img.shields.io/github/go-mod/go-version/neutrinocorp/ddderr?style=flat-square)][goversion]

[actions]: https://github.com/neutrinocorp/ddderr/workflows/Go/badge.svg?branch=master
[godocs]: https://godoc.org/github.com/neutrinocorp/ddderr
[goversion]: https://img.shields.io/github/go-mod/go-version/neutrinocorp/ddderr

`DDD Error` is a _generic domain-driven exception wrapper_ made for Go.

Using existing validators such as playground's implementation is overwhelming because tag validation and the need to rewrite descriptions. With DDD Error, you may still use 3rd-party validators or make your own validations in your value objects, entities or aggregates.

In addition, infrastructure exceptions were added so you may be able to _catch specific kind of infrastructure errors._

Exceptions _descriptions are based on the Google Cloud API Design Guideline_.

`DDD Error` is compatible with popular error-handling packages such as [Hashicorp's go-multierror](https://github.com/hashicorp/go-multierror)

In conclusion, `DDD Error` aims to _ease the lack of exception handling_ in The Go Programming Language by defining a _wide selection of common exceptions_ 
which happen inside the _domain and/or infrastructure_ layer(s).

_Note: **DDD Error** is dependency-free, it uses built-in packages such as errors package._

## Installation
Install `DDD Error` by running the command

    go get github.com/neutrinocorp/ddderr
    
Full documentation is available
[here](https://pkg.go.dev/github.com/neutrinocorp/ddderr)


## Common Use Cases
- Implement retry strategy and circuit breaker resiliency patterns by adding Network exception to the whitelist.
- Not Acknowledging messages from an event bus if got a Network or Infrastructure generic exception.
- Get an HTTP/gRPC/OpenCensus status code from an error.
- Implement multiple strategies when an specific (or generic) type of error was thrown in.
- Fine-grained exception logging on infrastructure layer by using GetParentDescription() function.

## Usage

**HTTP status codes**

Set an HTTP error code depending on the exception.

```go
err := ddderr.NewNotFound(errors.New("row not found"), "user")
log.Print(ddderr.GetDescription(err)) // prints: "user not found"
log.Print(ddderr.GetParentDescription(err)) // prints: "row not found"

if errors.Is(err, ddderr.NotFound) {
  log.Print(http.StatusNotFound) // prints: 404
  return
}

log.Print(http.StatusInternalServerError) // prints: 500
```

**Domain generic exceptions**

Create a generic domain exception when other domain errors don't fulfill your requirements.

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

Create a generic infrastructure exception when other infrastructure exceptions don't fulfill your requirements.

```go
dbErr := errors.New("specific db error")
err := ddderr.NewInfrastructure(dbErr, "custom error description")
log.Print(ddderr.GetDescription(err)) // prints: "custom error description"
log.Print(ddderr.GetParentDescription(err)) // prints: "specific db error"
```

**Implement multiple strategies depending on exception kind**

Take an specific action depending on the exception kind.

```go
dbErr := errors.New("connection to host 127.0.0.1 failed")
err := ddderr.NewFailedRemoteCall(dbErr, "apache cassandra")
log.Print(ddderr.GetDescription(err)) // prints: "remote call to apache cassandra has failed"
log.Print(ddderr.GetParentDescription(err)) // prints: "connection to host 127.0.0.1 failed"

if ddderr.IsDomain(err) {
  // handle domain errors
  log.Print("error comes from domain")
} else if ddderr.IsInfrastructure(err) {
  // handle infrastructure errors
  log.Print("error comes from infrastructure")
}
```

## Requirements
- Go version >= 1.13
