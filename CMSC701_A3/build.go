package main

import (
	"bufio"
	cryptorand "crypto/rand"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
)

func runBuild(fprStr, inputPath, outputPath string) error {
	fpr, err := strconv.ParseFloat(fprStr, 64)
	if err != nil {
		return fmt.Errorf("parse fpr: %w", err)
	}
	if !(fpr > 0 && fpr < 1) {
		return fmt.Errorf("fpr must be in (0, 1), got %v", fpr)
	}

	keys, err := readDistinctKeys(inputPath)
	if err != nil {
		return err
	}
	n := uint64(len(keys))
	m, k := OptimalParams(n, fpr)

	seeds, err := randomSeeds(k)
	if err != nil {
		return err
	}

	bf := NewBloomFilter(n, m, k, seeds)
	for key := range keys {
		bf.Insert([]byte(key))
	}

	out, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer out.Close()
	bw := bufio.NewWriter(out)
	if err := bf.Write(bw); err != nil {
		return err
	}
	if err := bw.Flush(); err != nil {
		return err
	}

	fmt.Printf("n\t%d\n", n)
	fmt.Printf("m\t%d\n", m)
	fmt.Printf("k\t%d\n", k)
	return nil
}

func readDistinctKeys(path string) (map[string]struct{}, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	keys := make(map[string]struct{})
	s := bufio.NewScanner(f)
	s.Buffer(make([]byte, 0, 64*1024), 1<<20)
	for s.Scan() {
		line := s.Text()
		if line == "" {
			continue
		}
		keys[line] = struct{}{}
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	return keys, nil
}

func randomSeeds(k uint64) ([]uint64, error) {
	buf := make([]byte, 8*k)
	if _, err := cryptorand.Read(buf); err != nil {
		return nil, err
	}
	seeds := make([]uint64, k)
	for i := range seeds {
		seeds[i] = binary.LittleEndian.Uint64(buf[i*8 : i*8+8])
	}
	return seeds, nil
}
