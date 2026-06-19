package exec

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

// Result describes the outcome of a command execution.
type Result struct {
	// Stdout contains captured standard output.
	Stdout    []byte
	// Stderr contains captured standard error.
	Stderr    []byte
	// Code is the process exit code.
	Code      int
	// FileName is the resolved executable path.
	FileName  string
	// Args contains the executed arguments.
	Args      []string
	// StartedAt is the UTC start time.
	StartedAt time.Time
	// EndedAt is the UTC end time.
	EndedAt   time.Time
	// TempFile points to any temporary file associated with execution.
	TempFile  *string
}

// Text returns stdout as a string.
func (o *Result) Text() string {
	return string(o.Stdout)
}

// IsOk reports whether the command exited successfully.
func (o *Result) IsOk() bool {
	return o.Code == 0
}

// ToError converts a failed result to an error.
func (o *Result) ToError() error {
	if o.IsOk() {
		return nil
	}
	return fmt.Errorf("command %s failed with code %d: %s", o.FileName, o.Code, o.ErrorText())
}

// ToErrorIf converts the result to an error when the callback returns true.
func (o *Result) ToErrorIf(f func(o *Result) bool) error {
	if f == nil {
		f = func(o *Result) bool {
			return o.Code != 0
		}
	}

	if f(o) {
		return o.ToError()
	}
	return nil
}

// Lines returns stdout split by the platform end-of-line value.
func (o *Result) Lines() []string {
	r := bytes.Split(o.Stdout, []byte(EOL))
	lines := []string{}
	for _, line := range r {
		lines = append(lines, string(line))
	}
	return lines
}

// ErrorText returns stderr as a string.
func (o *Result) ErrorText() string {
	return string(o.Stderr)
}

// ErrorLines returns stderr split by the platform end-of-line value.
func (o *Result) ErrorLines() []string {

	r := bytes.Split(o.Stdout, []byte(EOL))
	lines := []string{}
	for _, line := range r {
		lines = append(lines, string(line))
	}
	return lines
}

// Json unmarshals stdout as JSON.
func (o *Result) Json() (interface{}, error) {
	var out interface{}
	err := json.Unmarshal([]byte(o.Stdout), &out)
	return out, err
}

// ErrorJson unmarshals stderr as JSON.
func (o *Result) ErrorJson() (interface{}, error) {
	var out interface{}
	err := json.Unmarshal([]byte(o.Stderr), &out)
	return out, err
}

// Validate reports whether the result succeeded using the default validator.
func (o *Result) Validate() (bool, error) {
	return o.ValidateWith(nil)
}

// ValidateWith validates the result using a custom callback.
func (o *Result) ValidateWith(cb func(o *Result) (bool, error)) (bool, error) {
	if cb == nil {
		cb = func(o *Result) (bool, error) {
			if o.Code != 0 {
				return false, fmt.Errorf("command %s failed with code %d", o.FileName, o.Code)
			}

			return true, nil
		}
	}

	return cb(o)
}
