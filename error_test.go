package ddderr

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError_Setters(t *testing.T) {
	pqMockedErr := errors.New("pq: Generic PostgreSQL driver error")
	err := NewDomain("", "").
		SetParent(pqMockedErr).
		SetKind("CustomKind").
		SetTitle("generic title").
		SetDescription("specific description").
		SetProperty("foo").
		SetStatus("FooIsBad")

	assert.EqualValues(t, pqMockedErr, err.Parent())
	assert.Equal(t, "CustomKind", err.Kind())
	assert.Equal(t, "generic title", err.Title())
	assert.Equal(t, "specific description", err.Description())
	assert.Equal(t, "foo", err.Property())
	assert.Equal(t, "specific description", err.Error())
	assert.Equal(t, "FooIsBad", err.Status())

	err = NewOutOfRange("fooBar", 8, 16)
	err = err.SetProperty("barBaz")
	assert.Equal(t, "The property barBaz is out of range [8,16)", err.Error())

}

var dynamicFieldTests = []struct {
	In             Error
	InDynamicField string
	ExpDesc        string
	ExpStatus      string
}{
	{
		In:             NewOutOfRange("bar", 0, 3),
		InDynamicField: "foo",
		ExpDesc:        "The property foo is out of range [0,3)",
		ExpStatus:      "FooOutOfRange",
	},
	{
		In:             NewAlreadyExists("bar"),
		InDynamicField: "foo",
		ExpDesc:        "The resource foo already exists",
		ExpStatus:      "FooAlreadyExists",
	},
	{
		In:             NewInvalidFormat("bar", "jpeg", "gif"),
		InDynamicField: "foo",
		ExpDesc:        "The property foo has an invalid format, expected [jpeg,gif]",
		ExpStatus:      "FooInvalidFormat",
	},
	{
		In: Error{
			kind: invalidFormat,
		},
		InDynamicField: "foo",
		ExpDesc:        "",
		ExpStatus:      "FooInvalidFormat",
	},
	{
		In:             NewNotFound("bar"),
		InDynamicField: "foo",
		ExpDesc:        "The resource foo was not found",
		ExpStatus:      "FooNotFound",
	},
	{
		In:             NewRemoteCall("bar.com"),
		InDynamicField: "foo.org",
		ExpDesc:        "Failed to call external resource [foo.org]",
		ExpStatus:      "FooOrgFailedRemoteCall",
	},
	{
		In:             NewRequired("bar"),
		InDynamicField: "foo",
		ExpDesc:        "The property foo is required",
		ExpStatus:      "FooRequired",
	},
	{
		In:             Error{},
		InDynamicField: "",
		ExpDesc:        "",
		ExpStatus:      "",
	},
}

func TestDynamicFields(t *testing.T) {
	for _, tt := range dynamicFieldTests {
		t.Run("", func(t *testing.T) {
			err := tt.In.SetProperty(tt.InDynamicField)
			assert.Equal(t, tt.ExpDesc, err.Description())
			assert.Equal(t, tt.ExpStatus, err.Status())
		})
	}
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
			statusName:  "FailedRemoteCall",
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
			statusName:  "FailedRemoteCall",
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
			statusName:  "NotFound",
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
			statusName:  "FooNotFound",
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
			statusName:  "AlreadyExists",
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
			statusName:  "FooAlreadyExists",
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
			statusName:  "OutOfRange",
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
			statusName:  "FooOutOfRange",
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
			statusName:  "OutOfRange",
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
			statusName:  "FooOutOfRange",
		},
	},
}

func TestNewOutOfRange(t *testing.T) {
	for _, tt := range newOutOfRangeTestSuite {
		t.Run("", func(t *testing.T) {
			err := NewOutOfRange(tt.InProp, tt.InLimA, tt.InLimB)
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
			statusName:  "InvalidFormat",
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
			statusName:  "InvalidFormat",
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
			statusName:  "FooInvalidFormat",
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
			statusName:  "FooInvalidFormat",
		},
	},
}

func TestNewInvalidFormat(t *testing.T) {
	for _, tt := range newInvalidFormatTestSuite {
		t.Run("", func(t *testing.T) {
			err := NewInvalidFormat(tt.InProperty, tt.InFormats...)
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
			statusName:  "IsRequired",
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
			statusName:  "FooIsRequired",
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
