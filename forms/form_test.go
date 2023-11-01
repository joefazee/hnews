package forms

import (
	"reflect"
	"testing"
)

func TestErrs_Add(t *testing.T) {

	field := "first_field"
	msg := "error message"

	testErrors := make(errs)
	testErrors.Add(field, msg)

	res, ok := testErrors["first_field"]

	if !ok {
		t.Errorf("expected %v to be in the error map;", field)
	}

	if len(res) == 0 || res[0] != msg {
		t.Errorf("expected %v message to be in the error map", msg)
	}
}

func TestErrs_First(t *testing.T) {

	field := "first_field"
	msg := "error message"

	testErrors := make(errs)
	testErrors.Add(field, msg)
	testErrors.Add(field, "another message")
	res := testErrors.First(field)

	if res != msg {
		t.Errorf("expected error message to be %v, got; %v", msg, res)
	}
}

func TestNew(t *testing.T) {

	newF := New(nil)
	f := &Form{Errors: make(errs)}

	if !reflect.DeepEqual(newF, f) {
		t.Errorf("expected New() to return &Form{}; got %v", newF)
	}
}

func TestForm_Email(t *testing.T) {

	testCases := []struct {
		name  string
		field string
		value string
		want  string
	}{
		{"valid email", "valid_email", "sample@sample.com", ""},
		{"invalid email", "invalid_email", "sample sample.com", "the value provided is not a valid email address"},
		{"invalid email", "invalid_email", "", "the value provided is not a valid email address"},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			formData := make(map[string][]string)
			formData[tt.field] = append(formData[tt.field], tt.value)
			f := New(formData)
			f.Email(tt.field)
			got := f.Errors.First(tt.field)

			if got != tt.want {
				t.Errorf("expected %v; got %v", tt.want, got)
			}
		})
	}

}
