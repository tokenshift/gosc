package gosc

import (
	"fmt"
	"io"
)

func WriteMessage(out io.Writer, address OSCAddressPattern, args...OSCArg) (int, error) {
	if err := address.Valid(); err != nil {
		return 0, err
	}

	return 0, fmt.Errorf("NOT IMPLEMENTED")
}
