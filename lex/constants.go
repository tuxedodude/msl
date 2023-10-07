package lex

const (
	TOK_OPENPAREN TokenType = iota
	TOK_CLOSEPAREN
	TOK_SYMBOL
	TOK_STRING
	TOK_LISTQUOTE
	TOK_COMMASPLICE
	TOK_COMMENT
	TOK_INTEGER
	TOK_FLOAT

	patWhiteSpace = `^\s+`

	// comments start at ';' and go to end of line only.
	patComment = `^;+(.*)$`

	patSingleQuote = `^'`

	patOpenParen = `^\(`

	patCloseParen = `^\)`

	patInteger = `^(0|(-?[1-9]\d*))`

	patSymbol = `^[^\d\s':#][^\s]+`

	patString = `^"([^"\n\r\t]|(\\["\n\r\t]))*"`
)