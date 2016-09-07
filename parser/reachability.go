package sanparser

import (
	"fmt"

	ast "github.com/fgrehm/go-san/ast"
	token "github.com/fgrehm/go-san/token"
)

func parseReachability(f *ast.File, p *parser, firstToken token.Token) error {
	defer un(trace(p, "ParseReachabilityDefinition"))

	var err error
	reachabilityDef := &ast.ReachabilityDefinition{
		Tokens: []token.Token{firstToken},
	}
	f.Reachability = reachabilityDef

	tok := p.scan()
	if firstToken.Type == token.PARTIAL {
		if tok.Type != token.REACHABILITY {
			return fmt.Errorf("Unexpected token found: %s. Expected to find 'reachability'", tok.String())
		}
		reachabilityDef.Tokens = append(reachabilityDef.Tokens, tok)
	} else {
		p.unscan()
	}

	tok = p.scan()
	if tok.Type != token.ASSIGN {
		return fmt.Errorf("Unexpected token found: %s. Expected an =", tok.String())
	}

	reachabilityDef.Expression, err = p.scanExpression()
	if err != nil {
		return err
	}

	return nil
}
