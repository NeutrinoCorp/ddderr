package ddderr

import "net/http"

// HttpError is an RFC-compliant HTTP protocol problem object.
//
// For more information about the fields, please go to: https://datatracker.ietf.org/doc/html/rfc7807
type HttpError struct {
	Type     string `json:"type,omitempty"`
	Title    string `json:"title,omitempty"`
	Status   int    `json:"status,omitempty"`
	Detail   string `json:"detail,omitempty"`
	Instance string `json:"instance,omitempty"`
}

// NewHttpError builds an HttpError from the given DDD error
func NewHttpError(errType, instance string, err error) HttpError {
	code := http.StatusInternalServerError
	if err == nil {
		return HttpError{
			Type:     getHttpErrorType(errType, code),
			Title:    "",
			Status:   code,
			Detail:   "",
			Instance: instance,
		}
	}

	customErr, ok := err.(Error)
	if !ok {
		return HttpError{
			Type:     getHttpErrorType(errType, code),
			Title:    err.Error(),
			Status:   code,
			Detail:   err.Error(),
			Instance: instance,
		}
	}

	code = GetHttpStatusCode(customErr)
	return HttpError{
		Type:     getHttpErrorType(errType, code),
		Title:    customErr.Title(),
		Status:   code,
		Detail:   customErr.Description(),
		Instance: instance,
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
