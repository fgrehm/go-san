package san

import (
	ast "github.com/fgrehm/go-san/ast"
	model "github.com/fgrehm/go-san/model"
	parser "github.com/fgrehm/go-san/parser"
	token "github.com/fgrehm/go-san/token"
)

type translator func(*model.Model, *ast.File)

var translators = []translator{
	translateIdentifiers,
	translateEvents,
	translateReachabilityInfo,
	translateNetwork,
	translateResults,
}

// NewParserFunc is a function that creates the parser used by san.Parse, meant
// to be overwriten on unit tests
var NewParserFunc = func(src []byte) parser.Parser {
	return parser.New(src)
}

// Parse parses a textual san model into a machine friendly structure
func Parse(src []byte) (*model.Model, error) {
	file, err := NewParserFunc(src).Parse()
	if err != nil {
		return nil, err
	}
	return translateAstToModel(file), nil
}

func translateAstToModel(file *ast.File) *model.Model {
	model := model.New()
	for _, t := range translators {
		t(model, file)
	}
	return model
}

func translateIdentifiers(m *model.Model, f *ast.File) {
	if f.Identifiers == nil {
		return
	}

	for _, assignment := range f.Identifiers.Assignments {
		m.AddIdentifier(&model.Identifier{
			Name:  assignment.Identifier.Text,
			Type:  assignment.Expression.Type(),
			Value: assignment.Expression.Value(),
		})
	}
}

func translateEvents(m *model.Model, f *ast.File) {
	if f.Events == nil {
		return
	}

	for _, event := range f.Events.Descriptions {
		eventType := "local"
		if event.Type.Type == token.SYN {
			eventType = "synchronizing"
		}

		m.AddEvent(&model.Event{
			Name: event.Name.Text,
			Type: eventType,
			Rate: event.Rate.Text,
		})
	}
}

func translateReachabilityInfo(m *model.Model, f *ast.File) {
	if f.Reachability == nil {
		return
	}

	m.Reachability.Partial = f.Reachability.Tokens[0].Type == token.PARTIAL
	m.Reachability.Expression = f.Reachability.Expression.Text()
}

func translateNetwork(m *model.Model, f *ast.File) {
	if f.Network == nil {
		return
	}

	m.Network.Name = f.Network.Name.Text
	m.Network.Type = f.Network.Type.Text

	for _, automaton := range f.Network.Automata {
		translateAutomaton(m.Network, automaton)
	}
}

func translateAutomaton(n *model.Network, a *ast.AutomatonDescription) {
	aut := &model.Automaton{
		Name:        a.Name.Text,
		Transitions: model.Transitions{},
	}
	for _, transition := range a.Transitions {
		translateTransition(aut, transition)
	}
	n.AddAutomaton(aut)
}

func translateTransition(a *model.Automaton, t *ast.AutomatonTransition) {
	events := []string{}
	for _, e := range t.Events {
		events = append(events, e.Text)
	}
	a.AddTransition(&model.Transition{
		From:   t.From.Text,
		To:     t.To.Text,
		Events: events,
	})
}

func translateResults(m *model.Model, f *ast.File) {
	if f.Results == nil {
		return
	}

	for _, desc := range f.Results.Descriptions {
		m.AddResult(&model.Result{
			Label:      desc.Label.Text,
			Expression: desc.Expression.Text(),
		})
	}
}
