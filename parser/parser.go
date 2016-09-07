package sanparser

import (
	"errors"
	"fmt"

	ast "github.com/fgrehm/go-san/ast"
	scanner "github.com/fgrehm/go-san/scanner"
	token "github.com/fgrehm/go-san/token"
)

type parser struct {
	sc *scanner.Scanner

	// Last read token
	tok token.Token

	comments    []*ast.CommentGroup
	leadComment *ast.CommentGroup // last lead comment
	lineComment *ast.CommentGroup // last line comment

	enableTrace bool
	indent      int
	n           int // buffer size (max = 1)
}

type blockParserFunc func(f *ast.File, p *parser, firstToken token.Token) error

var parserMap = map[token.Type]blockParserFunc{
	token.IDENTIFIERS:  parseIdentifiers,
	token.EVENTS:       parseEvents,
	token.PARTIAL:      parseReachability,
	token.REACHABILITY: parseReachability,
	token.NETWORK:      parseNetwork,
	token.RESULTS:      parseResults,
}

// Parser defines a syntatic parser for SAN models
type Parser interface {
	// Parse parses a SAN model into an abstract syntax tree
	Parse() (*ast.File, error)
}

// New returns a new parser for the provided source
func New(src []byte) Parser {
	return &parser{
		sc: scanner.New(src),
		// enableTrace: true,
	}
}

// Parse returns the fully parsed source and returns the abstract syntax tree.
func Parse(src []byte) (*ast.File, error) {
	p := New(src)
	return p.Parse()
}

// Parse returns the fully parsed source and returns the abstract syntax tree.
func (p *parser) Parse() (*ast.File, error) {
	var err, scerr error
	p.sc.Error = func(pos token.Pos, msg string) {
		scerr = p.err(pos, errors.New(msg))
	}

	file, err := p.file()
	if scerr != nil {
		return nil, scerr
	}
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (p *parser) file() (*ast.File, error) {
	file := &ast.File{}

	defer un(trace(p, "ParseFile"))

	for {
		tok := p.scan()
		if tok.Type == token.EOF {
			return file, nil
		}
		if blockParser, ok := parserMap[tok.Type]; ok {
			if err := blockParser(file, p, tok); err != nil {
				return nil, err
			}
		} else {
			return nil, p.err(tok.Pos, fmt.Errorf("Unexpected token %q found", tok.Text))
		}
	}
}

func (p *parser) scanExpression() (*ast.Expression, error) {
	exp := &ast.Expression{[]token.Token{}}
	for {
		tok := p.scan()
		if tok.Type == token.EOF {
			return nil, p.err(tok.Pos, fmt.Errorf("Unexpected token found: %q. Expected a ;", tok.Text))
		}
		if tok.Type == token.SEMICOLON {
			if len(exp.Tokens) == 0 {
				return nil, p.err(tok.Pos, fmt.Errorf("Invalid expression"))
			}
			break
		}
		if tok.Type.IsKeyword() && tok.Type != token.ST {
			p.unscan()
			break
		}
		exp.Tokens = append(exp.Tokens, tok)
	}
	return exp, nil
}

// scan returns the next token from the underlying scanner. If a token has
// been unscanned then read that instead. In the process, it collects any
// comment groups encountered, and remembers the last lead and line comments.
func (p *parser) scan() token.Token {
	// If we have a token on the buffer, then return it.
	if p.n != 0 {
		p.n = 0
		return p.tok
	}

	// Otherwise read the next token from the scanner and Save it to the buffer
	// in case we unscan later.
	prev := p.tok
	p.tok = p.sc.Scan()

	if p.tok.Type == token.COMMENT {
		var comment *ast.CommentGroup
		var endline int

		// fmt.Printf("p.tok.Pos.Line = %+v prev: %d endline %d \n",
		// p.tok.Pos.Line, prev.Pos.Line, endline)
		if p.tok.Pos.Line == prev.Pos.Line {
			// The comment is on same line as the previous token; it
			// cannot be a lead comment but may be a line comment.
			comment, endline = p.consumeCommentGroup(0)
			if p.tok.Pos.Line != endline {
				// The next token is on a different line, thus
				// the last comment group is a line comment.
				p.lineComment = comment
			}
		}

		// consume successor comments, if any
		endline = -1
		for p.tok.Type == token.COMMENT {
			comment, endline = p.consumeCommentGroup(1)
		}

		if endline+1 == p.tok.Pos.Line {
			// The next token is following on the line immediately after the
			// comment group, thus the last comment group is a lead comment.
			p.leadComment = comment
		}
	}

	return p.tok
}

func (p *parser) consumeComment() (comment *ast.Comment, endline int) {
	endline = p.tok.Pos.Line

	// count the endline if it's multiline comment, ie starting with /*
	if len(p.tok.Text) > 1 && p.tok.Text[1] == '*' {
		// don't use range here - no need to decode Unicode code points
		for i := 0; i < len(p.tok.Text); i++ {
			if p.tok.Text[i] == '\n' {
				endline++
			}
		}
	}

	comment = &ast.Comment{Start: p.tok.Pos, Text: p.tok.Text}
	p.tok = p.sc.Scan()
	return
}

func (p *parser) consumeCommentGroup(n int) (comments *ast.CommentGroup, endline int) {
	var list []*ast.Comment
	endline = p.tok.Pos.Line

	for p.tok.Type == token.COMMENT && p.tok.Pos.Line <= endline+n {
		var comment *ast.Comment
		comment, endline = p.consumeComment()
		list = append(list, comment)
	}

	// add comment group to the comments list
	comments = &ast.CommentGroup{List: list}
	p.comments = append(p.comments, comments)

	return
}

// unscan pushes the previously read token back onto the buffer.
func (p *parser) unscan() {
	p.n = 1
}

// ----------------------------------------------------------------------------
// Parsing support

func (p *parser) err(pos token.Pos, err error) *PosError {
	return &PosError{Pos: pos, Err: err}
}

func (p *parser) printTrace(a ...interface{}) {
	if !p.enableTrace {
		return
	}

	const dots = ". . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . "
	const n = len(dots)
	fmt.Printf("%5d:%3d: ", p.tok.Pos.Line, p.tok.Pos.Column)

	i := 2 * p.indent
	for i > n {
		fmt.Print(dots)
		i -= n
	}
	// i <= n
	fmt.Print(dots[0:i])
	fmt.Println(a...)
}

func trace(p *parser, msg string) *parser {
	p.printTrace(msg, "(")
	p.indent++
	return p
}

// Usage pattern: defer un(trace(p, "..."))
func un(p *parser) {
	p.indent--
	p.printTrace(")")
}
