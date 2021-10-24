package ddderr

import (
	"net/http"
)

// HttpError is an RFC-compliant HTTP protocol problem object.
//
// For more information about the fields, please go to: https://datatracker.ietf.org/doc/html/rfc7807
type HttpError struct {
	Type       string `json:"type,omitempty"`
	Title      string `json:"title,omitempty"`
	Status     string `json:"status,omitempty"`
	StatusCode int    `json:"status_code,omitempty"`
	Detail     string `json:"detail,omitempty"`
	Instance   string `json:"instance,omitempty"`
}

// NewHttpError builds an HttpError from the given DDD error
func NewHttpError(errType, instance string, err error) HttpError {
	if err == nil {
		return HttpError{}
	}

	code := http.StatusInternalServerError
	errHttpType := getHttpErrorType(errType, code)

	customErr, ok := err.(Error)
	if !ok {
		return HttpError{
			Type:       errHttpType,
			Title:      err.Error(),
			Status:     errHttpType,
			StatusCode: code,
			Detail:     err.Error(),
			Instance:   instance,
		}
	}

	code = GetHttpStatusCode(customErr)
	return HttpError{
		Type:       errHttpType,
		Title:      customErr.Title(),
		Status:     getHttpDddErrorType(customErr, code),
		StatusCode: code,
		Detail:     customErr.Description(),
		Instance:   instance,
	}
}

// retrieves a generic HTTP problem object type.
//
// For more information, go to: https://datatracker.ietf.org/doc/html/rfc7807#section-4.2
func getHttpErrorType(rootType string, status int) string {
	if rootType != "" {
		return rootType
	}
	return http.StatusText(status)
}

// retrieves a status name from an Error or a generic HTTP problem object type.
//
// For more information, go to: https://datatracker.ietf.org/doc/html/rfc7807#section-4.2
func getHttpDddErrorType(err Error, status int) string {
	if statusName := err.Status(); statusName != "" {
		return statusName
	}
	return http.StatusText(status)
}

// GetHttpStatusCode retrieves an HTTP status code from the given error
func GetHttpStatusCode(err Error) int {
	switch {
	case err.IsAlreadyExists():
		return http.StatusConflict
	case err.IsNotFound():
		return http.StatusNotFound
	case err.IsInvalidFormat() || err.IsRequired() || err.IsOutOfRange() || err.IsDomain():
		return http.StatusBadRequest
	case err.IsRemoteCall():
		return http.StatusBadGateway
	default:
		return http.StatusInternalServerError
	}
}
