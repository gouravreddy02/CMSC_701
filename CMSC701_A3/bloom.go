package main

import (
	"encoding/binary"
	"errors"
	"io"
	"math"
)

var magic = [4]byte{'B', 'F', 'L', 'T'}

const formatVersion uint8 = 1


type BloomFilter struct {
	N     uint64 
	M     uint64 
	K     uint64 
	Seeds []uint64
	Bits  *BitVector
}

// OptimalParams returns m, k minimizing bit-array size for the target fpr.
func OptimalParams(n uint64, fpr float64) (m, k uint64) {
	if n == 0 {
		return 1, 1
	}
	ln2 := math.Ln2
	mFloat := -float64(n) * math.Log(fpr) / (ln2 * ln2)
	m = uint64(math.Ceil(mFloat))
	if m < 1 {
		m = 1
	}
	kFloat := (float64(m) / float64(n)) * ln2
	k = uint64(math.Ceil(kFloat))
	if k < 1 {
		k = 1
	}
	return m, k
}

// NewBloomFilter allocates an empty filter with the given parameters and seeds.
func NewBloomFilter(n, m, k uint64, seeds []uint64) *BloomFilter {
	return &BloomFilter{
		N:     n,
		M:     m,
		K:     k,
		Seeds: seeds,
		Bits:  NewBitVector(m),
	}
}

func (bf *BloomFilter) Insert(key []byte) {
	for _, s := range bf.Seeds {
		bf.Bits.Set(hashWithSeed(s, key, bf.M))
	}
}

func (bf *BloomFilter) Contains(key []byte) bool {
	for _, s := range bf.Seeds {
		if !bf.Bits.Get(hashWithSeed(s, key, bf.M)) {
			return false
		}
	}
	return true
}

// Write serializes the filter in the format described in bloom.go comments.
func (bf *BloomFilter) Write(w io.Writer) error {
	header := struct {
		Magic   [4]byte
		Version uint8
		_       [3]byte
		N       uint64
		M       uint64
		K       uint64
	}{Magic: magic, Version: formatVersion, N: bf.N, M: bf.M, K: bf.K}
	if err := binary.Write(w, binary.LittleEndian, &header); err != nil {
		return err
	}
	if uint64(len(bf.Seeds)) != bf.K {
		return errors.New("bloom: seed count does not match k")
	}
	if err := binary.Write(w, binary.LittleEndian, bf.Seeds); err != nil {
		return err
	}
	return binary.Write(w, binary.LittleEndian, bf.Bits.Words())
}

// ReadBloomFilter deserializes a filter previously produced by Write.
func ReadBloomFilter(r io.Reader) (*BloomFilter, error) {
	var header struct {
		Magic   [4]byte
		Version uint8
		_       [3]byte
		N       uint64
		M       uint64
		K       uint64
	}
	if err := binary.Read(r, binary.LittleEndian, &header); err != nil {
		return nil, err
	}
	if header.Magic != magic {
		return nil, errors.New("bloom: bad magic (not a bfilt file)")
	}
	if header.Version != formatVersion {
		return nil, errors.New("bloom: unsupported format version")
	}
	seeds := make([]uint64, header.K)
	if err := binary.Read(r, binary.LittleEndian, seeds); err != nil {
		return nil, err
	}
	nWords := (header.M + 63) >> 6
	words := make([]uint64, nWords)
	if err := binary.Read(r, binary.LittleEndian, words); err != nil {
		return nil, err
	}
	return &BloomFilter{
		N:     header.N,
		M:     header.M,
		K:     header.K,
		Seeds: seeds,
		Bits:  NewBitVectorFromWords(words, header.M),
	}, nil
}
