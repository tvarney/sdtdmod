package errors

import (
	"fmt"
	"strings"

	"github.com/tvarney/maputil/errctx"
	"github.com/tvarney/maputil/mpath"
)

// ErrorCollector is an ErrorHandler which collects errors.
type ErrorCollector struct {
	Errors []error
}

func (e *ErrorCollector) Add(p *mpath.Path, err error) {
	e.Errors = append(e.Errors, &ErrorWithContext{Err: err, Path: p.Copy()})
}

var _ errctx.ErrorHandler = &ErrorCollector{}

// ErrorWithContext is an error type which combines an error with a path.
type ErrorWithContext struct {
	Err  error
	Path *mpath.Path
}

func (e *ErrorWithContext) Error() string {
	return fmt.Sprintf("%v: %v", e.Path, e.Err)
}

func (e *ErrorWithContext) Unwrap() error {
	return e.Err
}

// ParseCountError is an error indicating that a number of errors occured
// while parsing.
type ParseCountError struct {
	Errors []error
	Name   string
	Count  int
}

func (e *ParseCountError) Error() string {
	if len(e.Errors) == 0 {
		return fmt.Sprintf("%d errors while parsing %q", e.Count, e.Name)
	}

	b := strings.Builder{}
	fmt.Fprintf(&b, "%d errors while parsing %q", e.Count, e.Name)
	for _, e := range e.Errors {
		b.WriteRune('\n')
		b.WriteString(e.Error())
	}
	return b.String()
}

type UnexpectedKeyError string

func (e UnexpectedKeyError) Error() string {
	return fmt.Sprintf("unexpected key %q", string(e))
}
