package dddeerr

import (
	"errors"
	"testing"
)

func TestNewDomain(t *testing.T) {
	err := NewDomain("custom domain exception")

	// A. error does not contain exception kind
	if !errors.Is(err, GenericDomain) {
		t.Error("error does not contain GenericDomain exception kind")
	}

	// B. error does not returns desired string
	format := "domain exception#custom domain exception"
	if err.Error() != format {
		t.Errorf("error does not return GenericDomain desired format, expected: %s", format)
	}
}

func TestNewRequiredField(t *testing.T) {
	err := NewRequired("foo")

	// A. error does not contain exception kind
	if !errors.Is(err, Required) {
		t.Error("error does not contain RequiredField exception kind")
	}

	// B. error does not returns desired string
	format := "missing required field#foo is required"
	if err.Error() != format {
		t.Errorf("error does not return RequiredField desired format, expected: %s", format)
	}
}

func TestNewInvalidFormat(t *testing.T) {
	err := NewInvalidFormat("foo", "int64")

	// A. error does not contain exception kind
	if !errors.Is(err, InvalidFormat) {
		t.Error("error does not contain InvalidFormat exception kind")
	}

	// B. error does not returns desired string
	format := "field with invalid format#foo has an invalid format, expected (int64)"
	if err.Error() != format {
		t.Errorf("error does not return InvalidFormat desired format, expected: %s", format)
	}
}

func TestNewOutOfRange(t *testing.T) {
	err := NewOutOfRange("foo", "0", "256")

	// A. error does not contain exception kind
	if !errors.Is(err, OutOfRange) {
		t.Error("error does not contain OutOfRange exception kind")
	}

	// B. error does not returns desired string
	format := "field out of range#foo is out of range [0, 256)"
	if err.Error() != format {
		t.Errorf("error does not return OutOfRange desired format, expected: %s", format)
	}
}
