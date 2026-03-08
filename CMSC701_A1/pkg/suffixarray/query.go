package suffixarray

func naiveBinarySearch(text string, sa []int, pattern string, findLower bool, lo, hi int) (int, int) {
	totalCmp := 0
	l := lo - 1
	r := hi

	for r-l > 1 {
		c := (l + r) / 2
		offset := sa[c]
		cmp := 0
		i := 0
		for i < len(pattern) && offset+i < len(text) {
			totalCmp++
			if pattern[i] < text[offset+i] {
				cmp = -1
				break
			} else if pattern[i] > text[offset+i] {
				cmp = 1
				break
			}
			i++
		}
		if cmp == 0 {
			if i == len(pattern) {

				if findLower {
					r = c
				} else {
					l = c
				}
			} else {

				l = c
			}
		} else if cmp < 0 {
			r = c
		} else {
			l = c
		}
	}
	return r, totalCmp
}

func QueryNaive(text string, sa []int, pattern string) (int, int, []int) {
	lb, cmpLB := naiveBinarySearch(text, sa, pattern, true, 0, len(sa))
	ub, cmpUB := naiveBinarySearch(text, sa, pattern, false, 0, len(sa))
	hits := make([]int, ub-lb)
	for i := lb; i < ub; i++ {
		hits[i-lb] = sa[i]
	}
	return cmpLB, cmpUB, hits
}

func simpaccelBinarySearch(text string, sa []int, pattern string, findLower bool, lo, hi int) (int, int) {
	totalCmp := 0
	l := lo - 1
	r := hi
	lcpL, lcpR := 0, 0

	for r-l > 1 {
		c := (l + r) / 2
		offset := sa[c]
		skip := lcpL
		if lcpR < skip {
			skip = lcpR
		}
		i := skip
		matched := skip
		cmp := 0
		for i < len(pattern) && offset+i < len(text) {
			totalCmp++
			if pattern[i] < text[offset+i] {
				cmp = -1
				break
			} else if pattern[i] > text[offset+i] {
				cmp = 1
				break
			}
			i++
			matched++
		}
		if cmp == 0 {
			if i == len(pattern) {
				if findLower {
					r = c
					lcpR = matched
				} else {
					l = c
					lcpL = matched
				}
			} else {
				l = c
				lcpL = matched
			}
		} else if cmp < 0 {
			r = c
			lcpR = matched
		} else {
			l = c
			lcpL = matched
		}
	}
	return r, totalCmp
}

func QuerySimpaccel(text string, sa []int, pattern string) (int, int, []int) {
	lb, cmpLB := simpaccelBinarySearch(text, sa, pattern, true, 0, len(sa))
	ub, cmpUB := simpaccelBinarySearch(text, sa, pattern, false, 0, len(sa))
	hits := make([]int, ub-lb)
	for i := lb; i < ub; i++ {
		hits[i-lb] = sa[i]
	}
	return cmpLB, cmpUB, hits
}

func QueryPrefaccel(text string, sa []int, prefixTable map[string][2]int, k int, pattern string) (int, int, []int) {
	prefix := pattern[:k]
	interval, found := prefixTable[prefix]
	if !found {
		return 0, 0, nil
	}
	lb, _ := naiveBinarySearch(text, sa, pattern, true, interval[0], interval[1])
	ub, _ := naiveBinarySearch(text, sa, pattern, false, interval[0], interval[1])
	hits := make([]int, ub-lb)
	for i := lb; i < ub; i++ {
		hits[i-lb] = sa[i]
	}
	return interval[0], interval[1], hits
}
