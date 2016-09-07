package sanast

import (
	token "github.com/fgrehm/go-san/token"
)

// EventsDefinition represents the set of events defined on the SAN file
type EventsDefinition struct {
	Token        token.Token
	Descriptions []*EventDescription
}

// EventDescription represents a single event description present on
// the events block
type EventDescription struct {
	Type token.Token // the type of event (local or synchronizing)
	Name token.Token // the name of the event
	Rate token.Token // the firing rate of the event
}
