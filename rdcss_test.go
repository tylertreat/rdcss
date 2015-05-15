package rdcss

import (
	"testing"
	"unsafe"
)

type field struct {
	x string
}

type container struct {
	x *field
	y *field
}

func TestRDCSS(t *testing.T) {
	expected := &field{"bar"}
	c := &container{
		x: &field{"foo"},
		y: expected,
	}

	if !RDCSS(
		(*unsafe.Pointer)(unsafe.Pointer(&c.x)),
		unsafe.Pointer(c.x),
		(*unsafe.Pointer)(unsafe.Pointer(&c.y)),
		unsafe.Pointer(c.y),
		unsafe.Pointer(&field{"baz"})) {
		t.Errorf("Expected RDCSS to succeed")
	}

	if c.y.x != "baz" {
		t.Errorf("Expected baz, got %s", c.y.x)
	}
}
