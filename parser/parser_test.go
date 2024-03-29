package sanparser

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	token "github.com/fgrehm/go-san/token"
)

// ----------------------------------------------------------------------------
// Identifiers block

type parsedIdentifierAssignment struct {
	line       int
	column     int
	identifier string
	expType    string
	value      interface{}
}
type parsedIdentifiersDefinition struct {
	line        int
	assignments []parsedIdentifierAssignment
}

func TestParseIdentifiersDefinition(t *testing.T) {
	var models = []struct {
		src      string
		expected parsedIdentifiersDefinition
	}{
		{
			"identifiers\nrate   = 3;\n  r_2=4;",
			parsedIdentifiersDefinition{
				line: 1,
				assignments: []parsedIdentifierAssignment{
					{line: 2, column: 1, identifier: "rate", expType: "constant", value: int64(3)},
					{line: 3, column: 3, identifier: "r_2", expType: "constant", value: int64(4)},
				},
			},
		},
		{
			"// Foo \nidentifiers\nF1 = (st Client == Working) * 1; r_2= 4/* */;",
			parsedIdentifiersDefinition{
				line: 2,
				assignments: []parsedIdentifierAssignment{
					{line: 3, column: 1, identifier: "F1", expType: "expression", value: "( st Client == Working ) * 1"},
					{line: 3, column: 34, identifier: "r_2", expType: "constant", value: int64(4)},
				},
			},
		},
	}

	for _, m := range models {
		file, err := Parse([]byte(m.src))
		if err != nil {
			t.Error(err)
			continue
		}

		parsed := parsedIdentifiersDefinition{
			line:        file.Identifiers.Token.Pos.Line,
			assignments: []parsedIdentifierAssignment{},
		}
		for _, identifierAssignment := range file.Identifiers.Assignments {
			parsed.assignments = append(parsed.assignments, parsedIdentifierAssignment{
				line:       identifierAssignment.Identifier.Pos.Line,
				column:     identifierAssignment.Identifier.Pos.Column,
				identifier: identifierAssignment.Identifier.Text,
				expType:    identifierAssignment.Expression.Type(),
				value:      identifierAssignment.Expression.Value(),
			})
		}

		equals(t, m.expected, parsed)
	}
}

func TestParseIdentifiersDefinition_Error(t *testing.T) {
	var models = []string{
		"identifiers f1",
		"identifiers f1 = ",
		"identifiers f1 = ;",
		"identifiers f1 = 1",
		"identifiers\nf1 = 1; a=; t = 3;",
	}

	for _, m := range models {
		_, err := Parse([]byte(m))
		if err == nil {
			t.Errorf("Expected to error with %q but did not", m)
		}
	}
}

// ----------------------------------------------------------------------------
// Events block

type parsedEventsDefinition struct {
	line   int
	events []parsedEvent
}

type parsedEvent struct {
	line   int
	column int
	evType string
	name   string
	rate   string
}

func TestParseEventsDefinition(t *testing.T) {
	var models = []struct {
		src      string
		expected parsedEventsDefinition
	}{
		{
			"events\nloc foo (bar);\n syn john (doe);",
			parsedEventsDefinition{
				line: 1,
				events: []parsedEvent{
					{line: 2, column: 1, evType: "loc", name: "foo", rate: "bar"},
					{line: 3, column: 2, evType: "syn", name: "john", rate: "doe"},
				},
			},
		},
		{
			"/* */events\n // loc bla; \nloc foo (bar); loc l_req (r_req); syn john (doe);",
			parsedEventsDefinition{
				line: 1,
				events: []parsedEvent{
					{line: 3, column: 1, evType: "loc", name: "foo", rate: "bar"},
					{line: 3, column: 16, evType: "loc", name: "l_req", rate: "r_req"},
					{line: 3, column: 35, evType: "syn", name: "john", rate: "doe"},
				},
			},
		},
	}

	for _, m := range models {
		file, err := Parse([]byte(m.src))
		if err != nil {
			t.Error(err)
			continue
		}

		parsed := parsedEventsDefinition{
			line:   file.Events.Token.Pos.Line,
			events: []parsedEvent{},
		}
		for _, eventDescription := range file.Events.Descriptions {
			parsed.events = append(parsed.events, parsedEvent{
				line:   eventDescription.Type.Pos.Line,
				column: eventDescription.Type.Pos.Column,
				name:   eventDescription.Name.Text,
				evType: eventDescription.Type.Text,
				rate:   eventDescription.Rate.Text,
			})
		}

		equals(t, m.expected, parsed)
	}
}

func TestParseEventsDefinition_Error(t *testing.T) {
	var models = []string{
		"events f1",
		"events f1;",
		"events loc ",
		"events loc ;",
		"events loc foo ;",
		"events loc foo ();",
		"events syn ",
		"events syn ;",
		"events syn foo ;",
		"events syn foo ();",
	}

	for _, m := range models {
		_, err := Parse([]byte(m))
		if err == nil {
			t.Errorf("Expected to error with %q but did not", m)
		}
	}
}

// ----------------------------------------------------------------------------
// Reachability block

func TestParseReachabilityInfo(t *testing.T) {
	var testData = []struct {
		src             string
		expectedPartial bool
		expectedExp     interface{}
	}{
		{
			"partial reachability = ((st Client == Idle) && (st Server == Idle));",
			true, "( ( st Client == Idle ) && ( st Server == Idle ) )",
		},
		{
			"reachability = (( st Client == Idle) && (st Server == Idle));",
			false, "( ( st Client == Idle ) && ( st Server == Idle ) )",
		},
		{
			"reachability = 1;",
			false, int64(1),
		},
	}

	for _, a := range testData {
		file, err := Parse(([]byte(a.src)))
		if err != nil {
			t.Error(err)
			continue
		}
		info := file.Reachability

		equals(t, a.expectedPartial, info.Tokens[0].Type == token.PARTIAL)
		equals(t, a.expectedExp, info.Expression.Value())
	}
}

func TestParseReachabilityInfo_Error(t *testing.T) {
	var models = []string{
		"partial ;",
		"partial reachability;",
		"partial reachability = ;",
		"reachability ;",
		"reachability = ;",
		"reachability events ;",
	}

	for _, m := range models {
		_, err := Parse([]byte(m))
		if err == nil {
			t.Errorf("Expected to error with %q but did not", m)
		}
	}
}

// ----------------------------------------------------------------------------
// Network block

type parsedNetworkDefinition struct {
	line     int
	name     string
	netType  string
	automata []parsedAutomaton
}

type parsedAutomaton struct {
	line        int
	column      int
	name        string
	transitions []parsedAutomatonTransition
}

type parsedAutomatonTransition struct {
	from   string
	to     string
	events []string
}

func TestParseNetworkDefinition(t *testing.T) {
	src := `network ClientServer (continuous)
aut Client
  stt A to (B) s_1
  stt B to (C) s_2
  stt C to (B) s_3(p_1)
        to (A) s_4(p_2) s_5(p_3)
aut Server stt D to (e) s_6`
	expected := parsedNetworkDefinition{
		line:    1,
		name:    "ClientServer",
		netType: "continuous",
		automata: []parsedAutomaton{
			{
				line:   2,
				column: 1,
				name:   "Client",
				transitions: []parsedAutomatonTransition{
					{from: "A", to: "B", events: []string{"s_1|"}},
					{from: "B", to: "C", events: []string{"s_2|"}},
					{from: "C", to: "B", events: []string{"s_3|p_1"}},
					{from: "C", to: "A", events: []string{"s_4|p_2", "s_5|p_3"}},
				},
			},
			{
				line:   7,
				column: 1,
				name:   "Server",
				transitions: []parsedAutomatonTransition{
					{from: "D", to: "e", events: []string{"s_6|"}},
				},
			},
		},
	}

	file, err := Parse([]byte(src))
	if err != nil {
		t.Fatal(err)
	}

	parsed := parsedNetworkDefinition{
		line:     file.Network.Token.Pos.Line,
		name:     file.Network.Name.Text,
		netType:  file.Network.Type.Text,
		automata: []parsedAutomaton{},
	}
	for _, automatonDescription := range file.Network.Automata {
		automaton := parsedAutomaton{
			line:        automatonDescription.Token.Pos.Line,
			column:      automatonDescription.Token.Pos.Column,
			name:        automatonDescription.Name.Text,
			transitions: []parsedAutomatonTransition{},
		}
		for _, automatonTransition := range automatonDescription.Transitions {
			events := []string{}
			for _, event := range automatonTransition.Events {
				events = append(events, fmt.Sprintf("%s|%s", event.EventName.Text, event.Probability.Text))
			}
			automaton.transitions = append(automaton.transitions, parsedAutomatonTransition{
				from:   automatonTransition.From.Text,
				to:     automatonTransition.To.Text,
				events: events,
			})
		}
		parsed.automata = append(parsed.automata, automaton)
	}
	equals(t, expected, parsed)
}

func TestParseNetworkDefinition_Error(t *testing.T) {
	var models = []string{
		"network",
		"network Foo",
		"network Foo\naut",
		"network Foo (continous) aut",
	}

	for _, m := range models {
		_, err := Parse([]byte(m))
		if err == nil {
			t.Errorf("Expected to error with %q but did not", m)
		}
	}
}

// ----------------------------------------------------------------------------
// Results block

type parsedResultsDefinition struct {
	line         int
	column       int
	descriptions []parsedResult
}

type parsedResult struct {
	line       int
	column     int
	label      string
	expression string
}

func TestParseResultsDefinition(t *testing.T) {
	var models = []struct {
		src      string
		expected parsedResultsDefinition
	}{
		{
			"results\nA_b = (st Foo == bar)\n && (st bla == foo); a = st Bla == state;",
			parsedResultsDefinition{
				line: 1,
				descriptions: []parsedResult{
					{line: 2, column: 1, label: "A_b", expression: "( st Foo == bar ) && ( st bla == foo )"},
					{line: 3, column: 22, label: "a", expression: "st Bla == state"},
				},
			},
		},
		{
			"results//A_b = (st Foo == bar)\n //&& (st bla == foo);\n a = st Bla == state\n;",
			parsedResultsDefinition{
				line: 1,
				descriptions: []parsedResult{
					{line: 3, column: 2, label: "a", expression: "st Bla == state"},
				},
			},
		},
	}

	for _, m := range models {
		file, err := Parse([]byte(m.src))
		if err != nil {
			t.Error(err)
			continue
		}

		parsed := parsedResultsDefinition{
			line:         file.Results.Token.Pos.Line,
			descriptions: []parsedResult{},
		}
		for _, resultDescription := range file.Results.Descriptions {
			parsed.descriptions = append(parsed.descriptions, parsedResult{
				line:       resultDescription.Label.Pos.Line,
				column:     resultDescription.Label.Pos.Column,
				label:      resultDescription.Label.Text,
				expression: resultDescription.Expression.Text(),
			})
		}

		equals(t, m.expected, parsed)
	}
}

func TestParseResultsDefinition_Error(t *testing.T) {
	var models = []string{
		"results",
		"results ;",
		"results a ;",
		"results a = ;",
		"results a = (st a & 2;",
	}

	for _, m := range models {
		_, err := Parse([]byte(m))
		if err == nil {
			t.Errorf("Expected to error with %q but did not", m)
		}
	}
}

// ----------------------------------------------------------------------------
// Bad models

func TestErrorsWhenUnexpectedTokenFoundAtRoot(t *testing.T) {
	_, err := Parse([]byte("ident\n"))
	if err == nil {
		t.Error("Expected to error but did not")
	}
	_, err = Parse([]byte("// Foo\naa"))
	if err == nil {
		t.Error("Expected to error but did not")
	}
}

// ----------------------------------------------------------------------------
// Utilities

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}
