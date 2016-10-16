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
	// ILLEGAL represents an invalid token
	ILLEGAL Type = iota
	// EOF represents the end of the file
	EOF
	// COMMENT represents a comment block
	COMMENT
	// SEMICOLON represents a semicolon
	SEMICOLON

	// literals
	literalBeg
	// IDENTIFIER represents an identifier
	IDENTIFIER
	// NUMBER represents an integer number
	NUMBER
	// FLOAT represents a float number
	FLOAT
	literalEnd

	// LPAREN represents a left parenthesis
	LPAREN
	// RPAREN represents a right parenthesis
	RPAREN
	// ASSIGN represents an equal sign used on assignments
	ASSIGN
	// SUM represents a sum
	SUM
	// SUB represents a subtraction
	SUB
	// MULT represents a multiplication
	MULT
	// DIV represents a division
	DIV

	// AND represents the && operator
	AND
	// EQUAL represents the == operator
	EQUAL
	// NEQUAL represents the != operator
	NEQUAL

	keywordBeg
	// IDENTIFIERS represents the identifiers keyword
	IDENTIFIERS
	// EVENTS represents the events keyword
	EVENTS
	// PARTIAL represents the partial keyword
	PARTIAL
	// REACHABILITY represents the reachability keyword
	REACHABILITY
	// NETWORK represents the network keyword
	NETWORK
	// CONTINUOUS represents the continuous keyword
	CONTINUOUS
	eventTypeBeg
	// LOC represents the loc keyword
	LOC
	// SYN represents the syn keyword
	SYN
	eventTypeEnd
	// AUT represents the aut keyword
	AUT
	// STT represents the stt keyword
	STT
	// ST represents the st keyword
	ST
	// TO represents the to keyword
	TO
	// RESULTS represents the results keyword
	RESULTS
	keywordEnd
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
	SUM:    "SUM",
	SUB:    "SUB",
	MULT:   "MULT",
	DIV:    "DIV",
	AND:    "AND",
	EQUAL:  "EQUAL",
	NEQUAL: "NEQUAL",

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
func (t Type) IsLiteral() bool { return literalBeg < t && t < literalEnd }

// IsEventType returns true for tokens corresponding to event types (currently
// loc or syn)
func (t Type) IsEventType() bool { return eventTypeBeg < t && t < eventTypeEnd }

// IsKeyword returns true for tokens corresponding to keywords; it returns
// false otherwise.
func (t Type) IsKeyword() bool { return keywordBeg < t && t < keywordEnd }

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
