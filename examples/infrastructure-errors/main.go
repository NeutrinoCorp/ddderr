package main

import (
	"errors"
	"log"

	"github.com/neutrinocorp/ddderr"
)

func main() {
	execRemoteCallExample()
}

func execRemoteCallExample() {
	esErr := errors.New("failed to connect to Elasticsearch host http://127.0.0.1:9300")

	err := ddderr.NewRemoteCall(esErr, "elasticsearch")
	log.Print("infrastructure error: ", err)                       // prints "failed to call external resource [elasticsearch]"
	log.Print("infrastructure error entity: ", err.Entity())       // empty string
	log.Print("is domain error: ", err.IsDomain())                 // false
	log.Print("is infrastructure error: ", err.IsInfrastructure()) // true
	log.Print("infrastructure error parent: ", err.Parent())       // prints "failed to connect to Elasticsearch host http://127.0.0.1:9300"
}
