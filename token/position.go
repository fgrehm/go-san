package santoken

import "fmt"

// Pos describes an arbitrary source position
// including the file, line, and column location.
// A Position is valid if the line number is > 0.
type Pos struct {
	Offset int // offset, starting at 0
	Line   int // line number, starting at 1
	Column int // column number, starting at 1 (character count)
}

// IsValid returns true if the position is valid.
func (p *Pos) IsValid() bool { return p.Line > 0 }

// String returns a string in one of two forms:
//
//	line:column         valid position without file name
//	-                   invalid position without file name
func (p Pos) String() string {
	s := ""
	if p.IsValid() {
		s = fmt.Sprintf("%d:%d", p.Line, p.Column)
	}
	if s == "" {
		s = "-"
	}
	return s
}
