package sanast

import (
	token "github.com/fgrehm/go-san/token"
)

// File represents a single SAN file
type File struct {
	Identifiers *IdentifiersDefinition
	Events      *EventsDefinition
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
