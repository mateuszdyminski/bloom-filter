package bitset

// the uint64Size of a bit set
const uint64Size = uint(64)

// log2Uint64Size is lg(uint64Size)
const log2Uint64Size = uint(6)

// BitSet efficient and fast set of bits.
type BitSet struct {
	length uint
	set    []uint64
}

// New creates a new BitSet with specified length.
func New(length uint) *BitSet {
	return &BitSet{length, make([]uint64, longsNeeded(length))}
}

// Len returns the length of the BitSet in longs(uint64).
func (b *BitSet) Len() uint {
	return b.length
}

// Test whether bit i is set.
func (b *BitSet) Test(i uint) bool {
	if i >= b.length {
		return false
	}
	return b.set[i>>log2Uint64Size]&(1<<(i&(uint64Size-1))) != 0
}

// Set bit i to 1.
func (b *BitSet) Set(i uint) *BitSet {
	b.set[i>>log2Uint64Size] |= 1 << (i & (uint64Size - 1))
	return b
}

// Clear bit i to 0.
func (b *BitSet) Clear(i uint) *BitSet {
	if i >= b.length {
		return b
	}
	b.set[i>>log2Uint64Size] &^= 1 << (i & (uint64Size - 1))
	return b
}

// Flip bit at i.
func (b *BitSet) Flip(i uint) *BitSet {
	if i >= b.length {
		return b.Set(i)
	}
	b.set[i>>log2Uint64Size] ^= 1 << (i & (uint64Size - 1))
	return b
}

// longsNeeded calculates the number of longs(uint64) needed for bits
func longsNeeded(bits uint) int {
	if bits > ((^uint(0)) - uint64Size + 1) {
		return int((^uint(0)) >> log2Uint64Size)
	}
	return int((bits + (uint64Size - 1)) >> log2Uint64Size)
}
