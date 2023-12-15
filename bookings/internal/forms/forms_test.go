package forms

import (
	"net/http"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r, err := http.NewRequest("POST", "/urlsser", nil)
	if err != nil {
		t.Error(err)
	}

	form := New(r.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Error("got invalid when form should be valid")
	}

}

func TestForm_Required(t *testing.T) {
	r, err := http.NewRequest("POST", "/urlsser", nil)
	if err != nil {
		t.Error(err)
	}

	form := New(r.PostForm)

	form.Required("name", "age", "email")
	if form.Valid() {
		t.Error("required validation failed for unprovided fields")
	}

	postData := url.Values{}
	postData.Add("name", "pablo")
	postData.Add("age", "56")
	postData.Add("email", "email")

	form = New(postData)
	form.Required("name", "age", "email")
	if !form.Valid() {
		t.Error("required validation failed for provided fields")
	}
}

func TestNew(t *testing.T) {
	postData := url.Values{}
	postData.Add("name", "pablo")
	postData.Add("age", "56")
	postData.Add("email", "email")

	form := New(postData)

	if form.Get("name") != "pablo" {
		t.Error("failed to add postform values in fields")
	}
}

func TestHas(t *testing.T) {
	postData := url.Values{}
	postData.Add("name", "August")

	form := New(postData)

	if !form.Has("name") {
		t.Error("Has gives false when form has a value")
	}
	if form.Has("email") {
		t.Error("Has gives true when form value not present")
	}
}

func TestMinLength(t *testing.T) {

	postData := url.Values{}
	postData.Add("name", "August")
	postData.Add("email", "August@")

	form := New(postData)

	ml := form.MinLength("name", 3)
	if !ml || !form.Valid() {
		t.Error("minLength didn't check length properly")
	}
	ml = form.MinLength("email", 10)
	if ml || form.Valid() {
		t.Error("minLength didn't check length properly")
	}
}

func TestIsEmail(t *testing.T) {
	postData := url.Values{}
	postData.Add("email1", "August@")
	postData.Add("email2", "August@go.cee")

	form := New(postData)
	form.IsEmail("email1")
	form.IsEmail("email2")
	if form.Errors.Get("email1") == "" {
		t.Errorf("isEmail validator failed for value %s", form.Get("email1"))
	}
	if form.Errors.Get("email2") != "" {
		t.Errorf("isEmail validator failed for value %s", form.Get("email2"))
	}
}
