package utils

import (
	"bytes"
	"io"
	"io/ioutil"
)

type byteReader struct {
	io.Reader
	buf [1]byte
}

func (b *byteReader) ReadByte() (byte, error) {
	_, err := io.ReadFull(b.Reader, b.buf[:])
	return b.buf[0], err
}

type dos2unix struct {
	r    io.ByteReader
	b    bool
	char byte
}

func (d *dos2unix) Read(b []byte) (int, error) {
	var n int
	for len(b) > 0 {
		if d.b {
			b[0] = d.char
			d.b = false
			b = b[1:]
			n++
			continue
		}
		c, err := d.r.ReadByte()
		if err != nil {
			return n, err
		}
		if c == '\r' {
			d.char, err = d.r.ReadByte()
			if err != io.EOF {
				if err != nil {
					return n, err
				}
				if d.char == '\n' {
					c = '\n'
				} else {
					d.b = true
				}
			}
		}
		b[0] = c
		b = b[1:]
		n++
	}
	return n, nil
}

// DOS2Unix wraps a byte reader with a reader that replaces all instances of
// \r\n with \n
func DOS2Unix(r io.Reader) io.Reader {
	br, ok := r.(io.ByteReader)
	if !ok {
		br = &byteReader{Reader: r}
	}
	return &dos2unix{r: br}
}

type unix2dos struct {
	r  io.ByteReader
	lf bool
}

func (u *unix2dos) Read(b []byte) (int, error) {
	var n int
	for len(b) > 0 {
		if u.lf {
			b[0] = '\n'
			u.lf = false
			b = b[1:]
			n++
			continue
		}
		c, err := u.r.ReadByte()
		if err != nil {
			return n, err
		}
		if c == '\n' {
			u.lf = true
			c = '\r'
		}
		b[0] = c
		b = b[1:]
		n++
	}
	return n, nil
}

// Unix2DOS wraps a byte reader with a reader that replaces all instances of \n
// with \r\n
func Unix2DOS(r io.Reader) io.Reader {
	br, ok := r.(io.ByteReader)
	if !ok {
		br = &byteReader{Reader: r}
	}
	return &unix2dos{r: br}
}

func ToUnix(b []byte) ([]byte, error) {
	return ioutil.ReadAll(DOS2Unix(bytes.NewReader(b)))
}
