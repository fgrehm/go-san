package sanparser

import (
	"fmt"

	ast "github.com/fgrehm/go-san/ast"
	token "github.com/fgrehm/go-san/token"
)

func parseIdentifiers(f *ast.File, p *parser, identifiersToken token.Token) error {
	defer un(trace(p, "parseIdentifiers"))

	var err error
	idDef := &ast.IdentifiersDefinition{
		Token:       identifiersToken,
		Assignments: []*ast.IdentifierAssignment{},
	}
	f.Identifiers = idDef

	for {
		tok := p.scan()
		if tok.Type == token.EOF {
			break
		}
		if tok.Type.IsKeyword() {
			p.unscan()
			break
		}
		if tok.Type != token.IDENTIFIER {
			return p.err(tok.Pos, fmt.Errorf("Unexpected token found: %q. Expected an identifier", tok.Text))
		}

		assignmentTrace := trace(p, "ParseIdentifier")
		assignment := &ast.IdentifierAssignment{Identifier: tok}

		tok = p.scan()
		if tok.Type != token.ASSIGN {
			return p.err(tok.Pos, fmt.Errorf("Unexpected token found: %q. Expected an =", tok.Text))
		}

		assignment.Expression, err = p.scanExpression()
		if err != nil {
			return err
		}
		un(assignmentTrace)

		idDef.Assignments = append(idDef.Assignments, assignment)
	}

	return nil
}
