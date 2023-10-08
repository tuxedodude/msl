package lex

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestRegex(t *testing.T) {

	assert := assert.New(t)

	// test case
	type tc struct {
		s    string
		want bool
	}

	var patternTests = []struct {
		pat   string
		desc  string
		tests []tc
	}{
		{
			pat:  patWhiteSpace,
			desc: "whitespace",
			tests: []tc{
				tc{"", false},
				tc{" ", true},
			},
		},
		{
			pat:  patComment,
			desc: "Comment",
			tests: []tc{
				tc{"", false},
				tc{" ", false},
				tc{"; hello", true},
				tc{" ; hello", false},
				tc{"bkdlfj;jfkdls", false},
				tc{";hello\n(lambda (x) (* x x))", true},
			},
		},
		{
			pat:  patSingleQuote,
			desc: "Singlequote (')",
			tests: []tc{
				tc{"", false},
				tc{" ", false},
				tc{"'", true},
				tc{" '", false},
			},
		},
		{
			pat:  patInteger,
			desc: "Integer",
			tests: []tc{
				tc{"", false},
				tc{" ", false},
				tc{" 1", false},
				tc{" -1", false},
				tc{"-", false},
				tc{"-1", true},
				tc{"-0", false},
				tc{" -0", false},
				tc{"-01234", false},
				tc{"-1234", true},
				tc{"1230123", true},
				tc{"123-123 fj kslf", true},
			},
		},
		{
			pat:  patSymbol,
			desc: "Symbol",
			tests: []tc{
				tc{"", false},
				tc{" ", false},
				tc{"'", false},
				tc{"'foobar", false},
				tc{"#anytthing", false},
				tc{":tag", false},
				tc{" blah blah foo", false},
				tc{"blah blah foo", true},
				tc{"a1234", true},
				tc{"a_stpiuf>JFd", true},
				tc{`"anything"`, false},
				tc{`baz)`, true},
				tc{`x) (* x x))`, true},
			},
		},
		{
			pat:  patString,
			desc: "Double quoted string (\"\")",
			tests: []tc{
				tc{"", false},
				tc{`"`, false},
				tc{`""`, true},
				tc{` "`, false},
				tc{` ""`, false},
				tc{`abcd`, false},
				tc{`"abcd"`, true},
				tc{`"a b c d `, false},
				tc{`"a b c d"`, true},
				// note: if we can't pass the test, then "hello""hello" will make it thru the lexer
				// as two tokens.
				//tc{`"abcd""`, false},
				tc{`"abcd")`, true},
				tc{"\"abc\tdef\"", false}, // tabs!
			},
		},
	}

	for _, p := range patternTests {
		re := regexp.MustCompile(p.pat)

		for _, test := range p.tests {
			result := re.MatchString(test.s)

			errstr := func(msg string) string {
				return fmt.Sprintf("%s regex pattern %s %s on %q", p.desc, p.pat, msg, test.s)
			}

			if test.want {
				assert.True(result, errstr("failed to match"))
			} else {
				assert.False(result, errstr("matched incorrectly"))
			}
		}
	}
}

func TestNewLexer(t *testing.T) {
	assert := assert.New(t)

	pats := defaultLexerPatterns()
	lex := newLexer(pats)

	assert.NotNil(lex.tokens, "lex object tokens are not nil")
	assert.NotNil(lex.patterns, "lex object patterns are not nil")
	assert.Equal(len(lex.patterns), len(lex.compiled), "lex object patterns and compiled regex have same number of objects")
	assert.Equal(len(lex.patterns), len(pats), "lexerObject.tokens [] has same length as patterns passed in")
	assert.Equal(len(lex.patterns), len(pats), "lexerObject.compiled [] has same length as patterns passed in")
	assert.Equal(len(lex.patterns), 8, "proper number of patterns in lex object")
	for _, re := range lex.compiled {
		assert.NotNil(re, "compiled regex objects in lexer are not nil")
	}

	foo := "foo"
	lex.init(foo)

	assert.Equal(lex.loc, 0, "after init, lexerObject has location 0")
	assert.Equal(lex.text, foo, "foo string is stored inside lexer object after init")

	tok, skip, success := lex.lexOnce(foo)
	assert.True(success, "lexer should have parsed `foo`")
	assert.Equal(skip, len(foo), "should have skipped the length of "+foo)
	assert.Equal(tok.Token, foo, "should parse "+foo)
	assert.Equal(tok.Typ, TOK_SYMBOL, "should have parsed "+foo+"as symbol")

	tests := []struct {
		s      string
		expect string
		skip   int
		typ    TokenType
		want   bool
	}{
		{"", "", 0, TOK_NONE, false},
		{" ", " ", 1, TOK_WHITESPACE, true},
		{"  ", "  ", 2, TOK_WHITESPACE, true},
		{"   ", "   ", 3, TOK_WHITESPACE, true},
		{" (", " ", 1, TOK_WHITESPACE, true},
		{" )", " ", 1, TOK_WHITESPACE, true},
		{" ()", " ", 1, TOK_WHITESPACE, true},

		{"()", "(", 1, TOK_OPENPAREN, true},
		{"( )", "(", 1, TOK_OPENPAREN, true},
		{"(a)", "(", 1, TOK_OPENPAREN, true},
		{"(a )", "(", 1, TOK_OPENPAREN, true},
		{"( a)", "(", 1, TOK_OPENPAREN, true},
		{"( a )", "(", 1, TOK_OPENPAREN, true},
		{"( aa)", "(", 1, TOK_OPENPAREN, true},
		{"(a aa aaa) ", "(", 1, TOK_OPENPAREN, true},
		{`("hello" world)`, "(", 1, TOK_OPENPAREN, true},

		{"))", ")", 1, TOK_CLOSEPAREN, true},
		{") )", ")", 1, TOK_CLOSEPAREN, true},
		{")a)", ")", 1, TOK_CLOSEPAREN, true},
		{")a )", ")", 1, TOK_CLOSEPAREN, true},
		{") a)", ")", 1, TOK_CLOSEPAREN, true},
		{") a )", ")", 1, TOK_CLOSEPAREN, true},
		{") aa)", ")", 1, TOK_CLOSEPAREN, true},

		{`"hello" foo`, `"hello"`, 7, TOK_STRING, true},
		{`baz)`, `baz`, 3, TOK_SYMBOL, true},

		{";hello\n(lambda (x) (* x x))", `;hello`, 6, TOK_COMMENT, true},
	}

	td := TokenTypeDict()

	for _, t := range tests {
		tok, skip, ok := lex.lexOnce(t.s)
		assert.Equal(ok, t.want, "unexpectedly found a token in "+t.s)
		assert.Equal(tok.Typ, t.typ, fmt.Sprintf("Token type %s != %s (expected)", td(tok.Typ), td(t.typ)))
		assert.Equal(skip, t.skip, "skip count matches")
		if t.want {
			assert.True(ok, "should have found first token in `"+t.s+"`")
		}
		assert.Equal(tok.Token, t.expect, tok.Token+" did not match "+t.expect)
	}
}

func TestLex_(t *testing.T) {
	assert := assert.New(t)

	lex := NewLexer()

	tokens := lex.Lex("(foo bar baz)")

	types := []TokenType{TOK_OPENPAREN, TOK_SYMBOL, TOK_WHITESPACE, TOK_SYMBOL,
		TOK_WHITESPACE, TOK_SYMBOL, TOK_CLOSEPAREN}

	for i, t := range tokens {
		assert.Equal(t.Typ, types[i], "matching token types")
	}
	assert.Equal(len(tokens), 7, fmt.Sprintf("should be 7 tokens, have %d", len(tokens)))
}

/*
func TestLex(t *testing.T) {
    assert := assert.New(t)

    lex := NewLexer()

    tests := []string {" ;hello\n(lambda (x) (* x x))"}
    //tests := []string {`()`, `( )`, `(())`, ` ;hello\n(lambda (x) (* x x))`}

    for _, t := range tests {
        fmt.Println("-----------------")
        fmt.Println("Testing `", t, "`")
        printTokens(lex.Lex(t))
        fmt.Println("-----------------")
    }

}
*/

/*
func TestFOO(t *testing.T) {
    assert := assert.New(t)
    assert.Fail("FIX ME")
}
*/
