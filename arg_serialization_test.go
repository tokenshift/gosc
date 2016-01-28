package gosc

/**
 * Tests for serialization of individual OSC arguments.
 */

import (
	"bytes"
	. "testing"
)

func TestOSTString(t *T) {
	var out bytes.Buffer
	var s OSCString
	var err error
	var ok bool
	var n int

	s = OSCString("testing")
	expectNil(t, s.Valid())
	expectSame(t, OST_TYPE_STRING, s.Tag())
	n, err = s.WriteTo(&out)
	expectNil(t, err)
	expectSame(t, 8, n)
	expectSame(t, []byte{116, 101, 115, 116, 105, 110, 103, 0}, out.Bytes())

	s = OSCString("tɘsting")
	if err, ok = s.Valid().(OSCArgumentError); !ok {
		t.Errorf("expected an OSCArgumentError, got %#v (%T)", err, err)
	}

	out.Reset()
	s = OSCString("")
	expectNil(t, s.Valid())
	n, err = s.WriteTo(&out)
	expectNil(t, err)
	expectSame(t, 4, n)
	expectSame(t, []byte{0,0,0,0}, out.Bytes())

	out.Reset()
	s = OSCString("1")
	expectNil(t, s.Valid())
	n, err = s.WriteTo(&out)
	expectNil(t, err)
	expectSame(t, 4, n)
	expectSame(t, []byte{49,0,0,0}, out.Bytes())

	out.Reset()
	s = OSCString("12")
	expectNil(t, s.Valid())
	n, err = s.WriteTo(&out)
	expectNil(t, err)
	expectSame(t, 4, n)
	expectSame(t, []byte{49,50,0,0}, out.Bytes())

	out.Reset()
	s = OSCString("123")
	expectNil(t, s.Valid())
	n, err = s.WriteTo(&out)
	expectNil(t, err)
	expectSame(t, 4, n)
	expectSame(t, []byte{49,50,51,0}, out.Bytes())

	out.Reset()
	s = OSCString("1234")
	expectNil(t, s.Valid())
	n, err = s.WriteTo(&out)
	expectNil(t, err)
	expectSame(t, 8, n)
	expectSame(t, []byte{49,50,51,52,0,0,0,0}, out.Bytes())
}
