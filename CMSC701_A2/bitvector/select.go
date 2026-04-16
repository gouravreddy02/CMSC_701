package bitvector

func Select(ri *RankIndex, r uint64) int64 {
	if r > ri.TotalOnes {
		return -1
	}
	// Find smallest i where Rank(i) > r
	lo, hi := uint64(0), ri.BV.N+1
	for lo < hi {
		mid := lo + (hi-lo)/2
		if ri.Rank(mid) <= r {
			lo = mid + 1
		} else {
			hi = mid
		}
	}
	return int64(lo - 1)
}

// Select1Pos returns the actual position of the k-th 1-bit (1-indexed).
func Select1Pos(ri *RankIndex, k uint64) int64 {
	if k == 0 || k > ri.TotalOnes {
		return -1
	}
	// Binary search for smallest i where Rank(i) >= k
	lo, hi := uint64(0), ri.BV.N+1
	for lo < hi {
		mid := lo + (hi-lo)/2
		if ri.Rank(mid) < k {
			lo = mid + 1
		} else {
			hi = mid
		}
	}
	return int64(lo - 1)
}
