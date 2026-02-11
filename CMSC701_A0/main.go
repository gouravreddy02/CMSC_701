package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"encoding/json"
)

type Stats struct {
	MinLen     int     `json:"min_len"`
	MaxLen     int     `json:"max_len"`
	MeanLen    float64 `json:"mean_len"`
	TotLen     int     `json:"tot_len"`
	NumRecords int     `json:"num_records"`
	CountA     int     `json:"count_a"`
	CountC     int     `json:"count_c"`
	CountG     int     `json:"count_g"`
	CountT     int     `json:"count_t"`
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: fasta_stats <input_file>")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	fmt.Fprintf(os.Stderr, "Processing file: %s\n", inputFile)

	file, err := os.Open(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	stats := Stats{
		MinLen: -1, 
		MaxLen: 0,
	}
	
	var currentSeq strings.Builder
	var hasCurrentSeq bool = false

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, ">") {
			if hasCurrentSeq {
				processSequence(&stats, currentSeq.String())
				currentSeq.Reset() 
			}

			hasCurrentSeq = true

		} else {
			currentSeq.WriteString(line)
		}
	}

	// Process the last sequence
	if hasCurrentSeq {
		processSequence(&stats, currentSeq.String())
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	if stats.NumRecords > 0 {
	stats.MeanLen = float64(stats.TotLen) / float64(stats.NumRecords)
	}

	// stats struct to JSON
	jsonData, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating JSON: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(jsonData))
}

// updates statistics for a each sequence
func processSequence(stats *Stats, sequence string) {
	seq_len := len(sequence)
	stats.NumRecords++

	stats.TotLen += seq_len
	if stats.MinLen == -1 || seq_len < stats.MinLen {
		stats.MinLen = seq_len
	}

	if seq_len > stats.MaxLen {
		stats.MaxLen = seq_len
	}

	for _, nucleotide := range sequence {
		switch nucleotide {
		case 'A', 'a':
			stats.CountA++
		case 'C', 'c':
			stats.CountC++
		case 'G', 'g':
			stats.CountG++
		case 'T', 't':
			stats.CountT++
		}
	}
}
