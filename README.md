# :japanese_goblin: DDD Error

![Go Build](https://github.com/neutrinocorp/ddderr/workflows/Go/badge.svg?branch=master)
[![GoDoc](https://pkg.go.dev/badge/github.com/neutrinocorp/ddderr/v3)][godocs]
[![Go Report Card](https://goreportcard.com/badge/github.com/neutrinocorp/ddderr)](https://goreportcard.com/report/github.com/neutrinocorp/ddderr)
[![codebeat badge](https://codebeat.co/badges/22d865b6-c99a-469a-bb85-6b2d6f44a6fe)](https://codebeat.co/projects/github-com-neutrinocorp-ddderr-master)
[![Coverage Status][cov-img]][cov]
[![Go Version][go-img]][go]

`DDD Error` is a _reflection-free_ Domain-Driven error wrapper made for _Go_.

Using existing _validators_ such as playground's implementation is _overwhelming because tag validation_ and the need to rewrite descriptions. 
With `DDD Error`, _you may still use 3rd-party validators_ or _make your own validations_ in your value objects, entities or aggregates.

In addition, infrastructure exceptions were added so you may be able to _catch specific kind of infrastructure errors._

Exceptions _descriptions are based on the [Google Cloud API Design Guidelines](https://cloud.google.com/apis/design/errors)_.

`DDD Error` is compatible with popular error-handling packages such as [Hashicorp's go-multierror](https://github.com/hashicorp/go-multierror)

In conclusion, `DDD Error` aims to _ease the lack of exception handling_ in The Go Programming Language by defining a _wide selection of common exceptions_ 
which happen inside the _domain and/or infrastructure_ layer(s).

_Note: `DDD Error` is dependency-free, it complies with Go's built-in error interface and avoids reflection to increase overall performance._

## Installation
Install `DDD Error` by running the command

    go get github.com/neutrinocorp/ddderr/v3
    
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
err := ddderr.NewNotFound("foo")
log.Print(err) // prints: "The resource foo was not found"

if err.IsNotFound() {
  log.Print(http.StatusNotFound) // prints: 404
  return
}

log.Print(http.StatusInternalServerError) // prints: 500
```

Or use the _builtin HTTP utils:_

```go
err := ddderr.NewNotFound("foo")
log.Print(err) // prints: "The resource foo was not found"

// HttpError struct is ready to be marshaled using JSON encoding libs
//
// Function accepts the following optional params (specified on the RFC spec):
// - Type
// - Instance
httpErr := ddderr.NewHttpError("", "", err)
// Will output -> If errType param is empty, then HTTP status text is used as type
// (e.g. Not Found, Internal Server Error)
log.Print(httpErr.Type)
// Will output -> 404 as we got a NotFound error type
log.Print(httpErr.Status)
// Will output -> The resource foo was not found
log.Print(httpErr.Detail)
```

**Domain generic exceptions**

Create a generic domain exception when other domain errors don't fulfill your requirements.

```go
err := ddderr.NewDomain("generic error title", "foo has returned a generic domain error")
log.Print(err) // prints: "foo has returned a generic domain error"

if err.IsDomain() {
  log.Print(http.StatusBadRequest) // prints: 400
  return
}

log.Print(http.StatusInternalServerError) // prints: 500
```

**Infrastructure generic exceptions**

Create a generic infrastructure exception when other infrastructure exceptions don't fulfill your requirements.

```go
msgErr := errors.New("sarama: Apache kafka consumer error")
err := ddderr.NewInfrastructure("generic error title", "error while consuming message from queue").
	AttachParent(msgErr)
log.Print(err) // prints: "error while consuming message from queue"
log.Print(err.Parent()) // prints: "sarama: Apache kafka consumer error"
```

**Implement multiple strategies depending on exception kind**

Take an specific action depending on the exception kind.

```go
esErr := errors.New("failed to connect to Elasticsearch host http://127.0.0.1:9300")

err := ddderr.NewRemoteCall("http://127.0.0.1:9300").
	AttachParent(esErr)
log.Print("infrastructure error: ", err)                       // prints "failed to call external resource [http://127.0.0.1:9300]"
log.Print("infrastructure error resource: ", err.Property())       // http://127.0.0.1:9300
log.Print("is domain error: ", err.IsDomain())                 // false
log.Print("is infrastructure error: ", err.IsInfrastructure()) // true
log.Print("infrastructure error parent: ", err.Parent())       // prints "failed to connect to Elasticsearch host http://127.0.0.1:9300"

if err.IsRemoteCall() {
    // implement retry and/or circuit breaker pattern(s)
}
```

See [examples][examples] for more details.

## Requirements
- Go version >= 1.13

[actions]: https://github.com/neutrinocorp/ddderr/workflows/Go/badge.svg?branch=master
[godocs]: https://pkg.go.dev/github.com/neutrinocorp/ddderr/v3
[cov-img]: https://codecov.io/gh/NeutrinoCorp/ddderr/branch/master/graph/badge.svg
[cov]: https://codecov.io/gh/NeutrinoCorp/ddderr
[go-img]: https://img.shields.io/github/go-mod/go-version/NeutrinoCorp/ddderr?style=square
[go]: https://github.com/NeutrinoCorp/ddderr/blob/master/go.mod
[examples]: https://github.com/neutrinocorp/ddderr/tree/master/examples
