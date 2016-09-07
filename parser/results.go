package sanparser

import (
	"fmt"

	ast "github.com/fgrehm/go-san/ast"
	token "github.com/fgrehm/go-san/token"
)

func parseResults(f *ast.File, p *parser, resultsToken token.Token) error {
	defer un(trace(p, "parseResults"))

	var err error

	resultsDef := &ast.ResultsDefinition{
		Token:        resultsToken,
		Descriptions: []*ast.ResultDescription{},
	}
	f.Results = resultsDef

	for {
		tok := p.scan()
		if tok.Type == token.EOF {
			if len(resultsDef.Descriptions) == 0 {
				return fmt.Errorf("Expected to find a list of results")
			}
			break
		}
		if tok.Type.IsKeyword() {
			p.unscan()
			break
		}
		if tok.Type != token.IDENTIFIER {
			return fmt.Errorf("Unexpected token found: %s. Expected an identifier", tok.String())
		}

		resultTrace := trace(p, "parseResultDescription")
		result := &ast.ResultDescription{Label: tok}

		tok = p.scan()
		if tok.Type != token.ASSIGN {
			return fmt.Errorf("Unexpected token found: %s. Expected an =", tok.String())
		}

		result.Expression, err = p.scanExpression()
		if err != nil {
			return err
		}
		un(resultTrace)

		resultsDef.Descriptions = append(resultsDef.Descriptions, result)
	}

	return nil
}
