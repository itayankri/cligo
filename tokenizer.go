package main

import (
	"github.com/pkg/errors"
	"strconv"
)

type TokenKind int

const (
	NAME TokenKind = iota + 1
	STRING
)

type Token struct {
	kind   TokenKind
	value  string
	start  int
	length int
}

func lex(comment string) ([]*Token, error) {
	var (
		tok       string
		tokens    []*Token
		tokenKind TokenKind
	)

	for i := 0; i < len(comment); i++ {
		switch comment[i] {
		case '"':
			{
				tokenKind = STRING
				tok += string(comment[i])

				for j := i + 1; j < len(comment); j++ {
					tok += string(comment[j])

					if comment[j] == '"' {
						i = j
						break
					}
				}

				stringValue, err := strconv.Unquote(tok)
				if err != nil {
					return nil, errors.Wrap(err, "bad quoted string")
				} else {
					tok = stringValue
				}
			}
		case '_', ':', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n',
			'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
			'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N',
			'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
			{
				tok += string(comment[i])
				tokenKind = NAME
			}
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			{
				if tok == "" {
					return nil, errors.New("bad name")
				} else {
					tok += string(comment[i])
				}
			}
		case ' ', '\t':
			{
				if tok != "" {
					tokens = append(tokens, &Token{tokenKind, tok, i - len(tok) + 1, len(tok)})
					tok = ""
				}
			}
		case '/':
		default:
			{
				return nil, errors.New("invalid character - " + string(comment[i]))
			}
		}
	}

	if tok != "" {
		tokens = append(tokens, &Token{tokenKind, tok, len(comment) - len(tok) + 1, len(tok)})
	}

	return tokens, nil
}
