package ddderr

import "errors"

// IsDomain verifies if the given exception is a domain exception
func IsDomain(err error) bool {
	switch {
	case errors.Is(err, Required):
		return true
	case errors.Is(err, InvalidFormat):
		return true
	case errors.Is(err, OutOfRange):
		return true
	case errors.Is(err, GenericDomain):
		return true
	default:
		return false
	}
}

// IsInfrastructure verifies if the given exception is an Infrastructure exception
func IsInfrastructure(err error) bool {
	switch {
	case errors.Is(err, AlreadyExists):
		return true
	case errors.Is(err, NotFound):
		return true
	case errors.Is(err, FailedRemoteCall):
		return true
	case errors.Is(err, GenericInfrastructure):
		return true
	default:
		return false
	}
}
