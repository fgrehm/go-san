package sanparser

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
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

type parsedEvent struct {
	line   int
	column int
	evType string
	name   string
	rate   string
}
type parsedEventsDefinition struct {
	line   int
	events []parsedEvent
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
// Invalid models

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
