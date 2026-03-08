package fasta

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Record struct {
	Name     string
	Sequence string
}

// ReadRecords reads a FASTA file and returns all records.
func ReadRecords(path string) ([]Record, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var records []Record
	var current *Record
	var seqBuilder strings.Builder

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue
		}
		if line[0] == '>' {
			if current != nil {
				current.Sequence = seqBuilder.String()
				records = append(records, *current)
			}
			current = &Record{Name: line[1:]}
			seqBuilder.Reset()
		} else {
			seqBuilder.WriteString(strings.ToUpper(line))
		}
	}
	if current != nil {
		current.Sequence = seqBuilder.String()
		records = append(records, *current)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return records, nil
}

// ReadReference reads a FASTA file and returns just the sequence as a single string.
func ReadReference(path string) (string, error) {
	records, err := ReadRecords(path)
	if err != nil {
		return "", err
	}
	if len(records) == 0 {
		return "", fmt.Errorf("no records found in FASTA file: %s", path)
	}
	return records[0].Sequence, nil
}