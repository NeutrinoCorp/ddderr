package main

import (
	"log"

	"github.com/neutrinocorp/ddderr"
)

func main() {
	_, err := getFooByID("123")
	if err != nil {
		// Will output -> The resource foo was not found
		log.Print(err)
	}
}

func getFooByID(_ string) (interface{}, error) {
	return nil, ddderr.NewNotFound("foo")
}
