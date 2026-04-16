package bitvector

import (
	"encoding/binary"
	"os"
)


type BitVector struct {
	Words []uint64 // packed 64-bit words holding the bitvector
	N     uint64   // length of the bitvector in bits
}

// AccessBit returns the value of bit at position i (0 or 1).
//   word = Words[i / 64]
//   bit  = (word >> (i % 64)) & 1
func (bv *BitVector) AccessBit(i uint64) uint64 {
	return (bv.Words[i/64] >> (i % 64)) & 1
}

// SetBit sets bit i to 1.
func (bv *BitVector) SetBit(i uint64) {
	bv.Words[i/64] |= 1 << (i % 64)
}

// NewBitVector creates a zero-initialized BitVector of n bits.
func NewBitVector(n uint64) *BitVector {
	numWords := (n + 63) / 64
	return &BitVector{Words: make([]uint64, numWords), N: n}
}

func ReadBitvector(path string) (*BitVector, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// First 4 bytes: bitvector length in bits (little-endian uint32).
	n := uint64(binary.LittleEndian.Uint32(data[0:4]))
	bvBytes := data[4:]

	// Repack bytes into uint64 words.
	numWords := (n + 63) / 64
	words := make([]uint64, numWords)

	// Process full 8-byte chunks: each group of 8 consecutive bytes
	fullChunks := uint64(len(bvBytes)) / 8
	for i := uint64(0); i < fullChunks; i++ {
		words[i] = binary.LittleEndian.Uint64(bvBytes[i*8 : i*8+8])
	}

	// Handle remaining bytes (fewer than 8) for the last partial word.
	remaining := uint64(len(bvBytes)) - fullChunks*8
	if remaining > 0 {
		var lastWord uint64
		for j := uint64(0); j < remaining; j++ {
			lastWord |= uint64(bvBytes[fullChunks*8+j]) << (j * 8)
		}
		words[fullChunks] = lastWord
	}

	return &BitVector{Words: words, N: n}, nil
}
