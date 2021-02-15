package main

import (
	"errors"
	"log"

	"github.com/neutrinocorp/ddderr"
)

func main() {
	execDomainExample()
	execInfraExample()
}

// generic domain error
func execDomainExample() {
	err := ddderr.NewDomain("foo", "foo has returned a generic domain error")
	log.Print("domain error: ", err)                               // prints error string
	log.Print("domain error entity: ", err.Entity())               // foo
	log.Print("is domain error: ", err.IsDomain())                 // true
	log.Print("is infrastructure error: ", err.IsInfrastructure()) // false
}

// generic infra error
func execInfraExample() {
	err := ddderr.NewInfrastructure(errors.New("apache kafka consumer error"), "error while consuming message from queue")
	log.Print("infrastructure error: ", err)                       // prints error string
	log.Print("infrastructure error entity: ", err.Entity())       // empty string
	log.Print("is domain error: ", err.IsDomain())                 // false
	log.Print("is infrastructure error: ", err.IsInfrastructure()) // true
	log.Print("infrastructure error parent: ", err.Parent())       // prints parent error string
}
