package main

import (
	"io"
	"os"
	"strings"
)

type rot13Reader struct {
	r io.Reader
}

func (rotr *rot13Reader) Read(b []byte) (int, error) {
	x, err := rotr.r.Read(b)
	if err == nil {
		for i := 0; i < x; i++ {
			if b[i] >= 'a' && b[i] <= 'z' {
				b[i] = 'a' + (b[i]-'a'+13)%26
			} else if b[i] >= 'A' && b[i] <= 'Z' {
				b[i] = 'A' + (b[i]-'A'+13)%26
			}
		}
	}
	return x, err
}

func main() {
	s := strings.NewReader("Lbh penpxrq gur pbqr!")
	r := rot13Reader{s}
	io.Copy(os.Stdout, &r)
}
