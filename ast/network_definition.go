package sanast

import (
	token "github.com/fgrehm/go-san/token"
)

// NetworkDefinition represents the network information defined on the SAN file
// and is composed of a set of automata
type NetworkDefinition struct {
	Token    token.Token
	Name     token.Token
	Type     token.Token
	Automata []*AutomatonDescription
}

// AutomatonDescription represents a single automaton definition present on
// the network block
type AutomatonDescription struct {
	Token       token.Token
	Name        token.Token
	Transitions []*AutomatonTransition
}

// AutomatonTransition represents a single automaton transition present on
// the automaton block inside the network block
type AutomatonTransition struct {
	From   token.Token
	To     token.Token
	Events []*TransitionEventDescription
}

// TransitionEventDescription represents a single automaton transition event
// present on the automaton block inside the network block
type TransitionEventDescription struct {
	EventName   token.Token
	Probability token.Token
}
