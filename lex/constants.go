package lex

const (
	TOK_NONE TokenType = iota
	TOK_WHITESPACE
	TOK_OPENPAREN
	TOK_CLOSEPAREN
	TOK_SYMBOL
	TOK_STRING
	TOK_SINGLEQUOTE
	TOK_COMMASPLICE
	TOK_COMMENT
	TOK_INTEGER
	TOK_FLOAT

	patWhiteSpace = `^\s+`

	// comments start at ';' and go to end of line only.
	// note: (?m) is necessary for proper end of line behavior.
	patComment = `(?m)^;[^\n]*`

	patSingleQuote = `^'`

	patOpenParen = `^\(`

	patCloseParen = `^\)`

	patInteger = `^(0|(-?[1-9]\d*))`

	patSymbol = `^[^\d\s':#"][^\s\)\(]*`

	patString = `^"([^"\n\r\t]|(\\["\n\r\t]))*"`
)
