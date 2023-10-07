package lex

import (
	"testing"
    "github.com/stretchr/testify/assert"
    "regexp"
)

func TestRegex(t *testing.T) {

    assert := assert.New(t)

    re := regexp.MustCompile(patWhiteSpace)

    var tests = []struct {
        s string
        msg string
        want bool
    }{
        {" ", "one space", true},
        {"", "Empty string", false},
        {"  ", "two spaces", true},
        {" \t\n \r ", "mixed spaces", true},
        {"a    ", "leading non-space char", false},
    }
    
    for _, test := range tests {
        if test.want {
            assert.True(re.MatchString(test.s), "Whitespace should succeed on" + test.msg)
        } else {
            assert.False(re.MatchString(test.s), "Whitespace should fail on", test.msg)
        }
    }

}

