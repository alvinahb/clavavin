package forms

import (
	"net/url"
	"strings"
)

// Forn creates a custom form struct, embeds a url.Values object
type Form struct {
	url.Values
	Errors errors
}

// New initializes a form struct
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Has checks if form field is in post and not empty
func (f *Form) Has(field string) bool {
	x := f.Get(field)
	if x == "" {
		return false
	}
	return true
}

// Required checks for required fields
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "Champ nécessaire")
		}
	}
}

// ContentIs checks if content of field is in values
func (f *Form) ContentIs(field string, values []string) {
	value := f.Get(field)
	for _, item := range values {
		if value == item {
			return
		}
	}
	f.Errors.Add(field, "Valeur invalide")
}

// Valid returns true if there are no errors, otherwise false
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
