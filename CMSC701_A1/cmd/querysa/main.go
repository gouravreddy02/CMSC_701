package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"cmsc701_a1/pkg/fasta"
	"cmsc701_a1/pkg/suffixarray"
)

func main() {
	if len(os.Args) != 5 {
		fmt.Fprintf(os.Stderr, "Usage: querysa <index> <queries> <mode> <output>\n")
		os.Exit(1)
	}

	idxPath := os.Args[1]
	queryPath := os.Args[2]
	mode := os.Args[3]
	outPath := os.Args[4]

	// Load the index
	idx, err := suffixarray.Load(idxPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading index: %v\n", err)
		os.Exit(1)
	}

	// Read query records
	queries, err := fasta.ReadRecords(queryPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading queries: %v\n", err)
		os.Exit(1)
	}

	// Open output file
	f, err := os.Create(outPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	// Process each query
	for _, q := range queries {
		var lb, ub int
		var hits []int

		switch mode {
		case "naive":
			lb, ub, hits = suffixarray.QueryNaive(idx.Text, idx.SA, q.Sequence)
		case "simpaccel":
			lb, ub, hits = suffixarray.QuerySimpaccel(idx.Text, idx.SA, q.Sequence)
		case "prefaccel":
			lb, ub, hits = suffixarray.QueryPrefaccel(idx.Text, idx.SA, idx.PrefixTable, idx.K, q.Sequence)
		default:
			fmt.Fprintf(os.Stderr, "Unknown mode: %s\n", mode)
			os.Exit(1)
		}

		// Format: name  lb  ub  k  hit1  hit2  ...
		parts := []string{q.Name, strconv.Itoa(lb), strconv.Itoa(ub), strconv.Itoa(len(hits))}
		for _, h := range hits {
			parts = append(parts, strconv.Itoa(h))
		}
		fmt.Fprintf(f, "%s\n", strings.Join(parts, "\t"))
	}
}
