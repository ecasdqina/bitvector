package bitvector

import (
	"log"
	"math/rand"
	"testing"
)

const (
	maxSize  = 256
	attempts = 3
)

func TestBitVector(t *testing.T) {
	fails := 0

	for size := 0; size < maxSize; size++ {
		for attempt := 0; attempt < attempts; attempt++ {
			bs, bv := random(size)
			if got, want := bv.Len(), len(bs); got != want {
				log.Fatalf("%q.Len() => got %v, want %v", bs, got, want)
			}
			var counts [2]int
			for i := 0; i <= len(bs) && fails < 30; i++ {
				if got, want := bv.Rank0(i), counts[0]; got != want {
					t.Errorf("%q.Rank0(%v) => got %v, want %v", bs, i, got, want)
					fails++
				}
				if got, want := bv.Rank1(i), counts[1]; got != want {
					t.Errorf("%q.Rank1(%v) => got %v, want %v", bs, i, got, want)
					fails++
				}
				if i != len(bs) {
					if got, want := bv.Get(i), zeroOne(bs[i]); got != want {
						t.Errorf("%q.Get(%v) => got %v, want %v", bs, i, got, want)
						fails++
					}
					if bs[i] == '0' {
						if got := bv.Select0(counts[0]); got != i {
							t.Errorf("%q.Select0(%v) => got %v, want %v", bs, counts[0], got, i)
							fails++
						}
					} else {
						if got := bv.Select1(counts[1]); got != i {
							t.Errorf("%q.Select1(%v) => got %v, want %v", bs, counts[0], got, i)
							fails++
						}
					}
					counts[zeroOne(bs[i])]++
				}
			}
		}
	}
}

func zeroOne(c byte) int {
	if c == '0' {
		return 0
	}
	return 1
}

func random(size int) (string, *BitVector) {
	var bs []byte
	b := NewBuilder(size)
	for i := 0; i < size; i++ {
		if rand.Intn(2) == 1 {
			bs = append(bs, '1')
			b.Set(i)
		} else {
			bs = append(bs, '0')
		}
	}
	return string(bs), b.Build()
}
