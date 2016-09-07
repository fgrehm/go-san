package sanparser

import (
	"fmt"

	ast "github.com/fgrehm/go-san/ast"
	token "github.com/fgrehm/go-san/token"
)

func parseEvents(f *ast.File, p *parser, eventsToken token.Token) error {
	defer un(trace(p, "parseEvents"))

	eventsDefinition := &ast.EventsDefinition{
		Token:        eventsToken,
		Descriptions: []*ast.EventDescription{},
	}
	f.Events = eventsDefinition

	for {
		tok := p.scan()
		if tok.Type == token.EOF {
			break
		}
		if !tok.Type.IsEventType() {
			if len(eventsDefinition.Descriptions) == 0 {
				return fmt.Errorf("Unexpected token found: %s. Expected an event type ('loc' or 'syn')", tok.String())
			}
			p.unscan()
			break
		}

		descriptionTrace := trace(p, "parseEventDescription")
		description := &ast.EventDescription{Type: tok}

		tok = p.scan()
		if tok.Type != token.IDENTIFIER {
			return fmt.Errorf("Unexpected token found: %s. Expected an identifier", tok.String())
		}
		description.Name = tok

		tok = p.scan()
		if tok.Type != token.LPAREN {
			return fmt.Errorf("Unexpected token found: %s. Expected a (", tok.String())
		}

		tok = p.scan()
		if tok.Type != token.IDENTIFIER {
			return fmt.Errorf("Unexpected token found: %s. Expected an identifier", tok.String())
		}
		description.Rate = tok

		tok = p.scan()
		if tok.Type != token.RPAREN {
			return fmt.Errorf("Unexpected token found: %s. Expected a )", tok.String())
		}
		tok = p.scan()
		if tok.Type != token.SEMICOLON {
			return fmt.Errorf("Unexpected token found: %s. Expected a ;", tok.String())
		}

		un(descriptionTrace)

		eventsDefinition.Descriptions = append(eventsDefinition.Descriptions, description)
	}

	return nil
}
