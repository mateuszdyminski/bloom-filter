package bloom

import (
	"encoding/binary"
	"testing"
)

func Testbasic(t *testing.T) {
	f := New(1000, 4)
	n1 := make([]byte, 4)
	n2 := make([]byte, 4)
	n3 := make([]byte, 4)
	n4 := make([]byte, 4)

	binary.BigEndian.PutUint32(n1, 100)
	binary.BigEndian.PutUint32(n2, 101)
	binary.BigEndian.PutUint32(n3, 102)
	binary.BigEndian.PutUint32(n4, 103)

	f.Add(n1)
	n3a := f.TestAndAdd(n3)
	n1b := f.Test(n1)
	n2b := f.Test(n2)
	n3b := f.Test(n3)
	f.Test(n4)
	if !n1b {
		t.Errorf("%v should be in.", n1)
	}
	if n2b {
		t.Errorf("%v should not be in.", n2)
	}
	if n3a {
		t.Errorf("%v should not be in the first time we look.", n3)
	}
	if !n3b {
		t.Errorf("%v should be in the second time we look.", n3)
	}
}

func TestString(t *testing.T) {
	f := New(1000, 4)
	n := make([]byte, 4)

	binary.BigEndian.PutUint32(n, 10232134)

	f.Add(n)
	if !f.Test(n) {
		t.Errorf("%v should be in.", n)
	}
}

func TestHeavyLoad(t *testing.T) {
	f := NewWithEstimates(10000000, 0.01)

	for i := uint32(0); i < 10000000; i++ {
		n := make([]byte, 4)
		binary.BigEndian.PutUint32(n, i)
		f.Add(n)
		if !f.Test(n) {
			t.Errorf("%v should be in.", n)
		}
	}
}
