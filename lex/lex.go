/*
To do:

*******Tests, tests, more tests
factor out Token stuff into token package
fix error handling in Lexer!!!
add single pass to recover line, col# information
fix the logic in the pretty printer (low priority though)
add floating point, scientific support
verify naming convention: no whitespace or non-printing characters
generate a lot more lisp lists to test
*/

package lex

import (
	"regexp"
)

type Lexer interface {
	//Add(pat string, typ TokenType)
	Lex(text string) []Token
}

type lexerObject struct {
	text     string
	tokens   []Token
	patterns []pattern
	compiled []*regexp.Regexp
	loc      int
}

type pattern struct {
	pat string
	typ TokenType
}

func NewLexer(patterns []pattern) *lexerObject {
	const defaultCapacity = 128

	lo := &lexerObject{}

    lo.tokens = make([]Token, 0, defaultCapacity)
	lo.patterns = patterns

	lo.compiled = make([]*regexp.Regexp, 0, len(patterns))
	for _, p := range lo.patterns {
        lo.compiled = append(lo.compiled, regexp.MustCompile(p.pat))
	}

	return lo
}

func (lex *lexerObject) lexOnce(here string) (token Token, skip int, success bool) {
	for i, re := range lex.compiled {

		result := re.FindStringIndex(here)

		if result == nil {
			continue
		}

		// FindStringIndex returns either nil or a pair of index integers
		_, end := result[0], result[1]

		// store the location
		token = Token{here[0:end], lex.patterns[i].typ, lex.loc}

		return token, end, true
	}
	return Token{}, 0, false
}

func (lex *lexerObject) init(text string) {
	lex.loc = 0
	lex.text = text
}

func (lex *lexerObject) Lex(text string) []Token {
	panic("TODO: implement")
}

/*
func (lex *lexerObject) Lex(text string) []Token {

	if lex == nil {
		panic("nil lexerObject: Must instantiate lexer with NewLexer()")
	}
	if lex.patterns == nil {
		panic("can't lex without regex actions attached!")
	}

    lex.init(text)

	for len(text) > 0 {
	scan:
    tok, skip, success := lex.lexOnce(text)

*/
