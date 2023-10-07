/*
To do:

New module according to Go Dev blog practices (internal folder)
Tests, tests, more tests
Package Lex
factor out Token stuff into token package
maybe refactor Lex to be recursive?
fix error handling in Lexer!!!
add single pass to recover line, col# information
fix the logic in the pretty printer (low priority though)
add floating point, scientific support
verify naming convention: no whitespace or non-printing characters
generate a lot more lisp lists to test
*/

package lex

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"text/tabwriter"
)

type TokenType int

const (
	debug = false

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
	// capture group to get only the comment body.
	patComment = `(?m)^;+(.*)$`

	// look for ' plus any whitespace after
	patSingleQuote = `'\s*`

	patOpenParen = `^\(`

	patCloseParen = `^\)`

	//Pattern for string literal

	patInteger = `^(-?\d*)`

    patSymbol = `^[^\d\s':][^\s]+`
)

func TokenTypeDict() func(TokenType) string {
	d := map[TokenType]string{
		TOK_OPENPAREN:  "TOK_OPENPAREN",
		TOK_CLOSEPAREN: "TOK_CLOSEPAREN",
		TOK_SYMBOL:     "TOK_SYMBOL",
		TOK_STRING:     "TOK_STRING",
		TOK_COMMENT:    "TOK_COMMENT",
		TOK_INTEGER:    "TOK_INTEGER",
		//TOK_LISTQUOTE:   "TOK_LISTQUOTE",
		//TOK_COMMASPLICE: "TOK_COMMASPLICE",
		//TOK_FLOAT:       "TOK_FLOAT",
	}

	return func(key TokenType) string {
		return d[key]
	}
}

// Holds location information at a point in a file for debugging
// and error messages.
type Location struct {
	Line  int
	Col   int
	Index int
}

type Token struct {
	Token string
	Typ   TokenType
	Loc   int
}

func (t *Token) String() string {
	return fmt.Sprintf("{Loc:\t%d\tType:\t%s\tToken:\t%q\t}",
		t.Loc,
		TokenTypeDict()(t.Typ),
		t.Token)
}

func (t *Token) isParen() bool {
	return t.Typ == TOK_OPENPAREN || t.Typ == TOK_CLOSEPAREN
}

func printTokens(tokens []Token) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	for _, t := range tokens {
		fmt.Fprintln(w, &t)
	}
	w.Flush()
}

func prettyPrintTokens(tokens []Token) {
	var sb strings.Builder

	for i, t := range tokens {
		if t.Typ == TOK_COMMENT {
			continue
		}
		sb.WriteString(t.Token)

		if t.Typ == TOK_OPENPAREN {
			continue
		}
		peekOK := i+1 < len(tokens)
		if !peekOK {
			continue
		}

		if !t.isParen() && tokens[i+1].Typ == TOK_CLOSEPAREN {
			continue
		}
		sb.WriteString(" ")
	}
	fmt.Println(sb.String())
}

// regexp.<Index> match methods return slices of ints, either
// a nil slice, or an even number of ints that denote pairs
// of start and end indices.
func reIndexGroup(groups []int, group int) int {
	i := group * 2
	return groups[i+1]
}

type Lexer interface {
	Add(pat string, typ TokenType)
	Lex(text string) []Token
}

type lexerObject struct {
	text     string
	tokens   []Token
	patterns []pattern
	loc      int
}

func NewLexer() Lexer {
	const defaultCapacity = 128

	tokens := make([]Token, 0, defaultCapacity)
	lo := &lexerObject{}
	lo.tokens = tokens

	return &lexerObject{}
}

type pattern struct {
	re  *regexp.Regexp
	typ TokenType
}

func (lex *lexerObject) Add(pat string, typ TokenType) {
	term := regexp.MustCompile(pat)
	lex.patterns = append(lex.patterns, pattern{term, typ})
}

func (lex *lexerObject) Lex(text string) []Token {

	if lex == nil {
		panic("nil lexerObject: Must instantiate lexer with NewLexer()")
	}
	if lex.patterns == nil {
		panic("can't lex without regex actions attached!")
	}

	lex.loc = 0
	lex.text = text

	for len(text) > 0 {
	scan:
		for _, pattern := range lex.patterns {

			re := pattern.re

			result := re.FindStringIndex(text)

			if result == nil {
				continue scan
			}

			start, end := result[0], result[1]

			// store the location
			token := Token{text[start:end], pattern.typ, lex.loc}

			lex.tokens = append(lex.tokens, token)

			skip := end
			if skip >= len(text) {
				panic("Tried to skip too many runes")
			}

			lex.loc += skip

			//skip to next unmatched portion
			text = text[skip:len(text)]
		}
		// haven't found anything; failed to find a match
		return nil
	}
	return lex.tokens
}

func main() {
    fmt.Println("hello world")
}
