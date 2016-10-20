package san

import (
	"bytes"

	model "github.com/fgrehm/go-san/model"
	parser "github.com/fgrehm/go-san/parser"
)

// Parse parses a textual san model into a machine friendly structure
func Parse(src []byte) (*model.Model, error) {
	file, err := parser.New(src).Parse()
	if err != nil {
		return nil, err
	}
	return translateAstToModel(file), nil
}
