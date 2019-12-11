package main

//type TokenKind int
//
//const (
//	PUNCTUATOR TokenKind = iota + 1
//	ANNOTATION
//)

type Token struct {
	value  string
	start  int
	length int
}

func lex(comment string) ([]*Token, error) {
	var (
		tok    string
		tokens []*Token
	)

	for index, char := range comment {
		switch char {
		case '(', ')', ':', '=', '[', ']', ',':
			if tok != "" {
				tokens = append(tokens, &Token{tok, index - len(tok) + 1, len(tok)})
				tok = ""
			}
			tokens = append(tokens, &Token{string(char), index, 1})
		case ' ':
			if tok != "" && tok != "..." {
				tokens = append(tokens, &Token{tok, index - len(tok) + 1, len(tok)})
				tok = ""
			}
		case '\n', '\r':
			if tok != "" {
				tokens = append(tokens, &Token{tok, index - len(tok) + 1, len(tok)})
				tok = ""
			}
		case '_', '@', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n',
			'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
			'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N',
			'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
			tok = tok + string(char)
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			tok = tok + string(char)
		case '/':
		default:
			//return nil, errors.New("invalid character '" + string(char) + "' at position " + strconv.Itoa(index))
		}
	}

	if tok != "" {
		tokens = append(tokens, &Token{tok, len(comment) - len(tok) + 1, len(tok)})
	}

	return tokens, nil
}
