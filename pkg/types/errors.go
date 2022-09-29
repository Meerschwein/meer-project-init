package types

import (
	"fmt"
	"strings"
)

type Errors struct {
	errors []string
}

func (errs *Errors) Append(new_error string) {
	errs.errors = append(errs.errors, new_error)
}

func (errs *Errors) Appendf(format string, a ...interface{}) {
	errs.errors = append(errs.errors, fmt.Sprintf(format, a...))
}

func (errs Errors) Error() string {
	return strings.Join(errs.errors, "\n")
}

func (errs Errors) Len() int {
	return len(errs.errors)
}
