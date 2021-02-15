package ddderr

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var errPackageTestingSuite = []struct {
	inRoot error
	inErr  error
	exp    bool
}{
	{NewInfrastructure(nil, "custom type"), NewInfrastructure(nil, "custom type"), true},
	{NewInfrastructure(nil, "custom type"), NewInfrastructure(nil, "custom type 2"), false},
	{NewInfrastructure(nil, "custom type"), NewDomain("custom type", "generic description"), false},
	{NewDomain("custom type", "description"), NewDomain("custom type", "generic description"), false},
	{NewDomain("custom type", "description"), NewDomain("custom type", "description"), true},
}

func TestErrorsPackageCompat(t *testing.T) {
	for _, tt := range errPackageTestingSuite {
		t.Run("Errors package compatibility", func(t *testing.T) {
			assert.Equal(t, errors.Is(tt.inRoot, tt.inErr), tt.exp)
		})
	}
}

var errDomainTestingSuite = []struct {
	inErr    string
	inEntity string
}{
	{"", ""},
	{"i am a custom error", ""},
	{"", "foo"},
	{"i am a custom error", "foo"},
	{"entity (foo) not found", "foo"},
}

func TestNewDomain(t *testing.T) {
	for _, tt := range errDomainTestingSuite {
		t.Run("New domain generic error", func(t *testing.T) {
			err := NewDomain(tt.inEntity, tt.inErr)
			assert.Equal(t, err.Error(), tt.inErr)
			assert.Equal(t, err.Entity(), tt.inEntity)
			assert.True(t, err.IsDomain())
			assert.False(t, err.IsInfrastructure())
			assert.False(t, err.IsNotFound())
			assert.False(t, err.IsNotFound())
			assert.False(t, err.IsInvalidFormat())
			assert.False(t, err.IsAlreadyExists())
			assert.False(t, err.IsOutOfRange())
			assert.False(t, err.IsRequired())
			assert.False(t, err.IsRemoteCall())
		})
	}
}

var errInfraTestingSuite = []struct {
	inErr     string
	inParent  error
	expParent string
}{
	{"", nil, ""},
	{"i am a custom error", errors.New("root error"), "root error"},
}

func TestNewInfrastructure(t *testing.T) {
	for _, tt := range errInfraTestingSuite {
		t.Run("New infrastructure generic error", func(t *testing.T) {
			err := NewInfrastructure(tt.inParent, tt.inErr)
			assert.Equal(t, err.Error(), tt.inErr)
			assert.True(t, err.IsInfrastructure())
			if err.Parent() != nil {
				assert.Equal(t, tt.expParent, err.Parent().Error())
			}
			assert.False(t, err.IsDomain())
			assert.False(t, err.IsNotFound())
			assert.False(t, err.IsInvalidFormat())
			assert.False(t, err.IsAlreadyExists())
			assert.False(t, err.IsOutOfRange())
			assert.False(t, err.IsRequired())
			assert.False(t, err.IsRemoteCall())
		})
	}
}

var errAlreadyExistsTestingSuite = []struct {
	expErr   string
	inEntity string
}{
	{"already exists", ""},
	{"foo already exists", "foo"},
	{"bar already exists", "bar"},
}

func TestNewAlreadyExists(t *testing.T) {
	for _, tt := range errAlreadyExistsTestingSuite {
		t.Run("New already exists error", func(t *testing.T) {
			err := NewAlreadyExists(tt.inEntity)
			assert.Equal(t, err.Error(), tt.expErr)
			assert.Equal(t, err.Entity(), tt.inEntity)
			assert.True(t, err.IsDomain())
			assert.True(t, err.IsAlreadyExists())
			assert.False(t, err.IsInfrastructure())
			assert.False(t, err.IsInvalidFormat())
			assert.False(t, err.IsNotFound())
			assert.False(t, err.IsOutOfRange())
			assert.False(t, err.IsRequired())
			assert.False(t, err.IsRemoteCall())
		})
	}
}

var errInvalidFmtTestingSuite = []struct {
	expErr     string
	inEntity   string
	inExpTypes []string
}{
	{"invalid format, expected []", "", []string{}},
	{"foo contains an invalid format, expected []", "foo", []string{}},
	{"bar contains an invalid format, expected []", "bar", []string{}},
	{"baz contains an invalid format, expected []", "baz", []string{""}},
	{"foo contains an invalid format, expected [string]", "foo", []string{"string"}},
	{"foo contains an invalid format, expected [string,bool]", "foo", []string{"string", "bool"}},
	{"foo contains an invalid format, expected [string,bool,custom_type]", "foo",
		[]string{"string", "bool", "custom_type"}},
}

func TestNewInvalidFormat(t *testing.T) {
	for _, tt := range errInvalidFmtTestingSuite {
		t.Run("New invalid format error", func(t *testing.T) {
			err := NewInvalidFormat(tt.inEntity, tt.inExpTypes...)
			assert.Equal(t, err.Error(), tt.expErr)
			assert.Equal(t, err.Entity(), tt.inEntity)
			assert.True(t, err.IsInvalidFormat())
			assert.True(t, err.IsDomain())
			assert.False(t, err.IsInfrastructure())
			assert.False(t, err.IsAlreadyExists())
			assert.False(t, err.IsNotFound())
			assert.False(t, err.IsOutOfRange())
			assert.False(t, err.IsRequired())
			assert.False(t, err.IsRemoteCall())
		})
	}
}

var errOutOfRangeTestingSuite = []struct {
	inEntity string
	intA     int
	intB     int
	expErr   string
}{
	{"", 0, 0, "out of range [0,0)"},
	{"", 1, 100, "out of range [1,100)"},
	{"", 0, 75, "out of range [0,75)"},
	{"", -10, -50, "out of range [-10,-50)"},
	{"foo", -10, -50, "foo is out of range [-10,-50)"},
	{"foo", 1, 50, "foo is out of range [1,50)"},
}

func TestNewOutOfRange(t *testing.T) {
	for _, tt := range errOutOfRangeTestingSuite {
		t.Run("New out of range error", func(t *testing.T) {
			err := NewOutOfRange(tt.inEntity, tt.intA, tt.intB)
			assert.Equal(t, err.Error(), tt.expErr)
			assert.Equal(t, err.Entity(), tt.inEntity)
			assert.True(t, err.IsDomain())
			assert.True(t, err.IsOutOfRange())
			assert.False(t, err.IsInfrastructure())
			assert.False(t, err.IsInvalidFormat())
			assert.False(t, err.IsNotFound())
			assert.False(t, err.IsAlreadyExists())
			assert.False(t, err.IsRequired())
			assert.False(t, err.IsRemoteCall())
		})
	}
}

var errNotFoundTestingSuite = []struct {
	inEntity string
	expErr   string
}{
	{"", "not found"},
	{"foo", "foo not found"},
	{"bar", "bar not found"},
}

func TestNewNotFound(t *testing.T) {
	for _, tt := range errNotFoundTestingSuite {
		t.Run("New not found error", func(t *testing.T) {
			err := NewNotFound(tt.inEntity)
			assert.Equal(t, err.Error(), tt.expErr)
			assert.Equal(t, err.Entity(), tt.inEntity)
			assert.True(t, err.IsDomain())
			assert.True(t, err.IsNotFound())
			assert.False(t, err.IsInfrastructure())
			assert.False(t, err.IsInvalidFormat())
			assert.False(t, err.IsAlreadyExists())
			assert.False(t, err.IsOutOfRange())
			assert.False(t, err.IsRequired())
			assert.False(t, err.IsRemoteCall())
		})
	}
}

var errRequiredTestingSuite = []struct {
	inEntity string
	expErr   string
}{
	{"", "required"},
	{"foo", "foo is required"},
	{"bar", "bar is required"},
}

func TestNewRequired(t *testing.T) {
	for _, tt := range errRequiredTestingSuite {
		t.Run("New required error", func(t *testing.T) {
			err := NewRequired(tt.inEntity)
			assert.Equal(t, err.Error(), tt.expErr)
			assert.Equal(t, err.Entity(), tt.inEntity)
			assert.True(t, err.IsDomain())
			assert.True(t, err.IsRequired())
			assert.False(t, err.IsInfrastructure())
			assert.False(t, err.IsInvalidFormat())
			assert.False(t, err.IsAlreadyExists())
			assert.False(t, err.IsOutOfRange())
			assert.False(t, err.IsNotFound())
			assert.False(t, err.IsRemoteCall())
		})
	}
}

var errRemoteCallTestingSuite = []struct {
	in        string
	inParent  error
	exp       string
	expParent string
}{
	{"", nil, "failed to call external resource", ""},
	{"", errors.New("root error"), "failed to call external resource", "root error"},
	{"foo", nil, "failed to call external resource [foo]", ""},
	{"foo", errors.New("root error"), "failed to call external resource [foo]", "root error"},
	{"bar", errors.New("root error"), "failed to call external resource [bar]", "root error"},
}

func TestNewRemoteCall(t *testing.T) {
	for _, tt := range errRemoteCallTestingSuite {
		t.Run("New remote call error", func(t *testing.T) {
			err := NewRemoteCall(tt.inParent, tt.in)
			assert.Equal(t, err.Error(), tt.exp)
			if err.Parent() != nil {
				assert.Equal(t, err.Parent().Error(), tt.expParent)
			}
			assert.Empty(t, err.Entity())
			assert.True(t, err.IsInfrastructure())
			assert.True(t, err.IsRemoteCall())
			assert.False(t, err.IsDomain())
			assert.False(t, err.IsInvalidFormat())
			assert.False(t, err.IsAlreadyExists())
			assert.False(t, err.IsOutOfRange())
			assert.False(t, err.IsNotFound())
			assert.False(t, err.IsRequired())
		})
	}
}
