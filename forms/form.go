package forms

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type errs map[string][]string

func (e errs) Add(field, message string) {
	e[field] = append(e[field], message)
}

func (e errs) First(field string) string {
	es := e[field]
	if len(es) == 0 {
		return ""
	}

	return es[0]

}

type Form struct {
	url.Values
	Errors errs
}

func New(data url.Values) *Form {
	return &Form{
		data,
		make(errs),
	}
}

func (f *Form) Email(field string) *Form {
	value := f.Get(field)
	if value == "" || !emailRegex.MatchString(value) {
		f.Errors.Add(field, "the value provided is not a valid email address")
	}
	return f
}

func (f *Form) MinLength(field string, d int) *Form {
	value := f.Get(field)
	if utf8.RuneCountInString(value) < d {
		f.Errors.Add(field, fmt.Sprintf("This %s is too short (minimum is %d characters)", field, d))
	}
	return f
}

func (f *Form) MaxLength(field string, d int) *Form {
	value := f.Get(field)
	if value == "" {
		return f
	}
	if utf8.RuneCountInString(value) > d {
		f.Errors.Add(field, fmt.Sprintf("This %s is too large (maximum is %d characters)", field, d))
	}
	return f
}

func (f *Form) Required(fields ...string) *Form {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, fmt.Sprintf("%s is required", field))
		}
	}

	return f
}

func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

func (f *Form) Fail(field, msg string) {
	f.Errors.Add(field, msg)
}

func (f *Form) Url(field string) *Form {
	value := f.Get(field)

	u, err := url.Parse(value)
	if err != nil || u.Scheme == "" || u.Host == "" {
		f.Errors.Add(field, fmt.Sprintf("%s is not a valid url", field))
	}

	return f
}

func (f *Form) GetInt(field string) int {

	v, err := strconv.Atoi(f.Get(field))
	if err != nil {
		return 0
	}

	return v
}
