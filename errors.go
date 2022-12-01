package asahi

import (
	"fmt"
)

type RequestError struct {
	Function string
	StatusCode int
	Err error
}

func (r *RequestError) Error() string {
	return fmt.Sprintf("%s(): responded with status %d: %s", r.Function, r.StatusCode, r.Err)
}

func (r *RequestError) Code() int {
	return r.StatusCode
}
