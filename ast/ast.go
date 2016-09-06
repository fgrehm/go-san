package sanast

import (
	"strings"

	token "github.com/fgrehm/go-san/token"
)

// File represents a single SAN file
type File struct {
	Identifiers *IdentifiersDefinition
}

// IdentifiersDefinition represents the set of identifiers defined on the SAN
// file
type IdentifiersDefinition struct {
	Token       token.Token
	Assignments []*IdentifierAssignment
}

// IdentifierAssignment represents a single identifier definition present on
// the identifiers block
type IdentifierAssignment struct {
	Identifier token.Token // the identifier name itself
	Expression *Expression // the value to be assigned to the identifier
}

// Expression represents an expression used on the identifiers or results
// definitions
type Expression struct {
	Tokens []token.Token
}

// Value returns the properly typed value for this expression. The type of
// the returned interface{} is guaranteed based on the tokens found.
func (e *Expression) Value() interface{} {
	if len(e.Tokens) == 1 && e.Tokens[0].Type.IsLiteral() {
		return e.Tokens[0].Value()
	}
	return e.Text()
}

// Type returns the type of an expression, being a constant if a single litral
// is present or an expression if multiple tokens make up for the expression
func (e *Expression) Type() string {
	if len(e.Tokens) == 1 && e.Tokens[0].Type.IsLiteral() {
		return "constant"
	}
	return "expression"
}

// Text returns the list of tokens separated by spaces
func (e *Expression) Text() string {
	text := []string{}
	for _, t := range e.Tokens {
		text = append(text, t.Text)
	}
	return strings.Join(text, " ")
}

// Comment node represents a single //, # style or /*- style commment
type Comment struct {
	Start token.Pos // position of / or #
	Text  string
}

// CommentGroup node represents a sequence of comments with no other tokens and
// no empty lines between.
type CommentGroup struct {
	List []*Comment // len(List) > 0
}
