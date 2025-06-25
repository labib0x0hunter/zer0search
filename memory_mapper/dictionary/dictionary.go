package dictionary

import (
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/tysonmote/gommap"
)

const (
	dictionaryIndexFile = "dictionary.index"
	lenSize             = 8
	entrySize           = 24       // [hash][offset][postingLen]
	MaxFileSize         = 10485760 // 10Mb
)

var (
	encoder = binary.BigEndian
)

type Dictionary struct {
	file   *os.File
	mmap   gommap.MMap // mmap
	len    uint64      // current size
	closed bool        // flag to check if the directory.index is closed
}

func NewDictionary(path string) (*Dictionary, error) {

	// os.Remove(filepath.Join(path, dictionaryIndexFile))

	dict := &Dictionary{}

	file, err := os.OpenFile(filepath.Join(path, dictionaryIndexFile), os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	dict.closed = false

	info, err := file.Stat()
	if err != nil {
		return nil, err
	}

	size := info.Size()

	dict.len = uint64(size)

	// allocate capacity
	if err := file.Truncate(int64(MaxFileSize)); err != nil {
		return nil, err
	}

	dict.file = file

	mmap, err := gommap.Map(
		dict.file.Fd(),
		gommap.PROT_READ|gommap.PROT_WRITE,
		gommap.MAP_SHARED,
	)

	if err != nil {
		file.Close()
		return nil, err
	}

	dict.mmap = mmap

	return dict, nil
}

// search in dictionary.index
// [wordHash][postingOffset][postingLength]
// [uint64][uint64][uint64]
// [8][8][8] -> 24bytes
// returns offsetOfWord, postingOffset, postingLen
func (d *Dictionary) Search(wordHash uint64) (bool, uint64, uint64, uint64, error) {
	if d.closed {
		return false, uint64(0), uint64(0), uint64(0), errors.New("file is closed")
	}
	for j := uint64(0); j+entrySize <= d.len; j += entrySize {
		offset := j
		storedHash := encoder.Uint64(d.mmap[offset : offset+lenSize])
		offset += lenSize
		postingOffset := encoder.Uint64(d.mmap[offset : offset+lenSize])
		offset += lenSize
		positionLength := encoder.Uint64(d.mmap[offset : offset+lenSize])

		fmt.Println("[debug] inside search: ", storedHash)

		if storedHash == wordHash {
			return true, uint64(j), postingOffset, positionLength, nil
		}
	}
	return false, uint64(0), uint64(0), uint64(0), nil
}

// append in dictionary.index
func (d *Dictionary) Append(hash, postingOffset, postingLen uint64) error {
	if d.closed {
		return errors.New("file is closed")
	}
	if d.IsFilled(entrySize) {
		return errors.New("max filesize reached")
	}

	offset := d.len
	d.Debug(offset)
	encoder.PutUint64(d.mmap[offset:offset+lenSize], hash)
	offset += lenSize
	d.Debug(offset)

	encoder.PutUint64(d.mmap[offset:offset+lenSize], postingOffset)
	offset += lenSize
	d.Debug(offset)

	encoder.PutUint64(d.mmap[offset:offset+lenSize], postingLen)
	offset += lenSize
	d.Debug(offset)

	fmt.Println()

	d.len += entrySize
	return nil
}

// update postingOffset and postingLen stored at offset
func (d *Dictionary) Update(offset, postingOffset, postingLen uint64) error {
	if d.closed {
		return errors.New("file is closed")
	}
	if offset+entrySize > d.len {
		return errors.New("[error] : SGMNT_FLT")
	}

	initialOffset := offset + lenSize

	encoder.PutUint64(d.mmap[initialOffset:initialOffset+lenSize], postingOffset)
	initialOffset += lenSize
	encoder.PutUint64(d.mmap[initialOffset:initialOffset+lenSize], postingLen)
	return nil
}

// check if there is enough space with size
func (d *Dictionary) IsFilled(size uint64) bool {
	return d.len+size > MaxFileSize
}

// close the dictionary.index
func (d *Dictionary) Close() error {
	if d.closed {
		return errors.New("file is closed")
	}
	if err := d.mmap.Sync(gommap.MS_SYNC); err != nil {
		return err
	}
	if err := d.file.Sync(); err != nil {
		return err
	}
	if err := d.file.Truncate(int64(d.len)); err != nil {
		return err
	}
	if err := d.file.Sync(); err != nil {
		return err
	}
	if err := d.mmap.UnsafeUnmap(); err != nil {
		return err
	}
	if err := d.file.Close(); err != nil {
		return err
	}
	d.closed = true
	d.mmap = nil
	d.file = nil
	return nil
}

// debug information
func (d *Dictionary) Debug(u uint64) {
	fmt.Println("bytes written: ", u)
}
