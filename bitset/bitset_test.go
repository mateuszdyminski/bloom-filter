package bitset

import "testing"

func TestLongsNeeded(t *testing.T) {
	if 4 != longsNeeded(194) {
		t.Errorf("Wrong number of needed uint64! %d; wanted %d", longsNeeded(194), 4)
	}

	if 3 != longsNeeded(190) {
		t.Errorf("Wrong number of needed uint64! %d; wanted %d", longsNeeded(190), 3)
	}

	if 4000 != longsNeeded(256000) {
		t.Errorf("Wrong number of needed uint64! %d; wanted %d", longsNeeded(256000), 4000)
	}

	if 288230376151711743 != longsNeeded(18446744073709551609) {
		t.Errorf("Wrong number of needed uint64! %d; wanted %d", longsNeeded(18446744073709551609), 288230376151711743)
	}
}

func TestSetTest(t *testing.T) {
	// given
	bs := New(1024)

	// when
	bs.Set(1011).Set(111)

	// then
	// should be ok
	if !bs.Test(111) {
		t.Errorf("Test failed! index %d should be set to 1", 111)
	}

	// should be ok
	if !bs.Test(1011) {
		t.Errorf("Test failed! index %d should be set to 1", 1011)
	}

	// test should return false
	if bs.Test(11) {
		t.Errorf("Test failed! index %d should be set to 0", 11)
	}

	// test should return false due to index > bitset length
	if bs.Test(11111) {
		t.Errorf("Test failed! index %d should be set to 0", 11111)
	}
}

func TestSetExceededIndex(t *testing.T) {
	// given
	defer func() {
		// recover from panic which should appear.
		err := recover()
		if err == nil {
			t.Errorf("Panic should be rised!")
		}
	}()

	bs := New(1024)

	// when.. then panic
	bs.Set(1025)
}

func TestClear(t *testing.T) {
	// given
	bs := New(1024)

	// when
	bs.Set(125)
	bs.Clear(125)

	// then
	if bs.Test(125) {
		t.Errorf("Test failed! index %d should be set to 0", 125)
	}
}

func TestFlip(t *testing.T) {
	// given
	bs := New(1024)

	// when
	if bs.Test(125) {
		t.Errorf("Test failed! index %d should be set to 0", 125)
	}
	bs.Flip(125)

	// then
	if !bs.Test(125) {
		t.Errorf("Test failed! index %d should be set to 0", 125)
	}
}
