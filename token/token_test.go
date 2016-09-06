package santoken

import (
	"testing"
)

func TestTypeString(t *testing.T) {
	var tokens = []struct {
		tt  Type
		str string
	}{
		{ILLEGAL, "ILLEGAL"},
		{EOF, "EOF"},
		{COMMENT, "COMMENT"},
		{SEMICOLON, "SEMICOLON"},

		{IDENTIFIER, "IDENTIFIER"},
		{NUMBER, "NUMBER"},
		{FLOAT, "FLOAT"},

		{LPAREN, "LPAREN"},
		{RPAREN, "RPAREN"},
		{ASSIGN, "ASSIGN"},
		{MULT, "MULT"},
		{AND, "AND"},
		{EQUAL, "EQUAL"},

		{IDENTIFIERS, "IDENTIFIERS"},
		{EVENTS, "EVENTS"},
		{PARTIAL, "PARTIAL"},
		{REACHABILITY, "REACHABILITY"},
		{NETWORK, "NETWORK"},
		{CONTINUOUS, "CONTINUOUS"},
		{LOC, "LOC"},
		{SYN, "SYN"},
		{AUT, "AUT"},
		{STT, "STT"},
		{ST, "ST"},
		{TO, "TO"},
		{RESULTS, "RESULTS"},
	}

	for _, token := range tokens {
		if token.tt.String() != token.str {
			t.Errorf("want: %q got:%q\n", token.str, token.tt)
		}
	}
}
