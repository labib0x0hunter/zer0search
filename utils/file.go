package utils

import (
	"encoding/gob"
	"hash/maphash"
	"os"
)

func FileExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func saveSeed(filename string, seed maphash.Seed) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	return encoder.Encode(seed)
}

func loadSeed(filename string) (maphash.Seed, error) {
	file, err := os.Open(filename)
	if err != nil {
		return maphash.Seed{}, err
	}
	defer file.Close()

	var seed maphash.Seed
	decoder := gob.NewDecoder(file)
	if err = decoder.Decode(&seed); err != nil {
		return maphash.Seed{}, err
	}
	return seed, nil
}

func SaveDocId(filename string, docId int64) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	return encoder.Encode(docId)
}

func LoadDocId(filename string) (int64, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	var docId int64
	decoder := gob.NewDecoder(file)
	if err = decoder.Decode(&docId); err != nil {
		return 0, err
	}
	return docId, nil
}