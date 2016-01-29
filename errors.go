package gosc

import (
	"fmt"
)

type OSCArgumentError string

func OSCArgumentErrorf(f string, args...interface{}) OSCArgumentError {
	return OSCArgumentError(fmt.Sprintf(f, args...))
}

func (e OSCArgumentError) Error() string {
	return string(e)
}
