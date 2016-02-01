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
