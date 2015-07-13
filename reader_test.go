package main

import (
	"bytes"
	"testing"
	"testing/iotest"
)

func TestAllZeros(t *testing.T) {
	for _, s := range []int{0, 1, 10, 1e3, 1e6, 10e6} {
		n, err := ReadZeros(bytes.NewReader(make([]byte, s)))
		if err != nil || n != s {
			t.Errorf("expected %d (<nil>), got %d (%v)", s, n, err)
		}
	}
}

func TestNonZero(t *testing.T) {
	for _, s := range []int{0, 1, 10, 1e3, 1e6, 10e6} {
		b := make([]byte, s+10) // add 10 so we'll leave some zeros at the end
		b[s] = 1                // set s-th byte to non-zero
		n, err := ReadZeros(bytes.NewReader(b))
		if err != errNonZeroByte || n != s {
			t.Errorf("expected %d (errNonZeroByte), got %d (%v)", s, n, err)
		}
	}
}

func TestOneByteReader(t *testing.T) {
	s := int(1e6)
	n, err := ReadZeros(iotest.OneByteReader(bytes.NewReader(make([]byte, s))))
	if err != nil || n != s {
		t.Errorf("expected %d (<nil>), got %d (%v)", s, n, err)
	}
}

func TestReadError(t *testing.T) {
	s := int(1e6)
	n, err := ReadZeros(iotest.TimeoutReader(bytes.NewReader(make([]byte, s))))
	// first read should be 512 bytes
	expected := 512
	if err != iotest.ErrTimeout || n != expected {
		t.Errorf("expected %d (timeout), got %d (%v)", expected, n, err)
	}
}
