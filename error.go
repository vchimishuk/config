package config

import "fmt"

type parseError struct {
	line int
	msg  string
}

func newError(line int, format string, args ...any) *parseError {
	return &parseError{line, fmt.Sprintf(format, args...)}
}

func (e parseError) Error() string {
	return fmt.Sprintf("%d: %s", e.line, e.msg)
}
