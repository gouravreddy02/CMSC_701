package suffixarray

func BuildPrefixTable(text string, sa []int, k int) map[string][2]int {
	table := make(map[string][2]int)
	n := len(text)
	prevPrefix := ""

	for i, offset := range sa {
		// Skip suffixes shorter than k
		if offset+k > n {
			continue
		}
		prefix := text[offset : offset+k]

		if prefix != prevPrefix {
			// Close the previous prefix range
			if prevPrefix != "" {
				entry := table[prevPrefix]
				entry[1] = i
				table[prevPrefix] = entry
			}
			// Open a new prefix range
			table[prefix] = [2]int{i, 0}
			prevPrefix = prefix
		}
	}
	// Close the last prefix range
	if prevPrefix != "" {
		entry := table[prevPrefix]
		entry[1] = len(sa)
		table[prevPrefix] = entry
	}

	return table
}
