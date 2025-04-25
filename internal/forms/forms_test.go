package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Error("got invalid when should have been valid")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)
	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form show svalid when required fields are missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "a")
	postedData.Add("c", "a")

	r, _ = http.NewRequest("POST", "/whatever", nil)

	r.PostForm = postedData
	form = New(r.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("shows form does not have required fields when it does")
	}
}

func TestForm_Has(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	has := form.Has("whatever")
	if has {
		t.Error("form shows has field when it does not")
	}

	postedData = url.Values{}
	postedData.Add("a", "a")

	form = New(postedData)

	has = form.Has("a")
	if !has {
		t.Error("shows form does not have field when it should")
	}
}

func TestForm_MinLength(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	form.MinLength("x", 10)
	if form.Valid() {
		t.Error("form shows minlength for non-existent field")
	}

	isError := form.Errors.Get("x")
	if isError == "" {
		t.Error("should have an error, but did not get one")
	}

	postedData := url.Values{}
	postedData.Add("some-field", "some-value")
	form = New(postedData)

	form.MinLength("some-field", 100)
	if form.Valid() {
		t.Error("shows minlength of 100 met when data is shorter")
	}

	postedData = url.Values{}
	postedData.Add("another-field", "abc123")
	form = New(postedData)

	form.MinLength("another-field", 1)
	if !form.Valid() {
		t.Error("shows form doesn't meet minlength when it does")
	}

	isError = form.Errors.Get("another-field")
	if isError != "" {
		t.Error("received an error when we shouldn't have")
	}
}

func TestForm_IsEmail(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	form.IsEmail("a")
	if form.Valid() {
		t.Error("form shows valid email for non-existent field")
	}

	postedData = url.Values{}
	postedData.Add("email", "me@exmaple.com")
	form = New(postedData)

	form.IsEmail("email")
	if !form.Valid() {
		t.Error("shows email isn't valid when it is")
	}

	postedData = url.Values{}
	postedData.Add("email", "a")
	form = New(postedData)

	form.IsEmail("email")
	if form.Valid() {
		t.Error("got valid for invalid email")
	}
}

// func TestHas(t *testing.T) {
// r := httptest.NewRequest("POST", "/whatever", nil)
// form = New(r.PostForm)

// form.Has("a", *r)
// if form.Valid() {
// t.Error("shows form has a field when it does not")
// }

// postedData := url.Values{}
// postedData.Add("a", "a")
// r.PostForm = postedData

// form.Has("a", *r)
// if !form.Valid() {
// t.Error("form should be valid but tests show it is not")
// }

// }
