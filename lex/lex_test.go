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
				// note: if we can't pass the test, then "hello""hello" will make it thru the lexer.
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
