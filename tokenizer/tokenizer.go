package tokenizer

import "strings"

func GetTokens(msg string) Token {
	msg = strings.TrimSpace(msg)
	var token Token
	token.Tokens = strings.Fields(msg)
	token.toLower()
	token.removePunctuation()
	return token
}