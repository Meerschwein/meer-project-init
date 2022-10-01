package types

import (
	"fmt"
	"strings"
)

type Errors struct {
	errors []string
}

func (errs *Errors) append(new_error string) {
	errs.errors = append(errs.errors, new_error)
}

func (errs *Errors) appendf(format string, a ...interface{}) {
	errs.append(fmt.Sprintf(format, a...))
}

func (errs Errors) Error() string {
	return strings.Join(errs.errors, "\n")
}

func (errs Errors) len() int {
	return len(errs.errors)
}
