package gosc

/**
 * Tests for deserialization (reading) of individual OSC arguments.
 */

import (
	"bytes"
	"io"
	"math"
	. "testing"
)

func TestReadOSCInt32(t *T) {
	var input io.Reader
	var i OSCInt32
	var err error

	input = bytes.NewReader([]byte{0,0,0,0})
	i, err = ReadOSCInt32(input)
	expectNil(t, err)
	expectSame(t, i, OSCInt32(0))

	input = bytes.NewReader([]byte{0,0,0,0x2a})
	i, err = ReadOSCInt32(input)
	expectNil(t, err)
	expectSame(t, i, OSCInt32(42))
}

func TestReadOSCFloat32(t *T) {
	var input io.Reader
	var f OSCFloat32
	var err error

	input = bytes.NewReader([]byte{0x3a,0x83,0x12,0x6f})
	f, err = ReadOSCFloat32(input)
	expectNil(t, err)
	expectSame(t, f, OSCFloat32(0.001))

	input = bytes.NewReader([]byte{0x7f,0x7f,0xff,0xff})
	f, err = ReadOSCFloat32(input)
	expectNil(t, err)
	expectSame(t, f, OSCFloat32(math.MaxFloat32))
}

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

func TestReadOSCBlob(t *T) {
	var input io.Reader
	var b OSCBlob
	var err error

	input = bytes.NewReader([]byte{0,0,0,5,1,2,3,4,5,0,0,0})
	b, err = ReadOSCBlob(input)
	expectNil(t, err)
	expectSame(t, b, OSCBlob([]byte{1,2,3,4,5}))
}
