package main

import (
	"log"

	"github.com/neutrinocorp/ddderr"
)

func main() {
	execAlreadyExistsExample()
	execInvalidFormatExample()
	execNotFoundExample()
	execOutOfRangeExample()
	execRequiredExample()
	execSimpleExample()
}

func execAlreadyExistsExample() {
	err := ddderr.NewAlreadyExists("foo")
	log.Print("domain error: ", err)                               // prints "foo already exists"
	log.Print("domain error entity: ", err.Entity())               // foo
	log.Print("is domain error: ", err.IsDomain())                 // true
	log.Print("is infrastructure error: ", err.IsInfrastructure()) // false
}

func execInvalidFormatExample() {
	err := ddderr.NewInvalidFormat("foo", "image", "video", "document") // using variadic inputs, up to n formats
	log.Print("domain error: ", err)                                    // prints "foo contains an invalid format, expected [image,video,document]"
}

func execNotFoundExample() {
	err := ddderr.NewNotFound("foo")
	log.Print("domain error: ", err) // prints "foo not found"
}

func execOutOfRangeExample() {
	err := ddderr.NewOutOfRange("foo", 1, 256)
	log.Print("domain error: ", err) // prints "foo is out of range [1,256)"
}

func execRequiredExample() {
	err := ddderr.NewRequired("foo")
	log.Print("domain error: ", err) // prints "foo is required"
}

func execSimpleExample() {
	err := doSomethingAndFail() // hides ddderr specific error using Go's built-in error
	log.Print(err)              // prints "foo already exists"

	// use ddderr properties by casting the error
	domainErr, ok := err.(ddderr.Error)
	if !ok {
		panic("cannot cast ddderr")
	}
	log.Print(domainErr.Entity()) // prints foo
}

func doSomethingAndFail() error {
	return ddderr.NewAlreadyExists("foo")
}
