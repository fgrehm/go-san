package sanast

import (
	token "github.com/fgrehm/go-san/token"
)

// ResultsDefinition represents the set of results defined on the SAN
// file
type ResultsDefinition struct {
	Token        token.Token
	Descriptions []*ResultDescription
}

// ResultDescription represents a single result description present on
// the results block
type ResultDescription struct {
	Label      token.Token // the result name itself
	Expression *Expression // the expression that represents the result
}
