package bitvector

import (
	"encoding/binary"
	"os"
)

func SaveIndex(path string, ri *RankIndex) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	header := []uint64{
		ri.BV.N,
		ri.SuperblockSize,
		ri.BlockSize,
		ri.SuperblockBits,
		ri.BlockBits,
		ri.TotalOnes,
		uint64(len(ri.BV.Words)),
		uint64(len(ri.Superblocks.Data)),
		ri.Superblocks.Count,
		uint64(len(ri.Blocks.Data)),
		ri.Blocks.Count,
	}

	for _, v := range header {
		if err := binary.Write(f, binary.LittleEndian, v); err != nil {
			return err
		}
	}

	if err := binary.Write(f, binary.LittleEndian, ri.BV.Words); err != nil {
		return err
	}
	if err := binary.Write(f, binary.LittleEndian, ri.Superblocks.Data); err != nil {
		return err
	}
	if err := binary.Write(f, binary.LittleEndian, ri.Blocks.Data); err != nil {
		return err
	}

	return nil
}

func LoadIndex(path string) (*RankIndex, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var (
		n              uint64
		superblockSize uint64
		blockSize      uint64
		superblockBits uint64
		blockBits      uint64
		totalOnes      uint64
		numBVWords     uint64
		numSBWords     uint64
		numSBCount     uint64
		numBKWords     uint64
		numBKCount     uint64
	)

	for _, ptr := range []*uint64{
		&n, &superblockSize, &blockSize, &superblockBits, &blockBits, &totalOnes,
		&numBVWords, &numSBWords, &numSBCount, &numBKWords, &numBKCount,
	} {
		if err := binary.Read(f, binary.LittleEndian, ptr); err != nil {
			return nil, err
		}
	}

	bvWords := make([]uint64, numBVWords)
	if err := binary.Read(f, binary.LittleEndian, bvWords); err != nil {
		return nil, err
	}

	sbData := make([]uint64, numSBWords)
	if err := binary.Read(f, binary.LittleEndian, sbData); err != nil {
		return nil, err
	}

	bkData := make([]uint64, numBKWords)
	if err := binary.Read(f, binary.LittleEndian, bkData); err != nil {
		return nil, err
	}

	return &RankIndex{
		BV: &BitVector{Words: bvWords, N: n},
		Superblocks: &IntVector{
			Data:           sbData,
			BitsPerElement: superblockBits,
			Count:          numSBCount,
		},
		Blocks: &IntVector{
			Data:           bkData,
			BitsPerElement: blockBits,
			Count:          numBKCount,
		},
		SuperblockSize: superblockSize,
		BlockSize:      blockSize,
		SuperblockBits: superblockBits,
		BlockBits:      blockBits,
		TotalOnes:      totalOnes,
	}, nil
}