package gosc

import (
	"io"
)

func WriteMessage(out io.Writer, address OSCAddressPattern, args...OSCArg) (int, error) {
	if err := address.Valid(); err != nil {
		return 0, err
	}

	// Validate all of the arguments and construct the complete tag string
	// before sending anything.
	tagstring := make([]rune, len(args) + 1)
	tagstring[0] = ','
	for i, arg := range args {
		if err := arg.Valid(); err != nil {
			return 0, err
		}

		tagstring[i+1] = rune(arg.Tag())
	}

	total := 0

	if sent, err := address.WriteTo(out); err != nil {
		return sent, err
	} else {
		total += sent
	}

	if sent, err := OSCString(tagstring).WriteTo(out); err != nil {
		return sent, err
	} else {
		total += sent
	}

	for _, arg := range args {
		if sent, err := arg.WriteTo(out); err != nil {
			return total+sent, err
		} else {
			total += sent
		}
	}

	return total, nil
}
