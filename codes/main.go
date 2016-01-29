package main

import (
	"bytes"
	"fmt"
	"os"
	"unicode/utf8"
)

/**
 * Command line utility to list character codes given an input string.
 *
 * Output format (one "block" for each argument):
 *
 * $ codes "eɘ"
 *
 * "eɘ"
 * Rune | Ascii           | Unicode           | UTF-8
 * -----------------------------------------------------------------
 *   e  | 101, 0x65, 0145 | 101,  0x65, 0145  | 101,   0x65,    0145
 *   ɘ  |                 | 600, 0x258, 01130 |      0xc998, 0144630
 */

func main() {
	buffer := make([]byte, 4)

	for _, arg := range os.Args[1:] {
		fmt.Printf("%#v\n", arg)
		fmt.Println("Rune | Ascii           | Unicode                 | UTF-8")
		for _, r := range arg {
			var asciiDec, asciiHex, asciiOct string
			var uniDec, uniHex, uniOct string
			var utf8Dec, utf8Hex, utf8Oct string

			uniDec = fmt.Sprintf("%d", r)
			uniHex = fmt.Sprintf("%#x", r)
			uniOct = fmt.Sprintf("%#o", r)

			utf8len := utf8.EncodeRune(buffer, r)
			utf8Hex = fmt.Sprintf("%#x", buffer[0:utf8len])
			utf8Oct = bytesToOctal(buffer[0:utf8len])

			if uint(r) < 128 {
				asciiDec = fmt.Sprintf("%d", r)
				asciiHex = fmt.Sprintf("%#x", r)
				asciiOct = fmt.Sprintf("%#o", r)

				utf8Dec = fmt.Sprintf("%d", r)

				fmt.Printf("  %s  | %3s, %4s, %4s | %5s, %5s, %5s | %3s, %6s, %8s\n",
					string(r),
					asciiDec, asciiHex, asciiOct,
					uniDec, uniHex, uniOct,
					utf8Dec, utf8Hex, utf8Oct)
			} else {
				fmt.Printf("  %s  |                 | %5s, %5s, %5s |      %6s, %8s\n",
					string(r),
					uniDec, uniHex, uniOct,
					utf8Hex, utf8Oct)
			}
		}

		fmt.Println("")
	}
}

func bytesToOctal(buffer []byte) string {
	var out bytes.Buffer

	fmt.Fprint(&out, "0")
	for _, b := range buffer {
		fmt.Fprintf(&out, "%o", int(b))
	}

	return out.String()
}
