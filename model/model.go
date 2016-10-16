package sanmodel

import (
	"bytes"
	"encoding/gob"
)

// Model represents a model that has been parsed from a .san file
type Model struct {
	Identifiers  Identifiers   `json:"identifiers"`
	Events       Events        `json:"events"`
	Reachability *Reachability `json:"reachability"`
	Network      *Network      `json:"network"`
	Results      Results       `json:"results"`
}

// Identifier represents a single identifier that has been parsed from the
// `identifiers` block
type Identifier struct {
	Name  string      `json:"name"`
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

// Identifiers represent a collection of identifiers present on the `identifiers` block
type Identifiers []*Identifier

// Event represents a single event that has been parsed from the `events` block
type Event struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Rate string `json:"rate"`
}

// Events represent a collection of events present on the `events` block
type Events []*Event

// Reachability represents the reachability information about the model network
type Reachability struct {
	Partial    bool   `json:"partial"`
	Expression string `json:"expression"`
}

// Network aggregates automata information from the `network` block
type Network struct {
	Name     string   `json:"name"`
	Type     string   `json:"type"`
	Automata Automata `json:"automata"`
}

// Automaton represents a single automaton present on the `network` block
type Automaton struct {
	Name        string      `json:"name"`
	Transitions Transitions `json:"transitions"`
}

// Transitions represent a collection of automaton transitions
type Transitions []*Transition

// Automata represents a collection of automatons
type Automata []*Automaton

// Transition represents a single automaton transition
type Transition struct {
	From   string           `json:"from"`
	To     string           `json:"to"`
	Events TransitionEvents `json:"events"`
}

// TransitionEvent represents a single automaton transition event
type TransitionEvent struct {
	EventName   string `json:"name"`
	Probability string `json:"probability"`
}

// TransitionEvents represents a collection of transition events
type TransitionEvents []*TransitionEvent

// Result represents a single result present on the `results` block
type Result struct {
	Label      string `json:"label"`
	Expression string `json:"expression"`
}

// Results represents a collection of results present on the `results` block
type Results []*Result

// New instantiates a new model struct
func New() *Model {
	return &Model{
		Identifiers:  Identifiers{},
		Events:       Events{},
		Reachability: &Reachability{},
		Network: &Network{
			Automata: Automata{},
		},
		Results: Results{},
	}
}

// Copy copies over a model to another variable so that it can be manipulated
// without side effects
func (m *Model) Copy(dest *Model) error {
	var mod bytes.Buffer
	enc := gob.NewEncoder(&mod)
	dec := gob.NewDecoder(&mod)

	err := enc.Encode(m)
	if err != nil {
		return err
	}

	return dec.Decode(&dest)
}

// AddIdentifier adds an Identifier to the model
func (m *Model) AddIdentifier(i *Identifier) {
	m.Identifiers = append(m.Identifiers, i)
}

// AddEvent adds an Event to the model
func (m *Model) AddEvent(e *Event) {
	m.Events = append(m.Events, e)
}

// AddResult adds a Result to the model
func (m *Model) AddResult(r *Result) {
	m.Results = append(m.Results, r)
}

// AddAutomaton adds an Automaton to a network
func (n *Network) AddAutomaton(a *Automaton) {
	n.Automata = append(n.Automata, a)
}

// AddTransition adds a Transition to the automaton
func (a *Automaton) AddTransition(t *Transition) {
	a.Transitions = append(a.Transitions, t)
}
