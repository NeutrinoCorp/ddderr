package ddderr

import (
	"errors"
	"testing"
)

func TestGetDescription(t *testing.T) {
	// A. error description returns invalid description
	err := NewRequired("foo")
	format := "foo is required"
	if GetDescription(err) != format {
		t.Errorf("error does not return desired description format, expected: %s", format)
	}
}

func TestGetParentDescription(t *testing.T) {
	// A. error description returns invalid description
	err := NewNotFound(errors.New("row was not found"), "foo")
	format := "foo not found"
	if GetDescription(err) != format {
		t.Errorf("error does not return desired description format, expected: %s", format)
	}

	// B. error description returns invalid parent description
	format = "row was not found"
	if GetParentDescription(err) != format {
		t.Errorf("error does not return desired parent description format, expected: %s", format)
	}

	// C. error description returns non-empty string when sending invalid error
	err = NewRequired("foo")
	format = ""
	if GetParentDescription(err) != format {
		t.Error("error does not return desired parent description format, expected: empty string")
	}
}
