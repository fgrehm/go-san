package santoken

import (
	"reflect"
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
		tt      Type
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

func TestLiterals(t *testing.T) {
	var tokens = []struct {
		tt      Type
		literal bool
	}{
		{ILLEGAL, false},
		{EOF, false},
		{COMMENT, false},
		{SEMICOLON, false},

		{IDENTIFIER, true},
		{NUMBER, true},
		{FLOAT, true},

		{LPAREN, false},
		{RPAREN, false},
		{ASSIGN, false},
		{MULT, false},
		{AND, false},
		{EQUAL, false},

		{IDENTIFIERS, false},
		{EVENTS, false},
		{PARTIAL, false},
		{REACHABILITY, false},
		{NETWORK, false},
		{CONTINUOUS, false},
		{LOC, false},
		{SYN, false},
		{AUT, false},
		{STT, false},
		{ST, false},
		{TO, false},
		{RESULTS, false},
	}

	for _, token := range tokens {
		if token.tt.IsLiteral() != token.literal {
			t.Errorf("want: %v got: %v for %s\n", token.literal, token.tt.IsLiteral(), token.tt.String())
		}
	}
}

func TestTokenValue(t *testing.T) {
	var tokens = []struct {
		tt Token
		v  interface{}
	}{
		{Token{Type: FLOAT, Text: `3.14`}, float64(3.14)},
		{Token{Type: NUMBER, Text: `42`}, int64(42)},
		{Token{Type: IDENTIFIER, Text: `foo`}, "foo"},
	}

	for _, token := range tokens {
		if val := token.tt.Value(); !reflect.DeepEqual(val, token.v) {
			t.Errorf("want: %v got:%v\n", token.v, val)
		}
	}
}
