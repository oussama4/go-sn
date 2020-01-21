package forms

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"
)

// EmailRx is a regular expression pattern of a valid email address
var EmailRx = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// Form holds form data and the errors associated with it
type Form struct {
	url.Values
	Errors Errors
}

// New creates a form struct given the form data
func New(data url.Values) *Form {
	return &Form{
		data,
		Errors(map[string][]string{}),
	}
}

// Required checks if a set of fields are not empty
// if any field is empty add an error message to form errors
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "this field is required")
		}
	}
}

// MaxLength checks the number of characters in a field don't pass a maximum value
func (f *Form) MaxLength(field string, max int) {
	value := f.Get(field)
	if value == "" {
		return
	}

	if utf8.RuneCountInString(value) > max {
		f.Errors.Add(field, fmt.Sprintf("%s must not pass %d characters", field, max))
	}
}

// MinLength checks  field has a minimum number of characters
func (f *Form) MinLength(field string, min int) {
	value := f.Get(field)
	if value == "" {
		return
	}

	if utf8.RuneCountInString(value) < min {
		f.Errors.Add(field, fmt.Sprintf("%s must have at least %d characters", field, min))
	}
}

// ValidEmail validate an email address
func (f *Form) ValidEmail(field string) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if !EmailRx.MatchString(value) {
		f.Errors.Add(field, "email is not valid")
	}
}

// StringsMatch checks if two strings are the same
func (f *Form) StringsMatch(field1, field2 string) {
	value1 := f.Get(field1)
	value2 := f.Get(field2)
	if value1 == "" && value2 == "" {
		return
	}
	if strings.TrimSpace(value1) != strings.TrimSpace(value2) {
		f.Errors.Add(field1, fmt.Sprintf("%s does not match %s", field1, field2))
	}
}

// Valid returns true if the form is valid
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
