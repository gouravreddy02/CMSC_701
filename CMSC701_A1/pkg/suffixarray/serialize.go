package suffixarray

import (
	"encoding/gob"
	"os"
)

// Index holds everything needed to query the suffix array.
type Index struct {
	Text        string
	SA          []int
	K           int
	PrefixTable map[string][2]int
}

// Save serializes the index to a binary file using gob encoding.
func Save(idx *Index, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return gob.NewEncoder(f).Encode(idx)
}

// Load deserializes the index from a binary file.
func Load(path string) (*Index, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var idx Index
	err = gob.NewDecoder(f).Decode(&idx)
	if err != nil {
		return nil, err
	}
	return &idx, nil
}
