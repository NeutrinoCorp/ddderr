package ddderr

import (
	"errors"
	"testing"
)

func TestNewInfrastructure(t *testing.T) {
	customErr := errors.New("infrastructure error")
	err := NewInfrastructure(customErr, "custom exception")

	// A. error does not contain exception kind
	if !errors.Is(err, GenericInfrastructure) {
		t.Error("error does not contain GenericInfrastructure exception kind")
	}

	// B. error does not returns desired string
	format := "infrastructure exception#custom exception#infrastructure error"
	if err.Error() != format {
		t.Errorf("error does not return GenericInfrastructure desired format, expected: %s", format)
	}
}

func TestNewAlreadyExists(t *testing.T) {
	err := NewAlreadyExists(errors.New("row replacing failed, already exists"), "foo")

	// A. error does not contain exception kind
	if !errors.Is(err, AlreadyExists) {
		t.Error("error does not contain AlreadyExists exception kind")
	}

	// B. error does not returns desired string
	format := "already exists#foo already exists#row replacing failed, already exists"
	if err.Error() != format {
		t.Errorf("error does not return AlreadyExists desired format, expected: %s", format)
	}
}

func TestNewNotFound(t *testing.T) {
	err := NewNotFound(errors.New("row was not found"), "foo")

	// A. error does not contain exception kind
	if !errors.Is(err, NotFound) {
		t.Error("error does not contain NotFound exception kind")
	}

	// B. error does not returns desired string
	format := "not found#foo not found#row was not found"
	if err.Error() != format {
		t.Errorf("error does not return NotFound desired format, expected: %s", format)
	}
}

func TestNewFailedRemoteCall(t *testing.T) {
	err := NewFailedRemoteCall(errors.New("connection to host 127.0.0.1 failed"), "foo service")

	// A. error does not contain exception kind
	if !errors.Is(err, FailedRemoteCall) {
		t.Error("error does not contain FailedRemoteCall exception kind")
	}

	// B. error does not returns desired string
	format := "remote call failed#remote call to foo service has failed#connection to host 127.0.0.1 failed"
	if err.Error() != format {
		t.Errorf("error does not return FailedRemoteCall desired format, expected: %s", format)
	}
}
