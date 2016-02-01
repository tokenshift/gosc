package gosc

import (
	"io"
)

const OSC_BYTE_ALIGNMENT = 4
const OSC_STRING_BUFFER_SIZE = 1024

// Arguments are all OSC types that can be transitted in a message.
type OSCArg interface {
	// WriteTo writes the argument to an output stream. Returns the number of
	// bytes written and an error, if the argument could not be written.
	// This error will either be a serialization error (the argument could not
	// be turned into a stream of bytes, usually because it was invalid in some
	// way) or an underlying transmission failure returned by the output Writer.
	WriteTo(out io.Writer) (int, error)

	// Tag returns the single-byte type tag that identifies the argument type
	// on the wire.
	Tag() OSCTypeTag

	// Valid checks its value and ensures that it can be serialized correctly.
	// Returns nil on success, or an error otherwise.
	Valid() error
}

type OSCTypeTag byte
const (
	// Standard argument types.
	OSC_TYPE_INT32        = OSCTypeTag('i')
	OSC_TYPE_FLOAT32      = OSCTypeTag('f')
	OSC_TYPE_STRING       = OSCTypeTag('s')
	OSC_TYPE_BLOB         = OSCTypeTag('b')

	// "Extended" argument types (not all clients support any/all of these).
	OSC_ETYPE_INT64       = OSCTypeTag('h')
	OSC_ETYPE_TIMETAG     = OSCTypeTag('t')
	OSC_ETYPE_FLOAT64     = OSCTypeTag('d')
	OSC_ETYPE_STRING_ALT  = OSCTypeTag('S')
	OSC_ETYPE_CHAR        = OSCTypeTag('c')
	OSC_ETYPE_RGBA        = OSCTypeTag('r')
	OSC_ETYPE_MIDI        = OSCTypeTag('m')
	OSC_ETYPE_TRUE        = OSCTypeTag('T')
	OSC_ETYPE_FALSE       = OSCTypeTag('F')
	OSC_ETYPE_NIL         = OSCTypeTag('N')
	OSC_ETYPE_INFINITY    = OSCTypeTag('I')
	OSC_ETYPE_ARRAY_START = OSCTypeTag('[')
	OSC_ETYPE_ARRAY_END   = OSCTypeTag(']')
)

// OSC-strings are more restrictive than go strings, so a []byte would be more
// appropriate; string is used purely for convenience, so that consumers of the
// library don't have to add conversion logic to their string literals. All
// strings are validated before transmission.
type OSCString string

func ReadOSCString(in io.Reader) (OSCString, error) {
	var s []byte

	var buf [1]byte
	var n int
	var err error

	for n, err = in.Read(buf[:]); err == nil && n > 0 && buf[0] != 0; n, err = in.Read(buf[:]) {
		s = append(s, buf[0])
	}

	if err != nil && err != io.EOF {
		return "", OSCReadErrorf("failed to read OSC-string from input: %v", err)
	}

	if n == 0 {
		return "", OSCReadErrorf("did not reach null terminator in OSC-string")
	}

	// Then discard null padding (OSC-strings are supposed to be padded to four
	// byte increments).

	for i := (4 - (len(s) + 1)) % 4; i > 0; i-- {
		n, err = in.Read(buf[:])
		if err != nil {
			return "", OSCReadErrorf("failed to read OSC-string from input: %v", err)
		}
		if n == 0 || buf[0] != 0 {
			//fmt.Println("HERE:", s, n, 
			return "", OSCReadErrorf("OSC-string was not padded properly")
		}
	}

	os := OSCString(s)
	return os, os.Valid()
}

func (s OSCString) Tag() OSCTypeTag {
	return OSC_TYPE_STRING
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

func (s OSCAddressPattern) WriteTo(out io.Writer) (int, error) {
	return OSCString(s).WriteTo(out)
}
