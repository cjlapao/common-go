// Copyright (c) 2015 Peter Strøiman, distributed under the MIT license

package automapper

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPanicWhenDestIsNotPointer(t *testing.T) {
	defer func() { recover() }()
	source, dest := SourceTypeA{}, DestTypeA{}
	err := Map(source, dest)

	assert.NotNilf(t, err, "Should have generated an error")
}

func TestDestinationIsUpdatedFromSource(t *testing.T) {
	source, dest := SourceTypeA{Foo: 42}, DestTypeA{}
	err := Map(source, &dest)
	assert.Equal(t, 42, dest.Foo)
	assert.Nilf(t, err, "Should not have generated an error")
}

func TestDestinationIsUpdatedFromSourceWhenSourcePassedAsPtr(t *testing.T) {
	source, dest := SourceTypeA{42, "Bar"}, DestTypeA{}
	err := Map(&source, &dest)
	assert.Equal(t, 42, dest.Foo)
	assert.Equal(t, "Bar", dest.Bar)
	assert.Nilf(t, err, "Should not have generated an error")
}

func TestWithNestedTypes(t *testing.T) {
	source := struct {
		Baz   string
		Child SourceTypeA
	}{}
	dest := struct {
		Baz   string
		Child DestTypeA
	}{}

	source.Baz = "Baz"
	source.Child.Bar = "Bar"
	err := Map(&source, &dest)
	assert.Equal(t, "Baz", dest.Baz)
	assert.Equal(t, "Bar", dest.Child.Bar)
	assert.Nilf(t, err, "Should not have generated an error")
}

func TestWithSourceSecondLevel(t *testing.T) {
	source := struct {
		Child DestTypeA
	}{}
	dest := SourceTypeA{}

	source.Child.Bar = "Bar"
	err := Map(&source, &dest)
	assert.Equal(t, "Bar", dest.Bar)
	assert.Nilf(t, err, "Should not have generated an error")
}

func TestWithDestSecondLevel(t *testing.T) {
	source := SourceTypeA{}
	dest := struct {
		Child DestTypeA
	}{}

	source.Bar = "Bar"
	err := Map(&source, &dest)
	assert.Equal(t, "Bar", dest.Child.Bar)
	assert.Nilf(t, err, "Should not have generated an error")
}

func TestWithSliceTypes(t *testing.T) {
	source := struct {
		Children []SourceTypeA
	}{}
	dest := struct {
		Children []DestTypeA
	}{}
	source.Children = []SourceTypeA{
		SourceTypeA{Foo: 1},
		SourceTypeA{Foo: 2}}

	err := Map(&source, &dest)
	assert.Equal(t, 1, dest.Children[0].Foo)
	assert.Equal(t, 2, dest.Children[1].Foo)
	assert.Nilf(t, err, "Should not have generated an error")
}

func TestWithMultiLevelSlices(t *testing.T) {
	source := struct {
		Parents []SourceParent
	}{}
	dest := struct {
		Parents []DestParent
	}{}
	source.Parents = []SourceParent{
		SourceParent{
			Children: []SourceTypeA{
				SourceTypeA{Foo: 42},
				SourceTypeA{Foo: 43},
			},
		},
		SourceParent{
			Children: []SourceTypeA{},
		},
	}

	err := Map(&source, &dest)
	assert.Nilf(t, err, "Should not have generated an error")
}

func TestWithEmptySliceAndIncompatibleTypes(t *testing.T) {
	defer func() { recover() }()

	source := struct {
		Children []struct{ Foo string }
	}{}
	dest := struct {
		Children []struct{ Bar int }
	}{}

	err := Map(&source, &dest)
	assert.NotNilf(t, err, "Should have generated an error")
}

func TestWhenSourceIsMissingField(t *testing.T) {
	defer func() { recover() }()
	source := struct {
		A string
	}{}
	dest := struct {
		A, B string
	}{}
	err := Map(&source, &dest)
	assert.NotNilf(t, err, "Should have generated an error")
}

func TestWithUnnamedFields(t *testing.T) {
	source := struct {
		Baz string
		SourceTypeA
	}{}
	dest := struct {
		Baz string
		DestTypeA
	}{}
	source.Baz = "Baz"
	source.SourceTypeA.Foo = 42

	err := Map(&source, &dest)
	assert.Equal(t, "Baz", dest.Baz)
	assert.Equal(t, 42, dest.DestTypeA.Foo)
	assert.Nilf(t, err, "Should not have generated an error")
}

func TestWithPointerFieldsNotNil(t *testing.T) {
	source := struct {
		Foo *SourceTypeA
	}{}
	dest := struct {
		Foo *DestTypeA
	}{}
	source.Foo = nil

	err := Map(&source, &dest)
	assert.Nil(t, dest.Foo)
	assert.Nilf(t, err, "Should not have generated an error")
}

func TestWithPointerFieldsNil(t *testing.T) {
	source := struct {
		Foo *SourceTypeA
	}{}
	dest := struct {
		Foo *DestTypeA
	}{}
	source.Foo = &SourceTypeA{Foo: 42}

	err := Map(&source, &dest)
	assert.NotNil(t, dest.Foo)
	assert.Equal(t, 42, dest.Foo.Foo)
	assert.Nilf(t, err, "Should not have generated an error")
}

func TestMapFromPointerToNonPointerTypeWithData(t *testing.T) {
	source := struct {
		Foo *SourceTypeA
	}{}
	dest := struct {
		Foo DestTypeA
	}{}
	source.Foo = &SourceTypeA{Foo: 42}

	err := Map(&source, &dest)
	assert.NotNil(t, dest.Foo)
	assert.Equal(t, 42, dest.Foo.Foo)
	assert.Nilf(t, err, "Should not have generated an error")
}

func TestMapFromPointerToNonPointerTypeWithoutData(t *testing.T) {
	source := struct {
		Foo *SourceTypeA
	}{}
	dest := struct {
		Foo DestTypeA
	}{}
	source.Foo = nil

	err := Map(&source, &dest)
	assert.NotNil(t, dest.Foo)
	assert.Equal(t, 0, dest.Foo.Foo)
	assert.Nilf(t, err, "Should not have generated an error")
}

func TestMapFromPointerToAnonymousTypeToFieldName(t *testing.T) {
	source := struct {
		*SourceTypeA
	}{}
	dest := struct {
		Foo int
	}{}
	source.SourceTypeA = nil

	err := Map(&source, &dest)
	assert.Equal(t, 0, dest.Foo)
	assert.Nilf(t, err, "Should not have generated an error")
}

func TestMapFromPointerToNonPointerTypeWithoutDataAndIncompatibleType(t *testing.T) {
	defer func() { recover() }()
	// Just make sure we stil panic
	source := struct {
		Foo *SourceTypeA
	}{}
	dest := struct {
		Foo struct {
			Baz string
		}
	}{}
	source.Foo = nil

	err := Map(&source, &dest)
	assert.NotNilf(t, err, "Should have generated an error")
}

func TestWhenUsingIncompatibleTypes(t *testing.T) {
	defer func() { recover() }()
	source := struct{ Foo string }{}
	dest := struct{ Foo int }{}
	err := Map(&source, &dest)
	assert.NotNilf(t, err, "Should have generated an error")
}

func TestWithLooseOption(t *testing.T) {
	source := struct {
		Foo string
		Baz int
	}{"Foo", 42}
	dest := struct {
		Foo string
		Bar int
	}{}
	err := Map(&source, &dest, Loose)
	assert.Equal(t, dest.Foo, "Foo")
	assert.Equal(t, dest.Bar, 0)
	assert.Nilf(t, err, "Should not have generated an error")
}

func TestWithRequestFormOption(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("foo=bar&bar=2"))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	dest := struct {
		Foo string
		Bar int
	}{}
	err := Map(r, &dest, RequestForm)
	assert.Equal(t, dest.Foo, "bar")
	assert.Equal(t, dest.Bar, 2)
	assert.Nilf(t, err, "Should not have generated an error")
}

func TestWithRequestFormOptionWithoutARequestPointer(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("testFoo=bar&bar=2"))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	dest := struct {
		Foo string `json:"testFoo"`
		Bar int    `json:"bar"`
	}{}

	err := Map(&r, &dest, RequestForm)
	assert.NotNilf(t, err, "Should have generated an error")
}
func TestWithRequestFormJsonOption(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("testFoo=bar&bar=2"))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	dest := struct {
		Foo string `json:"testFoo"`
		Bar int    `json:"bar"`
	}{}

	err := Map(r, &dest, RequestFormWithJsonTag)
	assert.Equal(t, dest.Foo, "bar")
	assert.Equal(t, dest.Bar, 2)
	assert.Nilf(t, err, "Should not have generated an error")
}

func TestWithRequestFormJsonOptionWithoutARequestPointer(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("testFoo=bar&bar=2"))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	dest := struct {
		Foo string `json:"testFoo"`
		Bar int    `json:"bar"`
	}{}

	err := Map(&r, &dest, RequestFormWithJsonTag)
	assert.NotNilf(t, err, "Should have generated an error")
}

type SourceParent struct {
	Children []SourceTypeA
}

type DestParent struct {
	Children []DestTypeA
}

type SourceTypeA struct {
	Foo int
	Bar string
}

type DestTypeA struct {
	Foo int
	Bar string
}
