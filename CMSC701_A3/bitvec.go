package main

type BitVector struct {
	words []uint64
	n     uint64 
}

func NewBitVector(n uint64) *BitVector {
	return &BitVector{
		words: make([]uint64, (n+63)>>6),
		n:     n,
	}
}

func (b *BitVector) Set(i uint64) {
	b.words[i>>6] |= 1 << (i & 63)
}

func (b *BitVector) Get(i uint64) bool {
	return b.words[i>>6]&(1<<(i&63)) != 0
}

func (b *BitVector) Len() uint64 { return b.n }

// Words exposes the underlying storage for serialization.
func (b *BitVector) Words() []uint64 { return b.words }

// NewBitVectorFromWords reconstructs a BitVector from serialized storage.
func NewBitVectorFromWords(words []uint64, n uint64) *BitVector {
	return &BitVector{words: words, n: n}
}
