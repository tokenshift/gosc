package gosc

import (
	"bytes"
	. "testing"
)

func TestWriteMessageEmpty(t *T) {
	var out bytes.Buffer
	var err error
	var n int

	n, err = WriteMessage(&out, OSCAddressPattern("/something"))
	expectNil(t, err)
	expectSame(t, 16, n)
	expectSame(t,
		[]byte{47,115,111,109,101,116,104,105,110,103,0,0, // "/something" (address)
		       44,0,0,0},                                  // "," (empty type string)
		out.Bytes())
}

func TestWriteMessageSingleString(t *T) {
	var out bytes.Buffer
	var err error
	var n int

	n, err = WriteMessage(&out, OSCAddressPattern("/something"), OSCString("foo"))
	expectNil(t, err)
	expectSame(t, 20, n)
	expectSame(t,
		[]byte{47,115,111,109,101,116,104,105,110,103,0,0, // "/something" (address)
		       44,115,0,0,                                 // ",s" (type string)
		       102,111,111,0},                             // "foo"
		out.Bytes())
}

func TestWriteMessageMultipleStrings(t *T) {
	var out bytes.Buffer
	var err error
	var n int

	n, err = WriteMessage(&out, OSCAddressPattern("/something"),
		OSCString("foo"), OSCString("bar"), OSCString("fizz"), OSCString("buzz"))
	expectNil(t, err)
	expectSame(t, 44, n)
	expectSame(t,
		[]byte{47,115,111,109,101,116,104,105,110,103,0,0, // "/something" (address)
		       44,115,115,115,115,0,0,0,                   // ",ssss" (type string)
		       102,111,111,0,                              // "foo"
		       98,97,114,0,                                // "bar"
		       102,105,122,122,0,0,0,0,                    // "fizz"
		       98,117,122,122,0,0,0,0},                    // "buzz"
		out.Bytes())
}

func TestWriteMessageStandardTypes(t *T) {
	var out bytes.Buffer
	var err error
	var n int

	n, err = WriteMessage(&out,
		OSCAddressPattern("/send/this/here"),
		OSCString("foo"),
		OSCInt32(1337),
		OSCFloat32(13.37),
		OSCString("bar"),
		OSCBlob([]byte{1,2,3,4,5}))

	expectNil(t, err)
	expectSame(t, 16+8+4+4+4+4+12, n)
	expectSame(t,
		[]byte{0x2f,0x73,0x65,0x6e,0x64,0x2f,0x74,0x68,0x69,0x73,0x2f,0x68,0x65,0x72,0x65,0x00,
		       0x2c,0x73,0x69,0x66,0x73,0x62,0x00,0x00, // ",sifsb"
		       0x66,0x6f,0x6f,0x00,
		       0x00,0x00,0x05,0x39,
		       0x41,0x55,0xeb,0x85,
		       0x62,0x61,0x72,0x00,
		       0x00,0x00,0x00,0x05,0x01,0x02,0x03,0x04,0x05,0x00,0x00,0x00},
		out.Bytes())
}

func TestReadMessageEmpty(t *T) {
	input := bytes.NewReader([]byte{
		47,115,111,109,101,116,104,105,110,103,0,0, // "/something" (address)
		44,0,0,0})                  // "," (empty type string)

	address, args, err := ReadMessage(input)
	expectNil(t, err)
	expectSame(t, OSCAddressPattern("/something"), address)
	expectSame(t, 0, len(args))
}

func TestReadMessageMultipleStrings(t *T) {
	input := bytes.NewReader([]byte{
		47,115,111,109,101,116,104,105,110,103,0,0, // "/something" (address)
		44,115,115,115,115,0,0,0,                   // ",ssss" (type string)
		102,111,111,0,                              // "foo"
		98,97,114,0,                                // "bar"
		102,105,122,122,0,0,0,0,                    // "fizz"
		98,117,122,122,0,0,0,0})                    // "buzz"

	address, args, err := ReadMessage(input)
	expectNil(t, err)
	expectSame(t, OSCAddressPattern("/something"), address)
	expectSame(t, []OSCArg{
		OSCString("foo"),
		OSCString("bar"),
		OSCString("fizz"),
		OSCString("buzz"),
	}, args)
}

func TestReadMessageStandardTypes(t *T) {
	input := bytes.NewReader([]byte{
		0x2f,0x73,0x65,0x6e,0x64,0x2f,0x74,0x68,0x69,0x73,0x2f,0x68,0x65,0x72,0x65,0x00,
		0x2c,0x73,0x69,0x66,0x73,0x62,0x00,0x00, // ",sifsb"
		0x66,0x6f,0x6f,0x00,
		0x00,0x00,0x05,0x39,
		0x41,0x55,0xeb,0x85,
		0x62,0x61,0x72,0x00,
		0x00,0x00,0x00,0x05,0x01,0x02,0x03,0x04,0x05,0x00,0x00,0x00})

	address, args, err := ReadMessage(input)
	expectNil(t, err)
	expectSame(t, OSCAddressPattern("/send/this/here"), address)
	expectSame(t, 5, len(args))
	expectSame(t, []OSCArg{
		OSCString("foo"),
		OSCInt32(1337),
		OSCFloat32(13.37),
		OSCString("bar"),
		OSCBlob([]byte{1,2,3,4,5}),
	}, args)
}
