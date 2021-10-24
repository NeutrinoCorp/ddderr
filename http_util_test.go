package ddderr

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var getCodeHttpTestSuite = []struct {
	InErr   Error
	ExpCode int
}{
	{
		InErr:   Error{},
		ExpCode: http.StatusInternalServerError,
	},
	{
		InErr:   NewInfrastructure("generic title", "specific description"),
		ExpCode: http.StatusInternalServerError,
	},
	{
		InErr:   NewRemoteCall("tcp:172.16.52.1"),
		ExpCode: http.StatusBadGateway,
	},
	{
		InErr:   NewDomain("generic title", "specific description"),
		ExpCode: http.StatusBadRequest,
	},
	{
		InErr:   NewInvalidFormat("foo", "1", "5", "9"),
		ExpCode: http.StatusBadRequest,
	},
	{
		InErr:   NewOutOfRange("foo", 8, 256),
		ExpCode: http.StatusBadRequest,
	},
	{
		InErr:   NewRequired("foo"),
		ExpCode: http.StatusBadRequest,
	},
	{
		InErr:   NewAlreadyExists("foo"),
		ExpCode: http.StatusConflict,
	},
	{
		InErr:   NewNotFound("foo"),
		ExpCode: http.StatusNotFound,
	},
}

func TestGetHttpStatusCode(t *testing.T) {
	for _, tt := range getCodeHttpTestSuite {
		t.Run("", func(t *testing.T) {
			code := GetHttpStatusCode(tt.InErr)
			assert.Equal(t, tt.ExpCode, code)
		})
	}
}

var newHttpErrorTestSuite = []struct {
	InErrType  string
	InInstance string
	InErr      error
	ExpHttpErr HttpError
}{
	{
		InErrType:  "",
		InInstance: "",
		InErr:      nil,
		ExpHttpErr: HttpError{
			Type:       "",
			Title:      "",
			Status:     "",
			StatusCode: 0,
			Detail:     "",
			Instance:   "",
		},
	},
	{
		InErrType:  "https://neutrinocorp.org/iam/probs/generic-error",
		InInstance: "/users/12345/msg/abc",
		InErr:      nil,
		ExpHttpErr: HttpError{
			Type:       "",
			Title:      "",
			Status:     "",
			StatusCode: 0,
			Detail:     "",
			Instance:   "",
		},
	},
	{
		InErrType:  "",
		InInstance: "/users/12345/msg/abc",
		InErr:      errors.New("generic error"),
		ExpHttpErr: HttpError{
			Type:       "Internal Server Error",
			Title:      "generic error",
			Status:     "Internal Server Error",
			StatusCode: http.StatusInternalServerError,
			Detail:     "generic error",
			Instance:   "/users/12345/msg/abc",
		},
	},
	{
		InErrType:  "https://neutrinocorp.org/iam/probs/generic-error",
		InInstance: "/users/12345/msg/abc",
		InErr:      errors.New("generic error"),
		ExpHttpErr: HttpError{
			Type:       "https://neutrinocorp.org/iam/probs/generic-error",
			Title:      "generic error",
			Status:     "https://neutrinocorp.org/iam/probs/generic-error",
			StatusCode: http.StatusInternalServerError,
			Detail:     "generic error",
			Instance:   "/users/12345/msg/abc",
		},
	},
	{
		InErrType:  "https://neutrinocorp.org/iam/probs/not-found",
		InInstance: "/users/12345/msg/abc",
		InErr:      NewNotFound("foo").SetStatus("FooNotFound"),
		ExpHttpErr: HttpError{
			Type:       "https://neutrinocorp.org/iam/probs/not-found",
			Title:      "Resource not found",
			Status:     "FooNotFound",
			StatusCode: http.StatusNotFound,
			Detail:     "The resource foo was not found",
			Instance:   "/users/12345/msg/abc",
		},
	},
	{
		InErrType:  "https://neutrinocorp.org/iam/probs/required",
		InInstance: "/users/12345/msg/abc",
		InErr:      NewRequired("foo"),
		ExpHttpErr: HttpError{
			Type:       "https://neutrinocorp.org/iam/probs/required",
			Title:      "Missing property",
			Status:     "FooIsRequired",
			StatusCode: http.StatusBadRequest,
			Detail:     "The property foo is required",
			Instance:   "/users/12345/msg/abc",
		},
	},
	{
		InErrType:  "https://neutrinocorp.org/iam/probs/generic-domain",
		InInstance: "",
		InErr:      NewDomain("generic title", "specific description"),
		ExpHttpErr: HttpError{
			Type:       "https://neutrinocorp.org/iam/probs/generic-domain",
			Title:      "generic title",
			Status:     "Bad Request",
			StatusCode: http.StatusBadRequest,
			Detail:     "specific description",
			Instance:   "",
		},
	},
	{
		InErrType:  "https://neutrinocorp.org/iam/probs/generic-domain",
		InInstance: "",
		InErr:      NewDomain("generic title", "specific description").SetStatus("GenericError"),
		ExpHttpErr: HttpError{
			Type:       "https://neutrinocorp.org/iam/probs/generic-domain",
			Title:      "generic title",
			Status:     "GenericError",
			StatusCode: http.StatusBadRequest,
			Detail:     "specific description",
			Instance:   "",
		},
	},
	{
		InErrType:  "https://neutrinocorp.org/iam/probs/generic-infra",
		InInstance: "",
		InErr:      NewInfrastructure("generic title", "specific description"),
		ExpHttpErr: HttpError{
			Type:       "https://neutrinocorp.org/iam/probs/generic-infra",
			Title:      "generic title",
			Status:     "Internal Server Error",
			StatusCode: http.StatusInternalServerError,
			Detail:     "specific description",
			Instance:   "",
		},
	},
}

func TestNewHttpError(t *testing.T) {
	for _, tt := range newHttpErrorTestSuite {
		t.Run("", func(t *testing.T) {
			err := NewHttpError(tt.InErrType, tt.InInstance, tt.InErr)
			assert.EqualValues(t, tt.ExpHttpErr, err)
		})
	}
}
