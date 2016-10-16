package sanscanner

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	token "github.com/fgrehm/go-san/token"
)

var f100 = strings.Repeat("f", 100)

type tokenPair struct {
	tok  token.Type
	text string
}

var tokenLists = map[string][]tokenPair{
	"semicolon": []tokenPair{
		{token.SEMICOLON, ";"},
	},
	"comment": []tokenPair{
		{token.COMMENT, "//"},
		{token.COMMENT, "////"},
		{token.COMMENT, "// comment"},
		{token.COMMENT, "// /* comment */"},
		{token.COMMENT, "// // comment //"},
		{token.COMMENT, "//" + f100},
		{token.COMMENT, "/**/"},
		{token.COMMENT, "/***/"},
		{token.COMMENT, "/* comment */"},
		{token.COMMENT, "/* // comment */"},
		{token.COMMENT, "/* /* comment */"},
		{token.COMMENT, "/*\n comment\n*/"},
		{token.COMMENT, "/*" + f100 + "*/"},
	},
	"operator": []tokenPair{
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.ASSIGN, "="},
		{token.SUM, "+"},
		{token.MULT, "*"},
		{token.DIV, "/"},
		{token.AND, "&&"},
		{token.EQUAL, "=="},
		{token.NEQUAL, "!="},
	},
	"ident": []tokenPair{
		{token.IDENTIFIER, "a"},
		{token.IDENTIFIER, "a0"},
		{token.IDENTIFIER, "foobar"},
		{token.IDENTIFIER, "foo-bar"},
		{token.IDENTIFIER, "abc123"},
		{token.IDENTIFIER, "LGTM"},
		{token.IDENTIFIER, "_"},
		{token.IDENTIFIER, "_abc123"},
		{token.IDENTIFIER, "abc123_"},
		{token.IDENTIFIER, "_abc_123_"},
		{token.IDENTIFIER, "_äöü"},
		{token.IDENTIFIER, "_本"},
		{token.IDENTIFIER, "äöü"},
		{token.IDENTIFIER, "本"},
		{token.IDENTIFIER, "a۰۱۸"},
		{token.IDENTIFIER, "foo६४"},
		{token.IDENTIFIER, "bar９８７６"},
	},
	"keyword": []tokenPair{
		{token.IDENTIFIERS, "identifiers"},
		{token.EVENTS, "events"},
		{token.PARTIAL, "partial"},
		{token.REACHABILITY, "reachability"},
		{token.NETWORK, "network"},
		{token.CONTINUOUS, "continuous"},
		{token.AUT, "aut"},
		{token.ST, "st"},
		{token.STT, "stt"},
		{token.TO, "to"},
		{token.RESULTS, "results"},
	},
	"number": []tokenPair{
		{token.NUMBER, "0"},
		{token.NUMBER, "1"},
		{token.NUMBER, "9"},
		{token.NUMBER, "42"},
		{token.NUMBER, "1234567890"},
		{token.NUMBER, "00"},
		{token.NUMBER, "01"},
		{token.NUMBER, "07"},
		{token.NUMBER, "042"},
		{token.NUMBER, "01234567"},
		{token.NUMBER, "0x0"},
		{token.NUMBER, "0x1"},
		{token.NUMBER, "0xf"},
		{token.NUMBER, "0x42"},
		{token.NUMBER, "0x123456789abcDEF"},
		{token.NUMBER, "0x" + f100},
		{token.NUMBER, "0X0"},
		{token.NUMBER, "0X1"},
		{token.NUMBER, "0XF"},
		{token.NUMBER, "0X42"},
		{token.NUMBER, "0X123456789abcDEF"},
		{token.NUMBER, "0X" + f100},
		{token.NUMBER, "-0"},
		{token.NUMBER, "-1"},
		{token.NUMBER, "-9"},
		{token.NUMBER, "-42"},
		{token.NUMBER, "-1234567890"},
		{token.NUMBER, "-00"},
		{token.NUMBER, "-01"},
		{token.NUMBER, "-07"},
		{token.NUMBER, "-29"},
		{token.NUMBER, "-042"},
		{token.NUMBER, "-01234567"},
		{token.NUMBER, "-0x0"},
		{token.NUMBER, "-0x1"},
		{token.NUMBER, "-0xf"},
		{token.NUMBER, "-0x42"},
		{token.NUMBER, "-0x123456789abcDEF"},
		{token.NUMBER, "-0x" + f100},
		{token.NUMBER, "-0X0"},
		{token.NUMBER, "-0X1"},
		{token.NUMBER, "-0XF"},
		{token.NUMBER, "-0X42"},
		{token.NUMBER, "-0X123456789abcDEF"},
		{token.NUMBER, "-0X" + f100},
	},
	"float": []tokenPair{
		{token.FLOAT, "0."},
		{token.FLOAT, "1."},
		{token.FLOAT, "42."},
		{token.FLOAT, "01234567890."},
		{token.FLOAT, "0.0"},
		{token.FLOAT, "1.0"},
		{token.FLOAT, "42.0"},
		{token.FLOAT, "01234567890.0"},
		{token.FLOAT, "0e0"},
		{token.FLOAT, "1e0"},
		{token.FLOAT, "42e0"},
		{token.FLOAT, "01234567890e0"},
		{token.FLOAT, "0E0"},
		{token.FLOAT, "1E0"},
		{token.FLOAT, "42E0"},
		{token.FLOAT, "01234567890E0"},
		{token.FLOAT, "0e+10"},
		{token.FLOAT, "1e-10"},
		{token.FLOAT, "42e+10"},
		{token.FLOAT, "01234567890e-10"},
		{token.FLOAT, "0E+10"},
		{token.FLOAT, "1E-10"},
		{token.FLOAT, "42E+10"},
		{token.FLOAT, "01234567890E-10"},
		{token.FLOAT, "01.8e0"},
		{token.FLOAT, "1.4e0"},
		{token.FLOAT, "42.2e0"},
		{token.FLOAT, "01234567890.12e0"},
		{token.FLOAT, "0.E0"},
		{token.FLOAT, "1.12E0"},
		{token.FLOAT, "42.123E0"},
		{token.FLOAT, "01234567890.213E0"},
		{token.FLOAT, "0.2e+10"},
		{token.FLOAT, "1.2e-10"},
		{token.FLOAT, "42.54e+10"},
		{token.FLOAT, "01234567890.98e-10"},
		{token.FLOAT, "0.1E+10"},
		{token.FLOAT, "1.1E-10"},
		{token.FLOAT, "42.1E+10"},
		{token.FLOAT, "01234567890.1E-10"},
		{token.FLOAT, "-0.0"},
		{token.FLOAT, "-1.0"},
		{token.FLOAT, "-42.0"},
		{token.FLOAT, "-01234567890.0"},
		{token.FLOAT, "-0e0"},
		{token.FLOAT, "-1e0"},
		{token.FLOAT, "-42e0"},
		{token.FLOAT, "-01234567890e0"},
		{token.FLOAT, "-0E0"},
		{token.FLOAT, "-1E0"},
		{token.FLOAT, "-42E0"},
		{token.FLOAT, "-01234567890E0"},
		{token.FLOAT, "-0e+10"},
		{token.FLOAT, "-1e-10"},
		{token.FLOAT, "-42e+10"},
		{token.FLOAT, "-01234567890e-10"},
		{token.FLOAT, "-0E+10"},
		{token.FLOAT, "-1E-10"},
		{token.FLOAT, "-42E+10"},
		{token.FLOAT, "-01234567890E-10"},
		{token.FLOAT, "-01.8e0"},
		{token.FLOAT, "-1.4e0"},
		{token.FLOAT, "-42.2e0"},
		{token.FLOAT, "-01234567890.12e0"},
		{token.FLOAT, "-0.E0"},
		{token.FLOAT, "-1.12E0"},
		{token.FLOAT, "-42.123E0"},
		{token.FLOAT, "-01234567890.213E0"},
		{token.FLOAT, "-0.2e+10"},
		{token.FLOAT, "-1.2e-10"},
		{token.FLOAT, "-42.54e+10"},
		{token.FLOAT, "-01234567890.98e-10"},
		{token.FLOAT, "-0.1E+10"},
		{token.FLOAT, "-1.1E-10"},
		{token.FLOAT, "-42.1E+10"},
		{token.FLOAT, "-01234567890.1E-10"},
	},
}

var orderedTokenLists = []string{
	"semicolon",
	"comment",
	"operator",
	"ident",
	"keyword",
	"number",
	"float",
}

func TestPosition(t *testing.T) {
	// create artifical source code
	buf := new(bytes.Buffer)

	for _, listName := range orderedTokenLists {
		for _, ident := range tokenLists[listName] {
			fmt.Fprintf(buf, "\t\t\t\t%s\n", ident.text)
		}
	}

	s := New(buf.Bytes())

	pos := token.Pos{4, 1, 5}
	s.Scan()
	for _, listName := range orderedTokenLists {

		for _, k := range tokenLists[listName] {
			curPos := s.tokPos
			// fmt.Printf("[%q] s = %+v:%+v\n", k.text, curPos.Offset, curPos.Column)

			if curPos.Offset != pos.Offset {
				t.Fatalf("offset = %d, want %d for %q", curPos.Offset, pos.Offset, k.text)
			}
			if curPos.Line != pos.Line {
				t.Fatalf("line = %d, want %d for %q", curPos.Line, pos.Line, k.text)
			}
			if curPos.Column != pos.Column {
				t.Fatalf("column = %d, want %d for %q", curPos.Column, pos.Column, k.text)
			}
			pos.Offset += 4 + len(k.text) + 1     // 4 tabs + token bytes + newline
			pos.Line += countNewlines(k.text) + 1 // each token is on a new line
			s.Scan()
		}
	}
	// make sure there were no token-internal errors reported by scanner
	if s.ErrorCount != 0 {
		t.Errorf("%d errors", s.ErrorCount)
	}
}

func TestNullChar(t *testing.T) {
	s := New([]byte("\"\\0"))
	s.Scan() // Used to panic
}

func TestSemicolon(t *testing.T) {
	testTokenList(t, tokenLists["semicolon"])
}

func TestComment(t *testing.T) {
	testTokenList(t, tokenLists["comment"])
}

func TestOperator(t *testing.T) {
	testTokenList(t, tokenLists["operator"])
}

func TestIdent(t *testing.T) {
	testTokenList(t, tokenLists["ident"])
}

func TestNumber(t *testing.T) {
	testTokenList(t, tokenLists["number"])
}

func TestFloat(t *testing.T) {
	testTokenList(t, tokenLists["float"])
}

func TestKeywords(t *testing.T) {
	testTokenList(t, tokenLists["keyword"])
}

func TestWindowsLineEndings(t *testing.T) {
	san := `// This should have Windows line endings
identifiers
	r_proc = 6;`

	sanWindowsEndings := strings.Replace(san, "\n", "\r\n", -1)

	literals := []struct {
		tokenType token.Type
		literal   string
	}{
		{token.COMMENT, "// This should have Windows line endings\r"},
		{token.IDENTIFIERS, `identifiers`},
		{token.IDENTIFIER, `r_proc`},
		{token.ASSIGN, `=`},
		{token.NUMBER, `6`},
		{token.SEMICOLON, `;`},
	}

	s := New([]byte(sanWindowsEndings))
	for _, l := range literals {
		tok := s.Scan()

		if l.tokenType != tok.Type {
			t.Errorf("got: %s want %s for %s\n", tok, l.tokenType, tok.String())
		}

		if l.literal != tok.Text {
			t.Errorf("got:\n%v\nwant:\n%v\n", []byte(tok.Text), []byte(l.literal))
		}
	}
}

func TestRealExample(t *testing.T) {
	complexSAN := `// This is based on the basic client server example
identifiers
  r_proc    = 6;
  F1 = (st Client == Working) * 1;
	F2 = r_proc / 2;

events
  loc l_proc    (r_proc);
	syn s_resp    (r_resp);

partial reachability = ((st Client == Idle) && (st Server == Idle));

network ClientServer (continuous)
  aut Client
    stt Idle         to (Transmitting) s_req
                     to (Idle)         l_no_more

results
  Client_processing      = (st Client == Working);
  Client_trans_Serv_rcv  = ((st Client == Transmitting)
                            && (st Server == Receiving));
`

	literals := []struct {
		tokenType token.Type
		literal   string
	}{
		{token.COMMENT, `// This is based on the basic client server example`},
		{token.IDENTIFIERS, `identifiers`},
		{token.IDENTIFIER, `r_proc`},
		{token.ASSIGN, `=`},
		{token.NUMBER, `6`},
		{token.SEMICOLON, `;`},
		{token.IDENTIFIER, `F1`},
		{token.ASSIGN, `=`},
		{token.LPAREN, `(`},
		{token.ST, `st`},
		{token.IDENTIFIER, `Client`},
		{token.EQUAL, `==`},
		{token.IDENTIFIER, `Working`},
		{token.RPAREN, `)`},
		{token.MULT, `*`},
		{token.NUMBER, `1`},
		{token.SEMICOLON, `;`},
		{token.IDENTIFIER, `F2`},
		{token.ASSIGN, `=`},
		{token.IDENTIFIER, `r_proc`},
		{token.DIV, `/`},
		{token.NUMBER, `2`},
		{token.SEMICOLON, `;`},

		{token.EVENTS, `events`},
		{token.LOC, `loc`},
		{token.IDENTIFIER, `l_proc`},
		{token.LPAREN, `(`},
		{token.IDENTIFIER, `r_proc`},
		{token.RPAREN, `)`},
		{token.SEMICOLON, `;`},
		{token.SYN, `syn`},
		{token.IDENTIFIER, `s_resp`},
		{token.LPAREN, `(`},
		{token.IDENTIFIER, `r_resp`},
		{token.RPAREN, `)`},
		{token.SEMICOLON, `;`},

		{token.PARTIAL, `partial`},
		{token.REACHABILITY, `reachability`},
		{token.ASSIGN, `=`},
		{token.LPAREN, `(`},
		{token.LPAREN, `(`},
		{token.ST, `st`},
		{token.IDENTIFIER, `Client`},
		{token.EQUAL, `==`},
		{token.IDENTIFIER, `Idle`},
		{token.RPAREN, `)`},
		{token.AND, `&&`},
		{token.LPAREN, `(`},
		{token.ST, `st`},
		{token.IDENTIFIER, `Server`},
		{token.EQUAL, `==`},
		{token.IDENTIFIER, `Idle`},
		{token.RPAREN, `)`},
		{token.RPAREN, `)`},
		{token.SEMICOLON, `;`},

		{token.NETWORK, `network`},
		{token.IDENTIFIER, `ClientServer`},
		{token.LPAREN, `(`},
		{token.CONTINUOUS, "continuous"},
		{token.RPAREN, `)`},
		{token.AUT, `aut`},
		{token.IDENTIFIER, `Client`},
		{token.STT, `stt`},
		{token.IDENTIFIER, `Idle`},
		{token.TO, `to`},
		{token.LPAREN, `(`},
		{token.IDENTIFIER, `Transmitting`},
		{token.RPAREN, `)`},
		{token.IDENTIFIER, `s_req`},
		{token.TO, `to`},
		{token.LPAREN, `(`},
		{token.IDENTIFIER, `Idle`},
		{token.RPAREN, `)`},
		{token.IDENTIFIER, `l_no_more`},

		{token.RESULTS, `results`},
		{token.IDENTIFIER, `Client_processing`},
		{token.ASSIGN, `=`},
		{token.LPAREN, `(`},
		{token.ST, `st`},
		{token.IDENTIFIER, `Client`},
		{token.EQUAL, `==`},
		{token.IDENTIFIER, `Working`},
		{token.RPAREN, `)`},
		{token.SEMICOLON, `;`},
		{token.IDENTIFIER, `Client_trans_Serv_rcv`},
		{token.ASSIGN, `=`},
		{token.LPAREN, `(`},
		{token.LPAREN, `(`},
		{token.ST, `st`},
		{token.IDENTIFIER, `Client`},
		{token.EQUAL, `==`},
		{token.IDENTIFIER, `Transmitting`},
		{token.RPAREN, `)`},
		{token.AND, `&&`},
		{token.LPAREN, `(`},
		{token.ST, `st`},
		{token.IDENTIFIER, `Server`},
		{token.EQUAL, `==`},
		{token.IDENTIFIER, `Receiving`},
		{token.RPAREN, `)`},
		{token.RPAREN, `)`},
		{token.SEMICOLON, `;`},
		{token.EOF, ``},
	}

	s := New([]byte(complexSAN))
	for _, l := range literals {
		tok := s.Scan()
		if l.tokenType != tok.Type {
			t.Errorf("got: %s want %s for %s\n", tok, l.tokenType, tok.String())
		}

		if l.literal != tok.Text {
			t.Errorf("got:\n%+v\n%s\n want:\n%+v\n%s\n", []byte(tok.String()), tok, []byte(l.literal), l.literal)
		}
	}

}

func TestError(t *testing.T) {
	testError(t, "\x80", "1:1", "illegal UTF-8 encoding", token.ILLEGAL)
	testError(t, "\xff", "1:1", "illegal UTF-8 encoding", token.ILLEGAL)

	testError(t, "ab\x80", "1:3", "illegal UTF-8 encoding", token.IDENTIFIER)
	testError(t, "abc\xff", "1:4", "illegal UTF-8 encoding", token.IDENTIFIER)

	testError(t, `&`, "1:1", "illegal char &", token.ILLEGAL)
	testError(t, `!`, "1:1", "illegal char !", token.ILLEGAL)

	testError(t, `01238`, "1:6", "illegal octal number", token.NUMBER)
	testError(t, `01238123`, "1:9", "illegal octal number", token.NUMBER)
	testError(t, `0x`, "1:3", "illegal hexadecimal number", token.NUMBER)
	testError(t, `0xg`, "1:3", "illegal hexadecimal number", token.NUMBER)

	testError(t, `/*/`, "1:4", "comment not terminated", token.COMMENT)
}

func testError(t *testing.T, src, pos, msg string, tok token.Type) {
	s := New([]byte(src))

	errorCalled := false
	s.Error = func(p token.Pos, m string) {
		if !errorCalled {
			if pos != p.String() {
				t.Errorf("pos = %q, want %q for %q", p, pos, src)
			}

			if m != msg {
				t.Errorf("msg = %q, want %q for %q", m, msg, src)
			}
			errorCalled = true
		}
	}

	tk := s.Scan()
	if tk.Type != tok {
		t.Errorf("tok = %s, want %s for %q", tk, tok, src)
	}
	if !errorCalled {
		t.Errorf("error handler not called for %q", src)
	}
	if s.ErrorCount == 0 {
		t.Errorf("count = %d, want > 0 for %q", s.ErrorCount, src)
	}
}

func testTokenList(t *testing.T, tokenList []tokenPair) {
	// create artifical source code
	buf := new(bytes.Buffer)
	for _, ident := range tokenList {
		fmt.Fprintf(buf, "%s\n", ident.text)
	}

	s := New(buf.Bytes())
	for _, ident := range tokenList {
		tok := s.Scan()
		if tok.Type != ident.tok {
			t.Errorf("tok = %q want %q for %q\n", tok, ident.tok, ident.text)
		}

		if tok.Text != ident.text {
			t.Errorf("text = %q want %q", tok.String(), ident.text)
		}

	}
}

func countNewlines(s string) int {
	n := 0
	for _, ch := range s {
		if ch == '\n' {
			n++
		}
	}
	return n
}
