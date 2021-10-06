package main

import (
	"log"

	"github.com/neutrinocorp/ddderr"
)

func main() {
	_, err := getFooByID("123")
	if err != nil {
		// HttpError struct is ready to be marshaled using JSON encoding libs
		httpErr := ddderr.NewHttpError("", "", err)
		// Will output -> If errType param is empty, then HTTP status text is used as type
		// (e.g. Not Found, Internal Server Error)
		log.Print(httpErr.Type)
		// Will output -> 404 as we got a NotFound error type
		log.Print(httpErr.Status)
		// Will output -> The resource foo was not found
		log.Print(httpErr.Detail)
	}
}

func getFooByID(_ string) (interface{}, error) {
	return nil, ddderr.NewNotFound("foo")
}
