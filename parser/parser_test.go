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

		equals(t, m.expected.assignments, parsed.assignments)
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
