package bloom

import (
	"hash"
	"log"
	"math"

	"github.com/mateuszdyminski/bloom-filter/bitset"
	"github.com/spaolacci/murmur3"
)

// A BloomFilter is a representation of a set of _n_ items, where the main
// requirement is to make membership queries; _i.e._, whether an item is a
// member of a set.
type BloomFilter struct {
	length    uint
	hashCount uint
	b         *bitset.BitSet
	h         Hash128
}

type Hash128 interface {
	hash.Hash
	Sum128() (uint64, uint64)
}

// New creates a new Bloom filter with _m_ bits and _k_ hashing functions
func New(length uint, hashCount uint) *BloomFilter {
	return &BloomFilter{length, hashCount, bitset.New(length), murmur3.New128()}
}

func (f *BloomFilter) hash(data []byte) (uint64, uint64) {
	f.h.Reset()
	f.h.Write(data)
	return f.h.Sum128()
}

// EstimateParameters estimates requirements for elements length and false positive rate.
func EstimateParameters(length uint, fpRate float64) (estLength uint, hashCount uint) {
	estLength = uint(math.Ceil(-1 * float64(length) * math.Log(fpRate) / math.Pow(math.Log(2), 2)))
	hashCount = uint(math.Ceil(math.Log(2) * float64(estLength) / float64(length)))
	return
}

// NewWithEstimates creates a new Bloom filter for about length items with fp
// false positive rate
func NewWithEstimates(length uint, fpRate float64) *BloomFilter {
	estLength, hashCount := EstimateParameters(length, fpRate)
	log.Printf("Estimates: bitmap length %d, number of hash functions %d \n", estLength, hashCount)
	return New(estLength, hashCount)
}

// indexes returns array...
func (f *BloomFilter) indexes(data []byte) []uint64 {
	base, inc := f.hash(data)
	indexes := make([]uint64, f.hashCount)

	for i := uint(0); i < f.hashCount; i++ {
		indexes[i] = abs(base % uint64(f.length))
		base += inc
	}

	return indexes
}

func abs(index uint64) uint64 {
	negbit := index >> 63
	return (index ^ negbit) - negbit
}

// Add data to the Bloom Filter. Returns the filter (allows chaining)
func (f *BloomFilter) Add(data []byte) *BloomFilter {
	indexes := f.indexes(data)
	for i := uint(0); i < f.hashCount; i++ {
		f.b.Set(uint(indexes[i]))
	}

	return f
}

// AddString to the Bloom Filter. Returns the filter (allows chaining)
func (f *BloomFilter) AddString(data string) *BloomFilter {
	return f.Add([]byte(data))
}

// Test returns true if the data is in the BloomFilter, false otherwise.
// If true, the result might be a false positive. If false, the data
// is definitely not in the set.
func (f *BloomFilter) Test(data []byte) bool {
	indexes := f.indexes(data)
	for i := uint(0); i < f.hashCount; i++ {
		if !f.b.Test(uint(indexes[i])) {
			return false
		}
	}

	return true
}

// TestString returns true if the string is in the BloomFilter, false otherwise.
// If true, the result might be a false positive. If false, the data
// is definitely not in the set.
func (f *BloomFilter) TestString(data string) bool {
	return f.Test([]byte(data))
}

// TestAndAdd is the equivalent to calling Test(data) then Add(data).
// Returns the result of Test.
func (f *BloomFilter) TestAndAdd(data []byte) bool {
	present := true
	indexes := f.indexes(data)
	for i := uint(0); i < f.hashCount; i++ {
		if !f.b.Test(uint(indexes[i])) {
			present = false
		}
		f.b.Set(uint(indexes[i]))
	}

	return present
}

// TestAndAddString is the equivalent to calling Test(string) then Add(string).
// Returns the result of Test.
func (f *BloomFilter) TestAndAddString(data string) bool {
	return f.TestAndAdd([]byte(data))
}
