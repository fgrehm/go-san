package sanparser

import (
	"fmt"

	ast "github.com/fgrehm/go-san/ast"
	token "github.com/fgrehm/go-san/token"
)

func parseNetwork(f *ast.File, p *parser, networkToken token.Token) error {
	defer un(trace(p, "parseNetwork"))

	networkDef := &ast.NetworkDefinition{
		Token:    networkToken,
		Automata: []*ast.AutomatonDescription{},
	}
	f.Network = networkDef

	tok := p.scan()
	if tok.Type != token.IDENTIFIER {
		return fmt.Errorf("Unexpected token found: %s. Expected an identifier", tok.String())
	}
	networkDef.Name = tok
	tok = p.scan()
	if tok.Type != token.LPAREN {
		return fmt.Errorf("Unexpected token found: %s. Expected a (", tok.String())
	}
	tok = p.scan()
	if tok.Type != token.CONTINUOUS {
		return fmt.Errorf("Unexpected token found: %s. Expected to find the 'continous' keyword", tok.String())
	}
	networkDef.Type = tok

	tok = p.scan()
	if tok.Type != token.RPAREN {
		return fmt.Errorf("Unexpected token found: %s. Expected a )", tok.String())
	}

	for {
		descriptionTrace := trace(p, "parseAutomatonDescription")
		tok = p.scan()
		if tok.Type != token.AUT {
			if len(networkDef.Automata) == 0 {
				return fmt.Errorf("Unexpected token found: %s. Expected to find the 'aut' keyword", tok.String())
			}
			if tok.Type != token.EOF {
				p.unscan()
			}
			break
		}
		// tok is the aut keyword
		automatonDesc, err := parseAutomatonDescription(p, tok)
		if err != nil {
			return err
		}
		networkDef.Automata = append(networkDef.Automata, automatonDesc)
		un(descriptionTrace)
	}

	return nil
}

func parseAutomatonDescription(p *parser, autToken token.Token) (*ast.AutomatonDescription, error) {
	automatonDesc := &ast.AutomatonDescription{
		Token:       autToken,
		Transitions: []*ast.AutomatonTransition{},
	}

	tok := p.scan()
	if tok.Type != token.IDENTIFIER {
		return nil, fmt.Errorf("Unexpected token found: %s. Expected to find an identifier", tok.String())
	}
	automatonDesc.Name = tok

	for {
		tok = p.scan()
		if tok.Type != token.STT {
			if len(automatonDesc.Transitions) == 0 {
				return nil, fmt.Errorf("Unexpected EOF. Expected to find the 'stt' keyword")
			}
			p.unscan()
			break
		}

		transitionsTrace := trace(p, "parseAutomatonTransitions")
		transitions, err := parseAutomatonTransitions(p, tok)
		if err != nil {
			return nil, err
		}
		automatonDesc.Transitions = append(automatonDesc.Transitions, transitions...)
		un(transitionsTrace)
	}

	return automatonDesc, nil
}

func parseAutomatonTransitions(p *parser, sttToken token.Token) ([]*ast.AutomatonTransition, error) {
	from := p.scan()
	if from.Type != token.IDENTIFIER {
		return nil, fmt.Errorf("Unexpected token found: %s. Expected to find an identifier", from.String())
	}

	transitions := []*ast.AutomatonTransition{}

	for {
		tok := p.scan()
		if tok.Type != token.TO {
			p.unscan()
			break
		}

		transition := &ast.AutomatonTransition{
			From:   from,
			Events: []token.Token{},
		}

		tok = p.scan()
		if tok.Type != token.LPAREN {
			return nil, fmt.Errorf("Unexpected token found: %s. Expected a (", tok.String())
		}
		tok = p.scan()
		if tok.Type != token.IDENTIFIER {
			return nil, fmt.Errorf("Unexpected token found: %s. Expected an identifier", tok.String())
		}
		transition.To = tok

		tok = p.scan()
		if tok.Type != token.RPAREN {
			return nil, fmt.Errorf("Unexpected token found: %s. Expected a )", tok.String())
		}

		tok = p.scan()
		if tok.Type != token.IDENTIFIER {
			return nil, fmt.Errorf("Unexpected token found: %s. Expected an identifier", tok.String())
		}
		transition.Events = append(transition.Events, tok)
		for {
			tok = p.scan()
			if tok.Type != token.IDENTIFIER {
				p.unscan()
				break
			}
			transition.Events = append(transition.Events, tok)
		}

		transitions = append(transitions, transition)
	}

	return transitions, nil
}
