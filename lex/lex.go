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

func defaultLexerPatterns() []pattern {
	return []pattern{
		{pat: patWhiteSpace, typ: TOK_WHITESPACE},
		{pat: patComment, typ: TOK_COMMENT},
		{pat: patSingleQuote, typ: TOK_SINGLEQUOTE},
		{pat: patOpenParen, typ: TOK_OPENPAREN},
		{pat: patCloseParen, typ: TOK_CLOSEPAREN},
		{pat: patInteger, typ: TOK_INTEGER},
		{pat: patSymbol, typ: TOK_SYMBOL},
		{pat: patString, typ: TOK_STRING},
	}
}

func NewLexer() *lexerObject {
	return newLexer(defaultLexerPatterns())
}

// helper method for NewLexer
// allows unit testing
func newLexer(patterns []pattern) *lexerObject {
	if patterns == nil {
		patterns = defaultLexerPatterns()
	}

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

	//defer fmt.Println()

	//td := TokenTypeDict()

	for i, re := range lex.compiled {
		//fmt.Printf("%d. Trying %s: %q\n", i, td(lex.patterns[i].typ), lex.patterns[i].pat)

		result := re.FindStringIndex(here)

		if result == nil {
			continue
		}

		// FindStringIndex returns either nil or a pair of index integers
		_, end := result[0], result[1]

		t := here[0:end]

		//fmt.Printf("\tsucceeded: %q\n", t)

		// store the location
		token = Token{t, lex.patterns[i].typ, lex.loc}

		return token, end, true
	}
	return Token{}, 0, false
}

func (lex *lexerObject) init(text string) {
	lex.loc = 0
	lex.text = text
}

func (lex *lexerObject) Lex(text string) []Token {

	if lex == nil {
		panic("nil lexerObject: Must instantiate lexer with NewLexer()")
	}
	if lex.patterns == nil {
		panic("can't lex without regex actions attached!")
	}

	lex.init(text)

	dbgcount := 0

	for len(text) > 0 {
		//fmt.Printf("Lexer Pass %d on %q\n", dbgcount, text)

		tok, skip, found := lex.lexOnce(text)

		if !found {
			return nil
		}

		tok.Loc = lex.loc
		lex.tokens = append(lex.tokens, tok)

		lex.loc += skip

		if skip == 0 {
			panic("Found a match but skipped 0 runes!")
		}

		text = text[skip:len(text)]

		dbgcount++
	}
	return lex.tokens
}
