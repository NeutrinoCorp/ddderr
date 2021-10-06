package main

import (
	"errors"
	"log"

	"github.com/neutrinocorp/ddderr"
)

func main() {
	_, err := getFooByID("123")
	if err != nil {
		log.Print(err)

		if customErr, ok := err.(ddderr.Error); ok {
			// Will output -> lib/pq mocked error
			log.Print(customErr.Parent())
			// Will output -> true
			log.Print(customErr.IsRemoteCall())
			// Will output -> true
			log.Print(customErr.IsInfrastructure())
		}
	}
}

func getFooByID(_ string) (interface{}, error) {
	return nil, ddderr.NewRemoteCall("localhost:5432").
		AttachParent(errors.New("pq: Failed to connect to PostgreSQL host"))
}
