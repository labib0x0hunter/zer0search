package utils

import (
	"fmt"
	"hash/maphash"
	"path/filepath"
)

var seedFile = "seed.gob"

type Hash struct {
	h *maphash.Hash
}

func NewHash(path string) (*Hash, error) {
	var seed maphash.Seed
	var err error

	fmt.Println("Inside HASH -> ", filepath.Join(path, seedFile))

	if !FileExists(filepath.Join(path, seedFile)) {
		seed = maphash.MakeSeed()
		if err = saveSeed(filepath.Join(path, seedFile), seed); err != nil {
			fmt.Println(seed)
			return nil, err
		}
	} else {
		if seed, err = loadSeed(filepath.Join(path, seedFile)); err != nil {
			return nil, err
		}
	}

	var h maphash.Hash
	h.SetSeed(seed)
	return &Hash{h: &h}, nil
}

func (h *Hash) WriteString(msg string) error {
	if _, err := h.h.WriteString(msg); err != nil {
		return err
	}
	return nil
}

func (h *Hash) Sum() uint64 {
	return h.h.Sum64()
}

func (h *Hash) Reset() {
	h.h.Reset()
}
