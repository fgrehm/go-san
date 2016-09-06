package santoken

import (
	"fmt"
	"strconv"
)

// Token defines a single token which can be obtained via the Scanner
type Token struct {
	Type Type
	Pos  Pos
	Text string
}

// Type is the set of lexical tokens of the SAN model
type Type int

const (
	// Special tokens
	ILLEGAL Type = iota
	EOF
	COMMENT
	SEMICOLON

	literal_beg
	IDENTIFIER // literals
	NUMBER     // 12345
	FLOAT      // 123.45
	literal_end

	LPAREN // (
	RPAREN // )
	ASSIGN // =
	// TODO: Support +
	// TODO: Support -
	MULT // *
	// TODO: Support /
	AND   // &&
	EQUAL // ==

	keyword_beg
	IDENTIFIERS
	EVENTS
	PARTIAL
	REACHABILITY
	NETWORK
	CONTINUOUS
	LOC
	SYN
	AUT
	STT
	ST
	TO
	RESULTS
	keyword_end
)

var tokens = [...]string{
	ILLEGAL:   "ILLEGAL",
	EOF:       "EOF",
	COMMENT:   "COMMENT",
	SEMICOLON: "SEMICOLON",

	IDENTIFIER: "IDENTIFIER",
	NUMBER:     "NUMBER",
	FLOAT:      "FLOAT",

	LPAREN: "LPAREN",
	RPAREN: "RPAREN",
	ASSIGN: "ASSIGN",
	MULT:   "MULT",
	AND:    "AND",
	EQUAL:  "EQUAL",

	IDENTIFIERS:  "IDENTIFIERS",
	EVENTS:       "EVENTS",
	LOC:          "LOC",
	SYN:          "SYN",
	PARTIAL:      "PARTIAL",
	REACHABILITY: "REACHABILITY",
	NETWORK:      "NETWORK",
	CONTINUOUS:   "CONTINUOUS",
	AUT:          "AUT",
	STT:          "STT",
	ST:           "ST",
	TO:           "TO",
	RESULTS:      "RESULTS",
}

// IsLiteral returns true for tokens corresponding to basic type literals; it
// returns false otherwise.
func (t Type) IsLiteral() bool { return literal_beg < t && t < literal_end }

// IsKeyword returns true for tokens corresponding to keywords; it returns
// false otherwise.
func (t Type) IsKeyword() bool { return keyword_beg < t && t < keyword_end }

// String returns the string corresponding to the token tok.
func (t Type) String() string {
	s := ""
	if 0 <= t && t < Type(len(tokens)) {
		s = tokens[t]
	}
	if s == "" {
		s = "token(" + strconv.Itoa(int(t)) + ")"
	}
	return s
}

// String returns the token's literal text.
func (t Token) String() string {
	return fmt.Sprintf("%s %s %s", t.Pos.String(), t.Type.String(), t.Text)
}

// Value returns the properly typed value for this token. The type of
// the returned interface{} is guaranteed based on the Type field.
//
// This can only be called for literal types. If it is called for any other
// type, this will panic.
func (t Token) Value() interface{} {
	switch t.Type {
	case FLOAT:
		v, err := strconv.ParseFloat(t.Text, 64)
		if err != nil {
			panic(err)
		}

		return float64(v)
	case NUMBER:
		v, err := strconv.ParseInt(t.Text, 0, 64)
		if err != nil {
			panic(err)
		}

		return int64(v)
	case IDENTIFIER:
		return t.Text
	default:
		panic(fmt.Sprintf("unimplemented Value for type: %s", t.Type))
	}
}
