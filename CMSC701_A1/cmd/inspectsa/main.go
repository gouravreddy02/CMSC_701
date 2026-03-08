package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"cmsc701_a1/pkg/suffixarray"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintf(os.Stderr, "Usage: inspectsa <index> <sample_rate> <output>\n")
		os.Exit(1)
	}

	idxPath := os.Args[1]
	sampleRate, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid sample_rate: %s\n", os.Args[2])
		os.Exit(1)
	}
	outPath := os.Args[3]

	// Load the index
	idx, err := suffixarray.Load(idxPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading index: %v\n", err)
		os.Exit(1)
	}

	text := idx.Text
	sa := idx.SA
	n := len(sa)

	// Compute LCP1 array
	lcp1 := make([]int, n-1)
	for i := 0; i < n-1; i++ {
		a, b := sa[i], sa[i+1]
		l := 0
		for a+l < len(text) && b+l < len(text) && text[a+l] == text[b+l] {
			l++
		}
		lcp1[i] = l
	}

	// Mean
	sum := 0
	for _, v := range lcp1 {
		sum += v
	}
	mean := float64(sum) / float64(len(lcp1))

	// Max
	maxVal := 0
	for _, v := range lcp1 {
		if v > maxVal {
			maxVal = v
		}
	}

	// Median
	sorted := make([]int, len(lcp1))
	copy(sorted, lcp1)
	sort.Ints(sorted)
	var median float64
	mid := len(sorted) / 2
	if len(sorted)%2 == 0 {
		median = float64(sorted[mid-1]+sorted[mid]) / 2.0
	} else {
		median = float64(sorted[mid])
	}

	// Spot check
	var spots []string
	for i := 0; i < n; i += sampleRate {
		spots = append(spots, strconv.Itoa(sa[i]))
	}

	// Write output
	f, err := os.Create(outPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	fmt.Fprintf(f, "%v\n", mean)
	fmt.Fprintf(f, "%v\n", median)
	fmt.Fprintf(f, "%d\n", maxVal)
	fmt.Fprintf(f, "%s\n", strings.Join(spots, "\t"))
}
