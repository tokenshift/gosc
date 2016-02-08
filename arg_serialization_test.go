package gosc

/**
 * Tests for serialization (writing) of individual OSC arguments.
 */

import (
	"bytes"
	"math"
	. "testing"
)

func TestWriteOSCInt32(t *T) {
	var out bytes.Buffer
	var i OSCInt32
	var err error
	var n int

	i = 0
	n, err = i.WriteTo(&out)
	expectNil(t, err)
	expectSame(t, 4, n)
	expectSame(t, []byte{0,0,0,0}, out.Bytes())
	out.Reset()

	i = 42
	n, err = i.WriteTo(&out)
	expectNil(t, err)
	expectSame(t, 4, n)
	expectSame(t, []byte{0,0,0,0x2a}, out.Bytes())
	out.Reset()

	i = -125135
	n, err = i.WriteTo(&out)
	expectNil(t, err)
	expectSame(t, 4, n)
	expectSame(t, []byte{0xff,0xfe,0x17,0x31}, out.Bytes())
	out.Reset()

	i = math.MaxInt32
	n, err = i.WriteTo(&out)
	expectNil(t, err)
	expectSame(t, 4, n)
	expectSame(t, []byte{0x7f,255,255,255}, out.Bytes())
	out.Reset()

	i = math.MinInt32
	n, err = i.WriteTo(&out)
	expectNil(t, err)
	expectSame(t, 4, n)
	expectSame(t, []byte{0x80,0,0,0}, out.Bytes())
	out.Reset()
}

func TestWriteOSCFloat32(t *T) {
	// http://www.h-schmidt.net/FloatConverter/IEEE754.html

	var out bytes.Buffer
	var f OSCFloat32
	var err error
	var n int

	f = 0.0
	n, err = f.WriteTo(&out)
	expectNil(t, err)
	expectSame(t, 4, n)
	expectSame(t, []byte{0,0,0,0}, out.Bytes())
	out.Reset()

	f = 0.001
	n, err = f.WriteTo(&out)
	expectNil(t, err)
	expectSame(t, 4, n)
	expectSame(t, []byte{0x3a,0x83,0x12,0x6f}, out.Bytes())
	out.Reset()

	f = -125.135
	n, err = f.WriteTo(&out)
	expectNil(t, err)
	expectSame(t, 4, n)
	expectSame(t, []byte{0xc2,0xfa,0x45,0x1f}, out.Bytes())
	out.Reset()

	f = math.MaxFloat32
	n, err = f.WriteTo(&out)
	expectNil(t, err)
	expectSame(t, 4, n)
	expectSame(t, []byte{0x7f,0x7f,0xff,0xff}, out.Bytes())
	out.Reset()

	f = math.SmallestNonzeroFloat32
	n, err = f.WriteTo(&out)
	expectNil(t, err)
	expectSame(t, 4, n)
	expectSame(t, []byte{0x00,0x00,0x00,0x01}, out.Bytes())
	out.Reset()
}

func TestWriteOSCString(t *T) {
	var out bytes.Buffer
	var s OSCString
	var err error
	var ok bool
	var n int

	s = OSCString("testing")
	expectNil(t, s.Valid())
	expectSame(t, OSC_TYPE_STRING, s.Tag())
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

func TestWriteOSCBlob(t *T) {
	var out bytes.Buffer
	var b OSCBlob
	var err error
	var n int

	b = []byte{}
	n, err = b.WriteTo(&out)
	expectNil(t, err)
	expectSame(t, 4, n)
	expectSame(t, []byte{0,0,0,0}, out.Bytes())
	out.Reset()

	b = []byte{1}
	n, err = b.WriteTo(&out)
	expectNil(t, err)
	expectSame(t, 8, n)
	expectSame(t, []byte{0,0,0,1,1,0,0,0}, out.Bytes())
	out.Reset()

	b = []byte{1,2}
	n, err = b.WriteTo(&out)
	expectNil(t, err)
	expectSame(t, 8, n)
	expectSame(t, []byte{0,0,0,2,1,2,0,0}, out.Bytes())
	out.Reset()

	b = []byte{1,2,3}
	n, err = b.WriteTo(&out)
	expectNil(t, err)
	expectSame(t, 8, n)
	expectSame(t, []byte{0,0,0,3,1,2,3,0}, out.Bytes())
	out.Reset()

	b = []byte{1,2,3,4}
	n, err = b.WriteTo(&out)
	expectNil(t, err)
	expectSame(t, 8, n)
	expectSame(t, []byte{0,0,0,4,1,2,3,4}, out.Bytes())
	out.Reset()

	b = []byte{1,2,3,4,5}
	n, err = b.WriteTo(&out)
	expectNil(t, err)
	expectSame(t, 12, n)
	expectSame(t, []byte{0,0,0,5,1,2,3,4,5,0,0,0}, out.Bytes())
	out.Reset()
}
