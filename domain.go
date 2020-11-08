package ddderr

import (
	"errors"
	"fmt"
)

var (
	// GenericDomain An error occurred inside a domain
	GenericDomain = errors.New("domain exception")
	// Required A field/resource is missing from the call/request
	Required    = errors.New("missing required field")
	requiredFmt = "%s is required"
	// InvalidFormat A field/resource has an invalid format/type
	InvalidFormat    = errors.New("field with invalid format")
	invalidFormatFmt = "%s has an invalid format, expected (%s)"
	// OutOfRange A field/resource is out of the defined range
	OutOfRange    = errors.New("field out of range")
	outOfRangeFmt = "%s is out of range [%s, %s)"
)

// NewDomain constructs a generic domain exception
func NewDomain(format string) error {
	return fmt.Errorf(errFormat, GenericDomain, format)
}

// NewRequired constructs a Required exception
func NewRequired(field string) error {
	return fmt.Errorf(errFormat, Required, fmt.Sprintf(requiredFmt, field))
}

// NewInvalidFormat constructs an InvalidFormat exception
func NewInvalidFormat(field, format string) error {
	return fmt.Errorf(errFormat, InvalidFormat, fmt.Sprintf(invalidFormatFmt, field, format))
}

// NewOutOfRange constructs an OutOfRange exception
func NewOutOfRange(field, x, y string) error {
	return fmt.Errorf(errFormat, OutOfRange, fmt.Sprintf(outOfRangeFmt, field, x, y))
}
