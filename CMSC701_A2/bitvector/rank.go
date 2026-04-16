package bitvector

import (
	"math"
	"math/bits"
)

//	Level 1 — Superblocks:
//	  Divide the bitvector into chunks of size s = (log₂ N)² bits.
//	Level 2 — Blocks:
//	  Divide each superblock into chunks of size b = floor((log₂ N) / 2) bits.
type RankIndex struct {
	BV *BitVector 
	Superblocks *IntVector 
	Blocks      *IntVector 

	SuperblockSize uint64 
	BlockSize      uint64 

	SuperblockBits uint64 
	BlockBits      uint64 

	TotalOnes uint64 
}

func BuildRankIndex(bv *BitVector) *RankIndex {
	n := bv.N
	logN := uint64(math.Ceil(math.Log2(float64(n))))
	if logN < 2 {
		logN = 2
	}

	blockSize := logN / 2  
	if blockSize == 0 {
		blockSize = 1
	}

	superblockSize := logN * logN 
	if superblockSize%blockSize != 0 {
		superblockSize = ((superblockSize / blockSize) + 1) * blockSize
	}

	numSuperblocks := (n + superblockSize - 1) / superblockSize + 1
	numBlocks := (n + blockSize - 1) / blockSize + 1
	superblockBits := bitsNeeded(n)
	blockBits := bitsNeeded(superblockSize)

	superblocks := NewIntVector(numSuperblocks, superblockBits)
	blocks := NewIntVector(numBlocks, blockBits)

	var cumulativeRank uint64
	var totalOnes uint64

	for i := uint64(0); i < n; i++ {
		if i%superblockSize == 0 {
			sbIdx := i / superblockSize
			superblocks.Set(sbIdx, cumulativeRank)
		}

		if i%blockSize == 0 {
			bIdx := i / blockSize
			sbStart := (i / superblockSize) * superblockSize
			relativeRank := cumulativeRank - superblocks.Get(i/superblockSize)
			_ = sbStart
			blocks.Set(bIdx, relativeRank)
		}

		bit := bv.AccessBit(i)
		cumulativeRank += bit
	}
	totalOnes = cumulativeRank

	lastSBIdx := n / superblockSize
	if n%superblockSize == 0 {
		superblocks.Set(lastSBIdx, cumulativeRank)
	}
	lastBIdx := n / blockSize
	if n%blockSize == 0 {
		blocks.Set(lastBIdx, cumulativeRank-superblocks.Get(n/superblockSize))
	}

	return &RankIndex{
		BV:             bv,
		Superblocks:    superblocks,
		Blocks:         blocks,
		SuperblockSize: superblockSize,
		BlockSize:      blockSize,
		SuperblockBits: superblockBits,
		BlockBits:      blockBits,
		TotalOnes:      totalOnes,
	}
}

func (ri *RankIndex) Rank(i uint64) uint64 {
	if i == 0 {
		return 0
	}
	if i >= ri.BV.N {
		return ri.TotalOnes
	}

	part1 := ri.Superblocks.Get(i / ri.SuperblockSize)
	part2 := ri.Blocks.Get(i / ri.BlockSize)

	blockStart := (i / ri.BlockSize) * ri.BlockSize
	part3 := uint64(0)

	pos := blockStart
	for pos+64 <= i {
		wordIdx := pos / 64
		if pos%64 == 0 {
			part3 += uint64(bits.OnesCount64(ri.BV.Words[wordIdx]))
			pos += 64
		} else {
			bit := ri.BV.AccessBit(pos)
			part3 += bit
			pos++
		}
	}

	for pos < i {
		part3 += ri.BV.AccessBit(pos)
		pos++
	}

	return part1 + part2 + part3
}

func bitsNeeded(v uint64) uint64 {
	if v == 0 {
		return 1
	}
	return uint64(bits.Len64(v))
}