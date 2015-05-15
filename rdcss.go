package rdcss

import (
	"sync/atomic"
	"unsafe"
)

// rdcssDescriptor is an intermediate struct which communicates the intent to
// replace the value at address a2 while ensuring the values at a1 and a2
// haven't changed before committing to the new value.
type rdcssDescriptor struct {
	a1        *unsafe.Pointer
	o1        unsafe.Pointer
	a2        *unsafe.Pointer
	o2        unsafe.Pointer
	n2        unsafe.Pointer
	committed bool
}

// RDCSS performs an RDCSS double-compare-single-swap operation. The value n2
// is swapped into address a2 iff the value at a1 is o1 and the value at a2 is
// o2. Returns whether or not the CAS succeeded.
func RDCSS(a1 *unsafe.Pointer, o1 unsafe.Pointer, a2 *unsafe.Pointer,
	o2, n2 unsafe.Pointer) bool {

	d := &rdcssDescriptor{
		a1: a1,
		o1: o1,
		a2: a2,
		o2: o2,
		n2: n2,
	}

	if atomic.CompareAndSwapPointer(
		(*unsafe.Pointer)(unsafe.Pointer(a2)), unsafe.Pointer(o2), unsafe.Pointer(d)) {
		complete(d)
		return d.committed
	}

	return false
}

func complete(d *rdcssDescriptor) {
	d = (*rdcssDescriptor)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&d))))
	if *d.a1 == d.o1 && atomic.CompareAndSwapPointer(d.a2, unsafe.Pointer(d), d.n2) {
		d.committed = true
	} else {
		atomic.CompareAndSwapPointer(d.a2, unsafe.Pointer(d), d.o2)
	}
}
