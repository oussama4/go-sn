package forms

// Errors holds  a list of error messages for each field
type Errors map[string][]string

// Add adds and error message to the errors map
func (e Errors) Add(field, msg string) {
	e[field] = append(e[field], msg)
}

// Get a list of errors for the given field
func (e Errors) Get(field string) []string {
	es := e[field]
	if len(es) == 0 {
		return []string{}
	}
	return es
}
