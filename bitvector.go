/*
Package bitvector provides a really simple bit vector implementation that supports only Get, Select
and Rank. Select is implemented by binary search on Rank, thus slow.

Please take a look at https://github.com/hideo55/go-sbvector if you are seeking for more advanced
usages. Also, if you don't need Rank or Select, https://code.google.com/p/go-bit/ is a very good
fit.
*/
package bitvector

import (
	"fmt"

	"github.com/hideo55/go-popcount"
)

// BitVector is a bit vector that supports Get, Rank and Select operation.
type BitVector struct {
	size int
	v    []uint64
	rank []int
}

// Get returns 1 or 0, the value of the i-th bit in the bit vector.
func (b *BitVector) Get(i int) int {
	return int((b.v[i/64] >> uint(i%64)) & 1)
}

// Len returns the size of the bit vector.
func (b *BitVector) Len() int {
	return b.size
}

// Rank1 returns the count of 1s before the i-th bit. It does not count i-th bit itself.
func (b *BitVector) Rank1(i int) int {
	const mask = uint64(0xffffffffffffffff)
	offset := uint(i % 64)
	return int(b.rank[i/64]) + popcnt(b.v[i/64]&^(mask<<offset))
}

// Rank0 returns the count of 0s before the i-th bit. It does not count i-th bit itself.
func (b *BitVector) Rank0(i int) int {
	return i - b.Rank1(i)
}

// Select1 is the inverse of Rank1, i.e. it returns the index of the r-th '1'.
// It is illegal to call this with r > bv.Rank1(bv.Len()).
func (b *BitVector) Select1(r int) int {
	return b.binarySearch(0, b.size+1, r, 1)
}

// Select0 is the inverse of Rank0, i.e. it returns the index of the r-th '0'.
// It is illegal to call this with r > bv.Rank0(bv.Len()).
func (b *BitVector) Select0(r int) int {
	return b.binarySearch(0, b.size+1, r, 0)
}

func (b *BitVector) binarySearch(l, h, r, bit int) int {
	for {
		if l == h {
			panic(fmt.Sprintf("bitvector: no such bit with Rank%v = %v", bit, r))
		}
		if l+1 == h {
			return l
		}

		m := (l + h) / 2
		var pivot int
		if bit == 0 {
			pivot = b.Rank0(m)
		} else {
			pivot = b.Rank1(m)
		}

		if pivot > r {
			h = m
		} else {
			l = m
		}
	}
}

// Builder is a builder of BitVector.
type Builder struct {
	size int
	v    []uint64
}

// NewBuilder makes a new builder of BitVector of the specified size.
func NewBuilder(size int) *Builder {
	bufSize := size/64 + 1

	b := &Builder{
		size: size,
		v:    make([]uint64, bufSize),
	}
	return b
}

// Set sets the i-th bit to 1.
func (b *Builder) Set(i int) {
	b.v[i/64] |= uint64(1) << uint(i%64)
}

// Get returns 1 or 0, the value of the i-th bit in the bit vector.
func (b *Builder) Get(i int) int {
	return int((b.v[i/64] >> uint(i%64)) & 1)
}

// Clear sets the i-th bit to 0. Note that all bits in the BitVector Builder is set to 0 initially,
// so there's no need to call Clear to cleanup all bits.
func (b *Builder) Clear(i int) {
	b.v[i/64] &^= uint64(1) << uint(i%64)
}

// Len returns the size of the bit vector.
func (b *Builder) Len() int {
	return b.size
}

// Build builds a BitVector from the builder.
func (b *Builder) Build() *BitVector {
	rank := make([]int, len(b.v))
	count := 0
	for i, x := range b.v {
		rank[i] = count
		count += popcnt(x)
	}
	return &BitVector{
		size: b.size,
		v:    b.v,
		rank: rank,
	}
}

func popcnt(x uint64) int {
	return int(popcount.Count(x))
}
