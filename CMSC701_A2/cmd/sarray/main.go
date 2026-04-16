package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"cmsc701_a2/bitvector"
)

type entry struct {
	pos uint64
	key string
}

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <command_script> <output_path>\n", os.Args[0])
		os.Exit(1)
	}

	sf, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening script: %v\n", err)
		os.Exit(1)
	}
	defer sf.Close()

	of, err := os.Create(os.Args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output: %v\n", err)
		os.Exit(1)
	}
	defer of.Close()

	w := bufio.NewWriter(of)
	defer w.Flush()

	scanner := bufio.NewScanner(sf)
	// Increase buffer for potentially long lines
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)

	var bv *bitvector.BitVector
	var entries []entry
	var ri *bitvector.RankIndex
	var values []string
	var n uint64

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, " ", 3)
		cmd := parts[0]

		switch cmd {
		case "init":
			n, _ = strconv.ParseUint(parts[1], 10, 64)
			bv = bitvector.NewBitVector(n)
			entries = nil

		case "insert":
			pos, _ := strconv.ParseUint(parts[1], 10, 64)
			key := parts[2]
			entries = append(entries, entry{pos, key})

		case "finish":
			sort.Slice(entries, func(i, j int) bool {
				return entries[i].pos < entries[j].pos
			})
			values = make([]string, len(entries))
			for i, e := range entries {
				bv.SetBit(e.pos)
				values[i] = e.key
			}
			ri = bitvector.BuildRankIndex(bv)

		case "query_index":
			idx, _ := strconv.ParseUint(parts[1], 10, 64)
			if bv.AccessBit(idx) == 0 {
				fmt.Fprintf(w, "idx:%d:-1\n", idx)
			} else {
				r := ri.Rank(idx)
				fmt.Fprintf(w, "idx:%d:%s\n", idx, values[r])
			}

		case "query_rank":
			r, _ := strconv.ParseUint(parts[1], 10, 64)
			k := uint64(len(values))
			if r >= k {
				fmt.Fprintf(w, "qr:%d:-1\n", r)
			} else {
				fmt.Fprintf(w, "qr:%d:%s\n", r, values[r])
			}

		case "rank_at_index":
			idx, _ := strconv.ParseUint(parts[1], 10, 64)
			fmt.Fprintf(w, "rai:%d:%d\n", idx, ri.Rank(idx))

		case "index_of_next_key":
			idx, _ := strconv.ParseUint(parts[1], 10, 64)
			// rank(idx+1) = number of 1s in [0, idx+1) = [0, idx]
			r := ri.Rank(idx + 1)
			// next 1-bit is the (r+1)-th one; use Select1Pos to get
			// the actual position of that 1-bit
			pos := bitvector.Select1Pos(ri, r+1)
			if pos < 0 {
				fmt.Fprintf(w, "ink:%d:-1\n", idx)
			} else {
				fmt.Fprintf(w, "ink:%d:%d\n", idx, pos)
			}

		case "density":
			k := uint64(len(values))
			density := float64(k) / float64(n)
			fmt.Fprintf(w, "%s\n", strconv.FormatFloat(density, 'f', -1, 64))
		}
	}
}
