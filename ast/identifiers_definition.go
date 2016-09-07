package sanast

import (
	token "github.com/fgrehm/go-san/token"
)

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

