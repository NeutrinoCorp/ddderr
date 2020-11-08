package ddderr

import (
	"strings"
)

// GetDescription returns an exception detailed description
func GetDescription(err error) string {
	sliced := strings.Split(err.Error(), "#")
	if len(sliced) == 1 {
		return sliced[0]
	}

	return sliced[1]
}

// GetParentDescription returns the parent's error description
//	Note: This function works only when using infrastructure exceptions
func GetParentDescription(err error) string {
	sliced := strings.Split(err.Error(), "#")
	if len(sliced) < 3 {
		return ""
	}

	return strings.Join(sliced[2:], "")
}
