package main

import (
	"errors"
	"io"
)

// errNonZeroByte will be returned when something non-zero was read
var errNonZeroByte = errors.New("non-zero byte read")

// ZeroReader reads from R until a non-zero byte is read
type ZeroReader struct {
	R io.Reader
}

// Read reads up to len(p) bytes. It returnes how many leading zeros was
// read. If a non-zero byte occures, errNonZeroByte is returned.
func (zr ZeroReader) Read(p []byte) (zeros int, err error) {
	for len(p) > 0 && err == nil {
		var n int
		n, err = zr.R.Read(p)
		for i := range p[:n] {
			if p[i] != 0 {
				zeros += i
				err = errNonZeroByte
				return
			}
		}
		p = p[n:]
		zeros += n
	}
	return
}

// ReadZeros will read from r until incoming bytes are zero.
// If non-zero byte occures errNonZeroByte will be returned.
// It returnes number of zeros read from the source and optional error;
// io.EOF is never returned, so if err is nil, it means source doesn't contain
// any non-zero bytes.
func ReadZeros(r io.Reader) (zeros int, err error) {
	buf := make([]byte, 512*1024) // read buffer 512kb
	// Create zero reader, ensure first read is small as most sources will have
	// non-zero byte somwhere at the beginning
	r = io.MultiReader(io.LimitReader(ZeroReader{r}, 512), ZeroReader{r})
	for err == nil {
		var n int
		n, err = r.Read(buf)
		zeros += n
	}
	if err == io.EOF {
		err = nil
	}
	return
}
