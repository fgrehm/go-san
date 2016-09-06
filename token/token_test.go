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

func TestKeywords(t *testing.T) {
	var tokens = []struct {
		tt  Type
		keyword bool
	}{
		{ILLEGAL, false},
		{EOF, false},
		{COMMENT, false},
		{SEMICOLON, false},

		{IDENTIFIER, false},
		{NUMBER, false},
		{FLOAT, false},

		{LPAREN, false},
		{RPAREN, false},
		{ASSIGN, false},
		{MULT, false},
		{AND, false},
		{EQUAL, false},

		{IDENTIFIERS, true},
		{EVENTS, true},
		{PARTIAL, true},
		{REACHABILITY, true},
		{NETWORK, true},
		{CONTINUOUS, true},
		{LOC, true},
		{SYN, true},
		{AUT, true},
		{STT, true},
		{ST, true},
		{TO, true},
		{RESULTS, true},
	}

	for _, token := range tokens {
		if token.tt.IsKeyword() != token.keyword {
			t.Errorf("want: %v got: %v for %s\n", token.keyword, token.tt.IsKeyword(), token.tt.String())
		}
	}
}
