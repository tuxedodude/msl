package lex

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

type TokenType int

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

type Token struct {
	Token string
	Typ   TokenType
	Loc   int
}

// Stringer for Token struct
func (t *Token) String() string {
	return fmt.Sprintf("{Loc:\t%d\tType:\t%s\tToken:\t%q\t}",
		t.Loc,
		TokenTypeDict()(t.Typ),
		t.Token)
}

// is this token a ( or )?
func (t *Token) isParen() bool {
	return t.Typ == TOK_OPENPAREN || t.Typ == TOK_CLOSEPAREN
}

// print a slice of tokens
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
