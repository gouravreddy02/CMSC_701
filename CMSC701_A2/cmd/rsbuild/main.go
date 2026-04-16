package main

import (
	"fmt"
	"os"
	"cmsc701_a2/bitvector"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <input_bitvector_file> <output_index_path>\n", os.Args[0])
		os.Exit(1)
	}

	bv, err := bitvector.ReadBitvector(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading bitvector: %v\n", err)
		os.Exit(1)
	}

	ri := bitvector.BuildRankIndex(bv)

	if err := bitvector.SaveIndex(os.Args[2], ri); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving index: %v\n", err)
		os.Exit(1)
	}
}