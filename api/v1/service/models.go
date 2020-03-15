package service

import "fmt"

type status string
type errorType string

const (
	statusSuccess status = "success"
	statusError   status = "error"
)

const (
	errorInternal errorType = "server_error"
	errorBadData  errorType = "bad_data"
)

// ApiError
type apiError struct {
	typ errorType
	err error
}

func (e *apiError) Error() string {
	return fmt.Sprintf("%s: %s", e.typ, e.err)
}
