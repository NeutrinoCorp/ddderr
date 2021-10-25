package ddderr

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var getSanitizedStatusNameTestSuite = []struct {
	InAttr string
	InOp   string
	Exp    string
}{
	{
		InAttr: "",
		InOp:   "",
		Exp:    "",
	},
	{
		InAttr: "-",
		InOp:   "",
		Exp:    "",
	},
	{
		InAttr: "-e",
		InOp:   "",
		Exp:    "E",
	},
	{
		InAttr: "foo-bar-baz",
		InOp:   "",
		Exp:    "FooBarBaz",
	},
	{
		InAttr: "foo-bar-baz",
		InOp:   "NotFound",
		Exp:    "FooBarBazNotFound",
	},
	{
		InAttr: "foo_bar_baz",
		InOp:   "NotFound",
		Exp:    "FooBarBazNotFound",
	},
	{
		InAttr: "_foo_bar_baz",
		InOp:   "NotFound",
		Exp:    "FooBarBazNotFound",
	},
	{
		InAttr: "foo#bar#baz",
		InOp:   "NotFound",
		Exp:    "FooBarBazNotFound",
	},
}

func TestGetSanitizedStatusName(t *testing.T) {
	for _, tt := range getSanitizedStatusNameTestSuite {
		t.Run("", func(t *testing.T) {
			status := getSanitizedStatusName(tt.InAttr, tt.InOp)
			assert.Equal(t, tt.Exp, status)
		})
	}
}

func BenchmarkGetSanitizedStatusName(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.ReportAllocs()
		getSanitizedStatusName("foo_bar_baz", "NotFound")
	}
}
