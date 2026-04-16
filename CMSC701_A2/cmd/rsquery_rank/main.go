package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"cmsc701_a2/bitvector"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintf(os.Stderr, "Usage: %s <index> <query_file> <output_path>\n", os.Args[0])
		os.Exit(1)
	}

	ri, err := bitvector.LoadIndex(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading index: %v\n", err)
		os.Exit(1)
	}

	qf, err := os.Open(os.Args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening query file: %v\n", err)
		os.Exit(1)
	}
	defer qf.Close()

	of, err := os.Create(os.Args[3])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer of.Close()

	w := bufio.NewWriter(of)
	defer w.Flush()

	scanner := bufio.NewScanner(qf)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) != 2 || parts[0] != "rank" {
			continue
		}
		idx, err := strconv.ParseUint(parts[1], 10, 64)
		if err != nil {
			continue
		}

		if idx > ri.BV.N {
			fmt.Fprintf(w, "%d:-1\n", idx)
		} else {
			fmt.Fprintf(w, "%d:%d\n", idx, ri.Rank(idx))
		}
	}
}