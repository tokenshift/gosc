package gosc

import (
	"fmt"
	"io"
)

type OSCTypeTag byte
const (
	// Standard argument types.
	OST_TYPE_INT32        = OSCTypeTag('i')
	OST_TYPE_FLOAT32      = OSCTypeTag('f')
	OST_TYPE_STRING       = OSCTypeTag('s')
	OST_TYPE_BLOB         = OSCTypeTag('b')

	// "Extended" argument types (not all clients support any/all of these).
	OST_ETYPE_INT64       = OSCTypeTag('h')
	OST_ETYPE_TIMETAG     = OSCTypeTag('t')
	OST_ETYPE_FLOAT64     = OSCTypeTag('d')
	OST_ETYPE_STRING_ALT  = OSCTypeTag('S')
	OST_ETYPE_CHAR        = OSCTypeTag('c')
	OST_ETYPE_RGBA        = OSCTypeTag('r')
	OST_ETYPE_MIDI        = OSCTypeTag('m')
	OST_ETYPE_TRUE        = OSCTypeTag('T')
	OST_ETYPE_FALSE       = OSCTypeTag('F')
	OST_ETYPE_NIL         = OSCTypeTag('N')
	OST_ETYPE_INFINITY    = OSCTypeTag('I')
	OST_ETYPE_ARRAY_START = OSCTypeTag('[')
	OST_ETYPE_ARRAY_END   = OSCTypeTag(']')
)

// Arguments are all OSC types that can be transitted in a message.
type OSCArg interface {
	// Writes the argument to an output stream. Returns the number of bytes
	// written and an error, if the argument could not be written.
	// This error will either be a serialization error (the argument could not
	// be turned into a stream of bytes, usually because it was invalid in some
	// way) or an underlying transmission failure returned by the output Writer.
	WriteTo(out io.Writer) (int, error)
	Tag() OSCTypeTag
}

// OSC-strings are more restrictive than go strings, so a []byte would be more
// appropriate; string is used purely for convenience, so that consumers of the
// library don't have to add conversion logic to their string literals. All
// strings are validated before transmission.
type OSCString string

func (s OSCString) Tag() OSCTypeTag {
	return OST_TYPE_STRING
}

func (s OSCString) Valid() error {
	// OSC-strings are null-terminated ASCII strings, padded to four bytes.
	// The null-termination and padding are left to the Write method; here we
	// simply validate that the string has no invalid characters.
	for i, r := range(s) {
		if int(r) == 0 || int(r) > 127 {
			return OSCArgumentErrorf("non-ascii character 0x%x found at position %d in string \"%s\"", r, i, s)
		}
	}

	return nil
}

func (s OSCString) WriteTo(out io.Writer) (int, error) {
	data := []byte(s)

	// OSC-strings are null-terminated.
	data = append(data, 0)

	// OSC-strings are null-padded to 4 bytes.
	for len(data) % 4 != 0 {
		data = append(data, 0)
	}

	return out.Write(data)
}

// An OSC address pattern is an OSC-string with some additional restrictions.
type OSCAddressPattern OSCString

func (s OSCAddressPattern) Valid() error {
	// An Address Pattern is an OSC-string starting with "/".
	if err := OSCString(s).Valid(); err != nil {
		return err
	}

	// At this point, the string has already been validated to contain only
	// ASCII characters, so it's safe to cast the first rune to a byte.
	if len(s) == 0 || s[0] != byte('/') {
		return OSCArgumentErrorf("OSCAddressPattern must start with a forward slash")
	}

	// Certain ASCII characters are disallowed in address patterns. Technically,
	// this is the list of disallowed characters for symbolic names, but I also
	// allow forward slashes here (the name separator).
	for i := 1; i < len(s); i++ {
		for _, invalid := range(" #*,?[]{}") {
			if s[i] == byte(invalid) {
				return OSCArgumentErrorf("disallowed character '%s' found at position %d in string \"%s\"", invalid, i, s)
			}
		}
	}

	return nil
}

type OSCArgumentError string
func OSCArgumentErrorf(f string, args...interface{}) OSCArgumentError {
	return OSCArgumentError(fmt.Sprintf(f, args...))
}
func (e OSCArgumentError) Error() string {
	return string(e)
}
