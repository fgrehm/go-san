package san

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	model "github.com/fgrehm/go-san/model"
)

type formatter func(*model.Model, *bytes.Buffer)

var formatters = []formatter{
	formatIdentifiers,
	formatEvents,
	formatReachability,
	formatNetwork,
	formatResults,
}

func formatIdentifiers(m *model.Model, buf *bytes.Buffer) {
	buf.WriteString("identifiers\n")
	for _, ident := range m.Identifiers {
		if ident.Type == "expression" {
			buf.WriteString(fmt.Sprintf("  %s = %s;\n", ident.Name, ident.Value))
		} else {
			if val, ok := ident.Value.(float64); ok {
				buf.WriteString(fmt.Sprintf("  %s = %f;\n", ident.Name, val))
			} else if val, ok := ident.Value.(int64); ok {
				buf.WriteString(fmt.Sprintf("  %s = %d;\n", ident.Name, val))
			}
		}
	}
}

func formatEvents(m *model.Model, buf *bytes.Buffer) {
	buf.WriteString("events\n")
	for _, event := range m.Events {
		if event.Type == "local" {
			buf.WriteString(fmt.Sprintf("  loc %s (%s);\n", event.Name, event.Rate))
		} else {
			buf.WriteString(fmt.Sprintf("  syn %s (%s);\n", event.Name, event.Rate))
		}
	}
}

func formatReachability(m *model.Model, buf *bytes.Buffer) {
	reachability := m.Reachability
	if reachability.Partial {
		buf.WriteString("partial ")
	}
	buf.WriteString(fmt.Sprintf("reachability = %s;\n", reachability.Expression))
}

func formatNetwork(m *model.Model, buf *bytes.Buffer) {
	network := m.Network

	buf.WriteString(fmt.Sprintf("network %s (%s)\n", network.Name, network.Type))
	for _, aut := range network.Automata {
		buf.WriteString(fmt.Sprintf("  aut %s\n", aut.Name))

		for _, state := range extractStates(aut.Transitions) {
			buf.WriteString(fmt.Sprintf("    stt %s\n", state))

			for _, transition := range aut.Transitions {
				if transition.From == state {
					events := []string{}
					for _, e := range transition.Events {
						event := e.EventName
						if e.Probability != "" {
							event += fmt.Sprintf("(%s)", e.Probability)
						}
						events = append(events, event)
					}
					buf.WriteString(fmt.Sprintf("      to (%s) %s\n", transition.To, strings.Join(events, " ")))
				}
			}
		}
	}
}

func extractStates(t model.Transitions) []string {
	statesMap := map[string]bool{}
	for _, transition := range t {
		statesMap[transition.From] = true
		statesMap[transition.To] = true
	}

	states := []string{}
	for state := range statesMap {
		states = append(states, state)
	}
	sort.Strings(states)
	return states
}

func formatResults(m *model.Model, buf *bytes.Buffer) {
	buf.WriteString("results\n")
	for _, res := range m.Results {
		buf.WriteString(fmt.Sprintf("  %s = %s;\n", res.Label, res.Expression))
	}
}
