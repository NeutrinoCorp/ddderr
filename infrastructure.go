package ddderr

import (
	"errors"
	"fmt"
)

var (
	// GenericInfrastructure An error occurred inside an infrastructure component
	GenericInfrastructure = errors.New("infrastructure exception")
	// AlreadyExists The given resource already exists
	AlreadyExists    = errors.New("already exists")
	alreadyExistsFmt = "%s already exists"
	// NotFound The requested resource was not found
	NotFound    = errors.New("not found")
	notFoundFmt = "%s not found"
	// FailedRemoteCall A remote call failed
	//	Note: Might be used for retry & circuit breaker resiliency patterns implementation(s)
	FailedRemoteCall    = errors.New("remote call failed")
	failedRemoteCallFmt = "remote call to %s has failed"
)

// NewInfrastructure constructs a generic domain exception
func NewInfrastructure(err error, format string) error {
	return fmt.Errorf(infraErrFormat, GenericInfrastructure, format, err.Error())
}

// NewAlreadyExists constructs an AlreadyExists exception
func NewAlreadyExists(err error, resource string) error {
	return fmt.Errorf(infraErrFormat, AlreadyExists, fmt.Sprintf(alreadyExistsFmt, resource), err)
}

// NewNotFound constructs a NotFound exception
func NewNotFound(err error, resource string) error {
	return fmt.Errorf(infraErrFormat, NotFound, fmt.Sprintf(notFoundFmt, resource), err)
}

// NewFailedRemoteCall constructs a FailedRemoteCall exception
func NewFailedRemoteCall(err error, resource string) error {
	return fmt.Errorf(infraErrFormat, FailedRemoteCall, fmt.Sprintf(failedRemoteCallFmt, resource), err)
}
