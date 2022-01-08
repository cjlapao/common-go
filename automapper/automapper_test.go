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
	Map(source, dest)

	t.Error("Should have panicked")
}

func TestDestinationIsUpdatedFromSource(t *testing.T) {
	source, dest := SourceTypeA{Foo: 42}, DestTypeA{}
	Map(source, &dest)
	assert.Equal(t, 42, dest.Foo)
}

func TestDestinationIsUpdatedFromSourceWhenSourcePassedAsPtr(t *testing.T) {
	source, dest := SourceTypeA{42, "Bar"}, DestTypeA{}
	Map(&source, &dest)
	assert.Equal(t, 42, dest.Foo)
	assert.Equal(t, "Bar", dest.Bar)
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
	Map(&source, &dest)
	assert.Equal(t, "Baz", dest.Baz)
	assert.Equal(t, "Bar", dest.Child.Bar)
}

func TestWithSourceSecondLevel(t *testing.T) {
	source := struct {
		Child DestTypeA
	}{}
	dest := SourceTypeA{}

	source.Child.Bar = "Bar"
	Map(&source, &dest)
	assert.Equal(t, "Bar", dest.Bar)
}

func TestWithDestSecondLevel(t *testing.T) {
	source := SourceTypeA{}
	dest := struct {
		Child DestTypeA
	}{}

	source.Bar = "Bar"
	Map(&source, &dest)
	assert.Equal(t, "Bar", dest.Child.Bar)
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

	Map(&source, &dest)
	assert.Equal(t, 1, dest.Children[0].Foo)
	assert.Equal(t, 2, dest.Children[1].Foo)
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

	Map(&source, &dest)
}

func TestWithEmptySliceAndIncompatibleTypes(t *testing.T) {
	defer func() { recover() }()

	source := struct {
		Children []struct{ Foo string }
	}{}
	dest := struct {
		Children []struct{ Bar int }
	}{}

	Map(&source, &dest)
	t.Error("Should have panicked")
}

func TestWhenSourceIsMissingField(t *testing.T) {
	defer func() { recover() }()
	source := struct {
		A string
	}{}
	dest := struct {
		A, B string
	}{}
	Map(&source, &dest)
	t.Error("Should have panicked")
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

	Map(&source, &dest)
	assert.Equal(t, "Baz", dest.Baz)
	assert.Equal(t, 42, dest.DestTypeA.Foo)
}

func TestWithPointerFieldsNotNil(t *testing.T) {
	source := struct {
		Foo *SourceTypeA
	}{}
	dest := struct {
		Foo *DestTypeA
	}{}
	source.Foo = nil

	Map(&source, &dest)
	assert.Nil(t, dest.Foo)
}

func TestWithPointerFieldsNil(t *testing.T) {
	source := struct {
		Foo *SourceTypeA
	}{}
	dest := struct {
		Foo *DestTypeA
	}{}
	source.Foo = &SourceTypeA{Foo: 42}

	Map(&source, &dest)
	assert.NotNil(t, dest.Foo)
	assert.Equal(t, 42, dest.Foo.Foo)
}

func TestMapFromPointerToNonPointerTypeWithData(t *testing.T) {
	source := struct {
		Foo *SourceTypeA
	}{}
	dest := struct {
		Foo DestTypeA
	}{}
	source.Foo = &SourceTypeA{Foo: 42}

	Map(&source, &dest)
	assert.NotNil(t, dest.Foo)
	assert.Equal(t, 42, dest.Foo.Foo)
}

func TestMapFromPointerToNonPointerTypeWithoutData(t *testing.T) {
	source := struct {
		Foo *SourceTypeA
	}{}
	dest := struct {
		Foo DestTypeA
	}{}
	source.Foo = nil

	Map(&source, &dest)
	assert.NotNil(t, dest.Foo)
	assert.Equal(t, 0, dest.Foo.Foo)
}

func TestMapFromPointerToAnonymousTypeToFieldName(t *testing.T) {
	source := struct {
		*SourceTypeA
	}{}
	dest := struct {
		Foo int
	}{}
	source.SourceTypeA = nil

	Map(&source, &dest)
	assert.Equal(t, 0, dest.Foo)
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

	Map(&source, &dest)
	t.Error("Should have panicked")
}

func TestWhenUsingIncompatibleTypes(t *testing.T) {
	defer func() { recover() }()
	source := struct{ Foo string }{}
	dest := struct{ Foo int }{}
	Map(&source, &dest)
	t.Error("Should have panicked")
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
	Map(&source, &dest, Loose)
	assert.Equal(t, dest.Foo, "Foo")
	assert.Equal(t, dest.Bar, 0)
}

func TestWithRequestFormOption(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("foo=bar&bar=2"))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	dest := struct {
		Foo string
		Bar int
	}{}
	Map(r, &dest, RequestForm)
	assert.Equal(t, dest.Foo, "bar")
	assert.Equal(t, dest.Bar, 2)
}

func TestWithRequestFormJsonOption(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("testFoo=bar&bar=2"))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	dest := struct {
		Foo string `json:"testFoo"`
		Bar int    `json:"bar"`
	}{}

	Map(r, &dest, RequestFormWithJsonTag)
	assert.Equal(t, dest.Foo, "bar")
	assert.Equal(t, dest.Bar, 2)
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
