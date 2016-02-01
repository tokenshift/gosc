package gosc

import (
	"io"
	"strings"
)

// Writes an OSC message to the output stream. Returns an error if any of the
// arguments were invalid, or if any transmission error occurred, and returns
// the total number of bytes sent in either case.
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

// Reads an OSC message from an input stream, returning the address and
// arguments, or an error if the message could not be read successfully.
func ReadMessage(in io.Reader) (OSCAddressPattern, []OSCArg, error) {
	address, err := ReadOSCString(in)
	if err != nil {
		return "", nil, err
	}

	oaddress := OSCAddressPattern(address)
	if err = oaddress.Valid(); err != nil {
		return "", nil, err
	}

	tagString, err := ReadOSCString(in)
	if err != nil {
		return oaddress, nil, err
	}
	if !strings.HasPrefix(string(tagString), ",") {
		return oaddress, nil, OSCReadErrorf("tag string (%s) must start with a comma", tagString)
	}
	if err = tagString.Valid(); err != nil {
		return oaddress, nil, err
	}

	args := make([]OSCArg, 0, len(tagString) - 1)
	for _, tag := range tagString[1:] {
		var arg OSCArg

		switch OSCTypeTag(tag) {
		case OSC_TYPE_STRING:
			arg, err = ReadOSCString(in)
		default:
			return oaddress, nil, OSCReadErrorf("unsupported type tag: %s", tag)
		}

		if err != nil {
			return oaddress, args, err
		}

		args = append(args, arg)
	}

	return oaddress, args, nil
}
