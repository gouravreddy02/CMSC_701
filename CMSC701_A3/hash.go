package main

import (
	"encoding/binary"
	"hash/fnv"
)

// The k hash functions of the Bloom filter are realized by varying `seed`.
func hashWithSeed(seed uint64, key []byte, m uint64) uint64 {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], seed)
	h := fnv.New64a()
	h.Write(buf[:])
	h.Write(key)
	return h.Sum64() % m
}
