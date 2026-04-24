package main

import (
	"bufio"
	"os"
)

func runQuery(bloomPath, queryPath string) error {
	f, err := os.Open(bloomPath)
	if err != nil {
		return err
	}
	defer f.Close()
	bf, err := ReadBloomFilter(bufio.NewReader(f))
	if err != nil {
		return err
	}

	qf, err := os.Open(queryPath)
	if err != nil {
		return err
	}
	defer qf.Close()

	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	s := bufio.NewScanner(qf)
	s.Buffer(make([]byte, 0, 64*1024), 1<<20)
	for s.Scan() {
		line := s.Bytes()
		if bf.Contains(line) {
			out.Write(line)
			out.WriteString("\tPROB_YES\n")
		} else {
			out.Write(line)
			out.WriteString("\tNO\n")
		}
	}
	return s.Err()
}
