package memorymapper

import (
	"errors"
	"fmt"
	"path/filepath"
	"searchengine/memory_mapper/dictionary"
	"searchengine/memory_mapper/posting"
	"searchengine/utils"
)

type Mapper struct {
	dict   *dictionary.Dictionary
	post   *posting.Posting
	hasher *utils.Hash
}

func NewMapper(path string) (*Mapper, error) {
	dict, err := dictionary.NewDictionary(path)
	if err != nil {
		return nil, err
	}

	post, err := posting.NewPosting(path)
	if err != nil {
		return nil, err
	}

	hasher, err := utils.NewHash(filepath.Join(path, "utils"))
	if err != nil {
		fmt.Println("HERE -> NewMapper : ", filepath.Join(path, "utils"))
		return nil, err
	}

	return &Mapper{
		dict:   dict,
		post:   post,
		hasher: hasher,
	}, nil
}

// search word in dictionary.index
// true : get docIds from posting.index, append new docId and update dictionary.index
// false : append docId in posting.index and append postingOffset in dictionary.index
func (m *Mapper) Update(word string, docId int64) error {
	wordHash := m.getHash(word)
	found, offset, postingOffset, postingLen, err := m.dict.Search(wordHash)
	if err != nil {
		return err
	}
	if !found {
		postingOffset, err = m.post.Append(uint64(docId), true)
		if err != nil {
			return err
		}
		return m.dict.Append(wordHash, postingOffset, uint64(1))
	}

	postingOffset, err = m.post.Update(postingOffset, postingLen, uint64(docId))
	if err != nil {
		return err
	}

	return m.dict.Update(offset, postingOffset, postingLen+1)
}

// search word in dictionary.index
// get docIds from posting.index
func (m *Mapper) GetDocIds(word string) ([]uint64, error) {
	wordHash := m.getHash(word)
	found, _, postingOffset, postingLen, err := m.dict.Search(wordHash)
	if err != nil {
		return []uint64{}, err
	}
	if !found {
		return []uint64{}, nil
	}
	return m.post.Search(postingOffset, postingLen)
}

// hash word
func (m *Mapper) getHash(word string) uint64 {
	m.hasher.WriteString(word)
	wordHash := m.hasher.Sum()
	m.hasher.Reset()
	return wordHash
}

// Close
func (m *Mapper) Close() error {
	errDict := m.dict.Close()
	errPost := m.post.Close()
	if errDict != nil || errPost != nil {
		return errors.New("[error] : " + errDict.Error() + " " + errPost.Error())
	}
	return nil
}
