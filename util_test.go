package ddderr

import (
	"errors"
	"testing"
)

func TestIsDomain(t *testing.T) {
	//	A. send valid domain exception
	//	A-1. Required
	err := NewRequired("foo")
	if !IsDomain(err) {
		t.Error("validation error, Required exception is a valid domain exception")
	}

	//	A-2. InvalidFormat
	err = NewInvalidFormat("foo", "int")
	if !IsDomain(err) {
		t.Error("validation error, InvalidFormat exception is a valid domain exception")
	}

	//	A-3. OutOfRange
	err = NewOutOfRange("foo", "0", "512")
	if !IsDomain(err) {
		t.Error("validation error, OutOfRange exception is a valid domain exception")
	}

	//	A-4. Generic
	err = NewDomain("generic exception")
	if !IsDomain(err) {
		t.Error("validation error, GenericDomain exception is a valid domain exception")
	}

	// B. send infrastructure errors
	//	B-1. NotFound
	err = NewNotFound(errors.New("row not found"), "foo")
	if IsDomain(err) {
		t.Error("validation error, NotFound exception is not a valid domain exception")
	}

	//	B-2. AlreadyExists
	err = NewAlreadyExists(errors.New("row already exists"), "foo")
	if IsDomain(err) {
		t.Error("validation error, AlreadyExists exception is not a valid domain exception")
	}

	//	B-3. FailedRemoteCall
	err = NewFailedRemoteCall(errors.New("call to 127.0.0.1 failed"), "apache cassandra")
	if IsDomain(err) {
		t.Error("validation error, FailedRemoteCall exception is not a valid domain exception")
	}

	//	B-4. Generic
	err = NewInfrastructure(errors.New("infrastructure error"), "generic exception")
	if IsDomain(err) {
		t.Error("validation error, GenericInfrastructure exception is not a valid domain exception")
	}
}

func TestIsInfrastructure(t *testing.T) {
	//	A. send valid infrastructure exception
	//	A-1. NotFound
	err := NewNotFound(errors.New("row not found"), "foo")
	if !IsInfrastructure(err) {
		t.Error("validation error, NotFound exception is a valid infrastructure exception")
	}

	//	A-2. AlreadyExists
	err = NewAlreadyExists(errors.New("row already exists"), "foo")
	if !IsInfrastructure(err) {
		t.Error("validation error, AlreadyExists exception is a valid infrastructure exception")
	}

	//	A-3. FailedRemoteCall
	err = NewFailedRemoteCall(errors.New("call to 127.0.0.1 failed"), "apache cassandra")
	if !IsInfrastructure(err) {
		t.Error("validation error, FailedRemoteCall exception is a valid infrastructure exception")
	}

	//	A-4. Generic
	err = NewInfrastructure(errors.New("infrastructure error"), "generic exception")
	if !IsInfrastructure(err) {
		t.Error("validation error, GenericInfrastructure exception is a valid infrastructure exception")
	}

	// B. send domain error
	//	B-1. Required
	err = NewRequired("foo")
	if IsInfrastructure(err) {
		t.Error("validation error, Required exception is not a valid infrastructure exception")
	}

	//	A-2. InvalidFormat
	err = NewInvalidFormat("foo", "int")
	if IsInfrastructure(err) {
		t.Error("validation error, InvalidFormat exception is not a valid infrastructure exception")
	}

	//	A-3. OutOfRange
	err = NewOutOfRange("foo", "0", "512")
	if IsInfrastructure(err) {
		t.Error("validation error, OutOfRange exception is not a valid infrastructure exception")
	}

	//	A-4. Generic
	err = NewDomain("generic exception")
	if IsInfrastructure(err) {
		t.Error("validation error, GenericDomain exception is not a valid infrastructure exception")
	}
}
