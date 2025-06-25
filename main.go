package main

import (
	"log/slog"
	"os"
	"searchengine/db"
	memorymapper "searchengine/memory_mapper"
	"searchengine/tokenizer"
	"searchengine/utils"
	"strings"
)

// What are the problems here ??
// Poor Structure
// Error handling is messy
// Testing is messy

var docId int64 = 1
var docIdStore = "docId.gob"

// The main function retrieves the current working directory, creates a memory mapper, loads a document
// ID if it exists, tokenizes a message, and performs search and insert operations using the tokens.
func main() {

	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	newDb, err := db.NewDbService()
	if err != nil {
		panic(err)
	}
	defer newDb.Close()

	if utils.FileExists(docIdStore) {
		if docId, err = utils.LoadDocId(docIdStore); err != nil {
			panic(err)
		}
		lastId, err := newDb.GetDocId()
		if err != nil || strings.HasPrefix(err.Error(), db.NoEntryError) {
			panic(err)
		}
		if lastId != int(docId) {
			panic(err)
		}
	}

	mapper, err := memorymapper.NewMapper(path)
	if err != nil {
		panic(err)
	}
	defer mapper.Close()

	msg := "Labib Al Faisal"
	token := tokenizer.GetTokens(msg)

	Search(&token, mapper, newDb)
	Insert(&token, mapper, newDb)

}

// The `Search` function searches for documents based on tokens using a memory mapper and then
// retrieves the documents from a MySQL database.
func Search(token *tokenizer.Token, mapper *memorymapper.Mapper, dbService *db.DbService) []string {
	foundDocIds := make(map[uint64]struct{})
	for index, tok := range token.Tokens {
		tempSlice, err := mapper.GetDocIds(tok)
		if err != nil {
			continue
		}

		// AND operation on tempSlice and foundDocIds
		if index == 0 {
			for _, docId := range tempSlice {
				foundDocIds[docId] = struct{}{}
			}
			continue
		}

		commonDocId := make(map[uint64]struct{}, 0)
		for _, docId := range tempSlice {
			if _, ok := foundDocIds[docId]; ok {
				commonDocId[docId] = struct{}{}
			}
		}
		foundDocIds, commonDocId = commonDocId, foundDocIds
	}

	result := make([]string, 0, len(foundDocIds))
	for docId, _ := range foundDocIds {
		// Search from MySql
		document, err := dbService.Query(int(docId))
		if err != nil {
			slog.Info(err.Error())
			continue
		}
		result = append(result, document)
	}
	return result
}

// The Insert function increments the document ID, updates the memory mapper with tokens, and inserts a
// message into MySQL.
func Insert(token *tokenizer.Token, mapper *memorymapper.Mapper, dbService *db.DbService) {
	var insertedFlag bool = false
	for _, tok := range token.Tokens {
		if err := mapper.Update(tok, docId); err != nil {
			slog.Info(err.Error())
			continue
		}
		// Insert msg to MySQL
		if err := dbService.Insert(tok); err != nil {
			// To-do : Rollback mapper
			slog.Info(err.Error())
			continue
		}
		insertedFlag = true
	}
	if insertedFlag {
		docId++
	}
}
