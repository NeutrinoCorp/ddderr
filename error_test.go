package ddderr

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError_Setters(t *testing.T) {
	pqMockedErr := errors.New("pq: Generic PostgreSQL driver error")
	err := NewDomain("", "").
		AttachParent(pqMockedErr).
		SetKind("CustomKind").
		SetTitle("generic title").
		SetDescription("specific description").
		SetProperty("foo")

	assert.EqualValues(t, pqMockedErr, err.Parent())
	assert.Equal(t, "CustomKind", err.Kind())
	assert.Equal(t, "generic title", err.Title())
	assert.Equal(t, "specific description", err.Description())
	assert.Equal(t, "foo", err.Property())
	assert.Equal(t, "specific description", err.Error())
}

var newDomainTestSuite = []struct {
	InTitle string
	InDesc  string
	Exp     Error
}{
	{
		InTitle: "",
		InDesc:  "",
		Exp: Error{
			parent:      nil,
			group:       domain,
			kind:        unknownDomain,
			property:    "",
			title:       "",
			description: "",
		},
	},
	{
		InTitle: "generic title",
		InDesc:  "foo description",
		Exp: Error{
			parent:      nil,
			group:       domain,
			kind:        unknownDomain,
			property:    "",
			title:       "generic title",
			description: "foo description",
		},
	},
}

func TestNewDomain(t *testing.T) {
	for _, tt := range newDomainTestSuite {
		t.Run("", func(t *testing.T) {
			err := NewDomain(tt.InTitle, tt.InDesc)
			assert.EqualValues(t, tt.Exp, err)
			assert.Equal(t, tt.InTitle, err.Title())
			assert.Equal(t, tt.InDesc, err.Description())
			assert.Equal(t, tt.Exp.Kind(), err.Kind())
			assert.Empty(t, err.Property())
			assert.True(t, err.IsDomain())
			assert.False(t, err.IsInfrastructure())
			assert.False(t, err.IsRequired())
			assert.False(t, err.IsOutOfRange())
			assert.False(t, err.IsInvalidFormat())
			assert.False(t, err.IsAlreadyExists())
			assert.False(t, err.IsRemoteCall())
			assert.False(t, err.IsNotFound())
		})
	}
}

var newInfraTestSuite = []struct {
	InTitle string
	InDesc  string
	Exp     Error
}{
	{
		InTitle: "",
		InDesc:  "",
		Exp: Error{
			parent:      nil,
			group:       infrastructure,
			kind:        unknownInfrastructure,
			property:    "",
			title:       "",
			description: "",
		},
	},
	{
		InTitle: "generic title",
		InDesc:  "foo description",
		Exp: Error{
			parent:      nil,
			group:       infrastructure,
			kind:        unknownInfrastructure,
			property:    "",
			title:       "generic title",
			description: "foo description",
		},
	},
}

func TestNewInfrastructure(t *testing.T) {
	for _, tt := range newInfraTestSuite {
		t.Run("", func(t *testing.T) {
			err := NewInfrastructure(tt.InTitle, tt.InDesc)
			assert.EqualValues(t, tt.Exp, err)
			assert.Equal(t, tt.InTitle, err.Title())
			assert.Equal(t, tt.InDesc, err.Description())
			assert.Equal(t, tt.Exp.Kind(), err.Kind())
			assert.Empty(t, err.Property())
			assert.True(t, err.IsInfrastructure())
			assert.False(t, err.IsDomain())
			assert.False(t, err.IsRequired())
			assert.False(t, err.IsOutOfRange())
			assert.False(t, err.IsInvalidFormat())
			assert.False(t, err.IsAlreadyExists())
			assert.False(t, err.IsRemoteCall())
			assert.False(t, err.IsNotFound())
		})
	}
}

var newRemoteCallTestSuite = []struct {
	InExternalResource string
	Exp                Error
}{
	{
		InExternalResource: "",
		Exp: Error{
			parent:      nil,
			group:       infrastructure,
			kind:        remoteCall,
			property:    "",
			title:       "Remote call failed",
			description: "Failed to call external resource",
		},
	},
	{
		InExternalResource: "https://foo.com",
		Exp: Error{
			parent:      nil,
			group:       infrastructure,
			kind:        remoteCall,
			property:    "https://foo.com",
			title:       "Remote call failed",
			description: "Failed to call external resource [https://foo.com]",
		},
	},
}

func TestNewRemoteCall(t *testing.T) {
	for _, tt := range newRemoteCallTestSuite {
		t.Run("", func(t *testing.T) {
			err := NewRemoteCall(tt.InExternalResource)
			assert.EqualValues(t, tt.Exp, err)
			assert.Equal(t, tt.Exp.Title(), err.Title())
			assert.Equal(t, tt.Exp.Description(), err.Description())
			assert.Equal(t, tt.Exp.Property(), err.Property())
			assert.Equal(t, tt.Exp.Kind(), err.Kind())
			assert.True(t, err.IsRemoteCall())
			assert.True(t, err.IsInfrastructure())
			assert.False(t, err.IsDomain())
			assert.False(t, err.IsAlreadyExists())
			assert.False(t, err.IsRequired())
			assert.False(t, err.IsOutOfRange())
			assert.False(t, err.IsInvalidFormat())
			assert.False(t, err.IsNotFound())
		})
	}
}

var newNotFoundTestSuite = []struct {
	InResource string
	Exp        Error
}{
	{
		InResource: "",
		Exp: Error{
			parent:      nil,
			group:       domain,
			kind:        notFound,
			property:    "",
			title:       "Resource not found",
			description: "not found",
		},
	},
	{
		InResource: "foo",
		Exp: Error{
			parent:      nil,
			group:       domain,
			kind:        notFound,
			property:    "foo",
			title:       "Resource not found",
			description: "The resource foo was not found",
		},
	},
}

func TestNewNotFound(t *testing.T) {
	for _, tt := range newNotFoundTestSuite {
		t.Run("", func(t *testing.T) {
			err := NewNotFound(tt.InResource)
			assert.EqualValues(t, tt.Exp, err)
			assert.Equal(t, tt.Exp.Title(), err.Title())
			assert.Equal(t, tt.Exp.Description(), err.Description())
			assert.Equal(t, tt.Exp.Property(), err.Property())
			assert.Equal(t, tt.Exp.Kind(), err.Kind())
			assert.True(t, err.IsNotFound())
			assert.True(t, err.IsDomain())
			assert.False(t, err.IsInfrastructure())
			assert.False(t, err.IsRemoteCall())
			assert.False(t, err.IsAlreadyExists())
			assert.False(t, err.IsRequired())
			assert.False(t, err.IsOutOfRange())
			assert.False(t, err.IsInvalidFormat())
		})
	}
}

var newAlreadyExistsTestSuite = []struct {
	InResource string
	Exp        Error
}{
	{
		InResource: "",
		Exp: Error{
			parent:      nil,
			group:       domain,
			kind:        alreadyExists,
			property:    "",
			title:       "Resource already exists",
			description: "already exists",
		},
	},
	{
		InResource: "foo",
		Exp: Error{
			parent:      nil,
			group:       domain,
			kind:        alreadyExists,
			property:    "foo",
			title:       "Resource already exists",
			description: "The resource foo already exists",
		},
	},
}

func TestNewAlreadyExists(t *testing.T) {
	for _, tt := range newAlreadyExistsTestSuite {
		t.Run("", func(t *testing.T) {
			err := NewAlreadyExists(tt.InResource)
			assert.EqualValues(t, tt.Exp, err)
			assert.Equal(t, tt.Exp.Title(), err.Title())
			assert.Equal(t, tt.Exp.Description(), err.Description())
			assert.Equal(t, tt.Exp.Property(), err.Property())
			assert.Equal(t, tt.Exp.Kind(), err.Kind())
			assert.True(t, err.IsAlreadyExists())
			assert.True(t, err.IsDomain())
			assert.False(t, err.IsInfrastructure())
			assert.False(t, err.IsRequired())
			assert.False(t, err.IsOutOfRange())
			assert.False(t, err.IsInvalidFormat())
			assert.False(t, err.IsRemoteCall())
			assert.False(t, err.IsNotFound())
		})
	}
}

var newOutOfRangeTestSuite = []struct {
	InProp string
	InLimA int
	InLimB int
	Exp    Error
}{
	{
		InProp: "",
		InLimA: 0,
		InLimB: 0,
		Exp: Error{
			parent:      nil,
			group:       domain,
			kind:        outOfRange,
			property:    "",
			title:       "Property is out of the specified range",
			description: "out of range [0,0)",
		},
	},
	{
		InProp: "foo",
		InLimA: 0,
		InLimB: 0,
		Exp: Error{
			parent:      nil,
			group:       domain,
			kind:        outOfRange,
			property:    "foo",
			title:       "Property is out of the specified range",
			description: "The property foo is out of range [0,0)",
		},
	},
	{
		InProp: "",
		InLimA: 8,
		InLimB: 256,
		Exp: Error{
			parent:      nil,
			group:       domain,
			kind:        outOfRange,
			property:    "",
			title:       "Property is out of the specified range",
			description: "out of range [8,256)",
		},
	},
	{
		InProp: "foo",
		InLimA: 8,
		InLimB: 256,
		Exp: Error{
			parent:      nil,
			group:       domain,
			kind:        outOfRange,
			property:    "foo",
			title:       "Property is out of the specified range",
			description: "The property foo is out of range [8,256)",
		},
	},
}

func TestNewOutOfRange(t *testing.T) {
	for _, tt := range newOutOfRangeTestSuite {
		t.Run("", func(t *testing.T) {
			err := NewOutOfRange(tt.InProp, tt.InLimA, tt.InLimB)
			assert.EqualValues(t, tt.Exp, err)
			assert.Equal(t, tt.Exp.Title(), err.Title())
			assert.Equal(t, tt.Exp.Description(), err.Description())
			assert.Equal(t, tt.Exp.Property(), err.Property())
			assert.Equal(t, tt.Exp.Kind(), err.Kind())
			assert.True(t, err.IsOutOfRange())
			assert.True(t, err.IsDomain())
			assert.False(t, err.IsInfrastructure())
			assert.False(t, err.IsRequired())
			assert.False(t, err.IsNotFound())
			assert.False(t, err.IsRemoteCall())
			assert.False(t, err.IsAlreadyExists())
			assert.False(t, err.IsInvalidFormat())
		})
	}
}

var newInvalidFormatTestSuite = []struct {
	InProperty string
	InFormats  []string
	Exp        Error
}{
	{
		InProperty: "",
		InFormats:  nil,
		Exp: Error{
			parent:      nil,
			group:       domain,
			kind:        invalidFormat,
			property:    "",
			title:       "Property is not a valid format",
			description: "invalid format, expected []",
		},
	},
	{
		InProperty: "",
		InFormats:  []string{"foo"},
		Exp: Error{
			parent:      nil,
			group:       domain,
			kind:        invalidFormat,
			property:    "",
			title:       "Property is not a valid format",
			description: "invalid format, expected [foo]",
		},
	},
	{
		InProperty: "foo",
		InFormats:  []string{"bar"},
		Exp: Error{
			parent:      nil,
			group:       domain,
			kind:        invalidFormat,
			property:    "foo",
			title:       "Property is not a valid format",
			description: "The property foo has an invalid format, expected [bar]",
		},
	},
	{
		InProperty: "foo",
		InFormats:  []string{"bar", "baz"},
		Exp: Error{
			parent:      nil,
			group:       domain,
			kind:        invalidFormat,
			property:    "foo",
			title:       "Property is not a valid format",
			description: "The property foo has an invalid format, expected [bar,baz]",
		},
	},
}

func TestNewInvalidFormat(t *testing.T) {
	for _, tt := range newInvalidFormatTestSuite {
		t.Run("", func(t *testing.T) {
			err := NewInvalidFormat(tt.InProperty, tt.InFormats...)
			assert.EqualValues(t, tt.Exp, err)
			assert.Equal(t, tt.Exp.Title(), err.Title())
			assert.Equal(t, tt.Exp.Description(), err.Description())
			assert.Equal(t, tt.Exp.Property(), err.Property())
			assert.Equal(t, tt.Exp.Kind(), err.Kind())
			assert.True(t, err.IsInvalidFormat())
			assert.True(t, err.IsDomain())
			assert.False(t, err.IsInfrastructure())
			assert.False(t, err.IsAlreadyExists())
			assert.False(t, err.IsRequired())
			assert.False(t, err.IsOutOfRange())
			assert.False(t, err.IsRemoteCall())
			assert.False(t, err.IsNotFound())
		})
	}
}

var newRequiredTestSuite = []struct {
	InProp string
	Exp    Error
}{
	{
		InProp: "",
		Exp: Error{
			parent:      nil,
			group:       domain,
			kind:        required,
			property:    "",
			title:       "Missing property",
			description: "required",
		},
	},
	{
		InProp: "foo",
		Exp: Error{
			parent:      nil,
			group:       domain,
			kind:        required,
			property:    "foo",
			title:       "Missing property",
			description: "The property foo is required",
		},
	},
}

func TestNewRequired(t *testing.T) {
	for _, tt := range newRequiredTestSuite {
		t.Run("", func(t *testing.T) {
			err := NewRequired(tt.InProp)
			assert.EqualValues(t, tt.Exp, err)
			assert.Equal(t, tt.Exp.Title(), err.Title())
			assert.Equal(t, tt.Exp.Description(), err.Description())
			assert.Equal(t, tt.Exp.Property(), err.Property())
			assert.Equal(t, tt.Exp.Kind(), err.Kind())
			assert.True(t, err.IsRequired())
			assert.True(t, err.IsDomain())
			assert.False(t, err.IsInfrastructure())
			assert.False(t, err.IsNotFound())
			assert.False(t, err.IsRemoteCall())
			assert.False(t, err.IsAlreadyExists())
			assert.False(t, err.IsOutOfRange())
			assert.False(t, err.IsInvalidFormat())
		})
	}
}
