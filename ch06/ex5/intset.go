package intset

import (
	"bytes"
	"fmt"
)

const bitsize = 32 << (^uint(0) >> 63)

// An IntSet is a set of small non-negative integers.
// Its zero value represents the empty set.
type IntSet struct {
	words []uint
}

// Has reports whether the set contains the non-negative value x.
func (s *IntSet) Has(x int) bool {
	word, bit := x/bitsize, uint(x%bitsize)
	return word < len(s.words) && s.words[word]&(1<<bit) != 0
}

// Add adds the non-negative value x to the set.
func (s *IntSet) Add(x int) {
	word, bit := x/bitsize, uint(x%bitsize)
	for word >= len(s.words) {
		s.words = append(s.words, 0)
	}
	s.words[word] |= 1 << bit
}

// AddAll adds the given elements to the set.
func (s *IntSet) AddAll(xs ...int) {
	for i := range xs {
		s.Add(xs[i])
	}
}

// UnionWith sets s to the union of s and t.
func (s *IntSet) UnionWith(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			// OR logic
			s.words[i] |= tword
		} else {
			s.words = append(s.words, tword)
		}
	}
}

// IntersectWith sets s to the intersect of s and t.
func (s *IntSet) IntersectWith(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			// AND logic
			s.words[i] &= tword
		}
	}
}

// DifferenceWith sets s to the difference of s and t.
func (s *IntSet) DifferenceWith(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			// NAND logic
			s.words[i] &^= tword
		}
	}
}

// SymmetricDifference sets s to the symmetric difference of s and t.
func (s *IntSet) SymmetricDifference(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			// XOR logic
			s.words[i] ^= tword
		} else {
			s.words = append(s.words, tword)
		}
	}
}

// String returns the set as a string of the form "{1 2 3}".
func (s *IntSet) String() string {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < bitsize; j++ {
			if word&(1<<uint(j)) != 0 {
				if buf.Len() > len("{") {
					buf.WriteByte(' ')
				}
				fmt.Fprintf(&buf, "%d", bitsize*i+j)
			}
		}
	}
	buf.WriteByte('}')
	return buf.String()
}

// Len returns the number of elements.
func (s *IntSet) Len() int {
	var n int
	for _, word := range s.words {
		for i := 0; i < bitsize; i++ {
			if word&(1<<uint(i)) != 0 {
				n++
			}
		}
	}
	return n
}

// Remove removes x from the set.
func (s *IntSet) Remove(x int) {
	word, bit := x/bitsize, uint(x%bitsize)
	s.words[word] &= ^(1 << bit)
}

// Clear removes all the elements from the set.
func (s *IntSet) Clear() {
	for i := range s.words {
		s.words[i] &= 0
	}
}

// Copy returns a copy of the set.
func (s *IntSet) Copy() *IntSet {
	cpy := &IntSet{}
	for i := range s.words {
		cpy.words = append(cpy.words, s.words[i])
	}
	return cpy
}

// Elem returns a slice containing the elements of s.
func (s *IntSet) Elem() []int {
	var elems = make([]int, 0)
	for i := range s.words {
		for j := 0; j < bitsize; j++ {
			if val := s.words[i] & (1 << uint(j)); val != 0 {
				elems = append(elems, i*bitsize+j)
			}
		}
	}
	return elems
}
