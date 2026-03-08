package main

import (
	"fmt"
	"os"
	"strconv"

	"cmsc701_a1/pkg/fasta"
	"cmsc701_a1/pkg/suffixarray"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintf(os.Stderr, "Usage: buildsa <reference> <k> <output>\n")
		os.Exit(1)
	}

	refPath := os.Args[1]
	k, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid k value: %s\n", os.Args[2])
		os.Exit(1)
	}
	outPath := os.Args[3]

	//Read the FASTA reference
	seq, err := fasta.ReadReference(refPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading reference: %v\n", err)
		os.Exit(1)
	}

	//Append sentinel
	text := seq + "$"

	//Build suffix array
	sa := suffixarray.BuildSuffixArray(text)

	//Build prefix table
	prefixTable := suffixarray.BuildPrefixTable(text, sa, k)

	//Save to binary file
	idx := &suffixarray.Index{
		Text:        text,
		SA:          sa,
		K:           k,
		PrefixTable: prefixTable,
	}
	err = suffixarray.Save(idx, outPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error saving index: %v\n", err)
		os.Exit(1)
	}
}
