package gosc

/**
 * Tests for deserialization (reading) of individual OSC arguments.
 */

import (
	"bytes"
	"io"
	. "testing"
)

func TestReadOSCString(t *T) {
	var input io.Reader
	var s OSCString
	var err error

	input = bytes.NewReader([]byte{116,101,115,116,105,110,103,0}) // "testing"
	s, err = ReadOSCString(input)
	expectNil(t, err)
	expectSame(t, OSCString("testing"), s)

	input = bytes.NewReader([]byte{0,0,0,0}) // ""
	s, err = ReadOSCString(input)
	expectNil(t, err)
	expectSame(t, OSCString(""), s)

	input = bytes.NewReader([]byte{49,0,0,0}) // "1"
	s, err = ReadOSCString(input)
	expectNil(t, err)
	expectSame(t, OSCString("1"), s)

	input = bytes.NewReader([]byte{49,50,0,0}) // "12"
	s, err = ReadOSCString(input)
	expectNil(t, err)
	expectSame(t, OSCString("12"), s)

	input = bytes.NewReader([]byte{49,50,51,0}) // "123"
	s, err = ReadOSCString(input)
	expectNil(t, err)
	expectSame(t, OSCString("123"), s)

	input = bytes.NewReader([]byte{49,50,51,52,0,0,0,0}) // "1234"
	s, err = ReadOSCString(input)
	expectNil(t, err)
	expectSame(t, OSCString("1234"), s)
}

func TestReadMultipleOSCStrings(t *T) {
	var input io.Reader
	var s OSCString
	var err error

	input = bytes.NewReader([]byte{
		116,101,115,116,105,110,103,0, // "testing"
		49,50,51,52,0,0,0,0,           // "1234"
	})

	s, err = ReadOSCString(input)
	expectNil(t, err)
	expectSame(t, OSCString("testing"), s)

	s, err = ReadOSCString(input)
	expectNil(t, err)
	expectSame(t, OSCString("1234"), s)
}
