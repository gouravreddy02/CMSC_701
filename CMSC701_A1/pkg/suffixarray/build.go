package suffixarray

import "sort"

func BuildSuffixArray(text string) []int {
	n := len(text)
	sa := make([]int, n)
	for i := 0; i < n; i++ {
		sa[i] = i
	}

	sort.Slice(sa, func(a, b int) bool {
		i, j := sa[a], sa[b]
		for i < n && j < n {
			if text[i] != text[j] {
				return text[i] < text[j]
			}
			i++
			j++
		}
		return (n - sa[a]) < (n - sa[b])
	})

	return sa
}