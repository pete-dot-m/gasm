package parser

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

/*

	^\w+:$					Labels
	^\t*[a-zA-Z]+ or ^\t*[a-zA-Z]{3}	Opcode (MOV, lda, etc.)
	;[a-zA-Z ]*$			Comment
	#\$[0-9a-fA-F]{1,4}		literal hex value #$ (2-4)
	\$[0-9a-fA-F]{4}		word (address)
	\$[0-9a-fA-F]{2}		byte
	\"[a-zA-Z ]*\"$		string literal
	[AXYZ]			register (A X Y Z)
	^\.[a-z]*		directive
*/

type Token int

const (
	EOF = iota
	ILLEGAL
	COMMA
	OP
	SEMI
	LABEL
	COMMENT
	BYTE
	WORD
	LITERAL
	STRING
	DIRECTIVE
)

var tokens = []string{
	EOF:       "EOF",
	ILLEGAL:   "ILLEGAL",
	OP:        "OP",
	SEMI:      ";",
	LABEL:     "LABEL",
	COMMENT:   "COMMENT",
	BYTE:      "BYTE",
	WORD:      "WORD",
	LITERAL:   "LITERAL",
	STRING:    "STRING",
	DIRECTIVE: "DIRECTIVE",
}

const (
	labelRegex     = `^\w+:$`
	commentRegex   = `;[a-zA-Z ]*`
	directiveRegex = `^\s*\.[a-z]*`
	stringRegex    = `\"[a-zA-Z ]*\"`
)

func (t Token) String() string {
	return tokens[t]
}

type Line struct {
	Index int
	Type  Token
	Value string
}

type Lexer struct {
	Line    int
	Scanner *bufio.Scanner
	Lines   []Line
}

func NewLexer(reader io.Reader) *Lexer {
	return &Lexer{
		Line:    1,
		Scanner: bufio.NewScanner(reader),
		Lines:   []Line{},
	}
}

func (l *Lexer) String() string {
	s := fmt.Sprintf("lines lexed: %d\n", l.Line)

	for _, d := range l.Lines {
		s = s + fmt.Sprintf("%d %v %s\n", d.Index, d.Type, d.Value)
	}

	return s
}

func (l *Lexer) Lex() error {
	// initial pass..
	// remove comments
	// collect lines and identify LABELs, DIRECTIVEs, and OPs
	return l.linePass()
}

func (l *Lexer) linePass() error {
	labelRe, err := regexp.Compile(labelRegex)
	if err != nil {
		return err
	}
	commentRe, err := regexp.Compile(commentRegex)
	if err != nil {
		return err
	}
	directiveRe, err := regexp.Compile(directiveRegex)
	if err != nil {
		return err
	}

	for l.Scanner.Scan() {
		line := l.Scanner.Text()

		// labels and directives live on their own
		// line, so pull those out first
		if label := labelRe.FindString(line); label != "" {
			trimmed := strings.TrimSpace(label)
			if trimmed != "" {
				l.Lines = append(l.Lines, Line{Index: l.Line, Type: LABEL, Value: trimmed})
			}
		} else if directiveRe.MatchString(line) {
			trimmed := strings.TrimSpace(line)
			if trimmed != "" {
				l.Lines = append(l.Lines, Line{Index: l.Line, Type: DIRECTIVE, Value: trimmed})
			}
		} else {
			if loc := commentRe.FindStringIndex(line); loc != nil {
				// clip off the comment and save the line
				t := strings.TrimSpace(line[:loc[0]])
				if t != "" {
					l.Lines = append(l.Lines, Line{Index: l.Line, Type: OP, Value: t})
				}
			} else if line != "" {
				t := strings.TrimSpace(line)
				l.Lines = append(l.Lines, Line{Index: l.Line, Type: OP, Value: t})
			}
		}
		l.Line++
	}
	return nil
}

/*
func (l *Lexer) Lex() (Position, Token, string) {
	// loop until we return a token
	for {
		r, _, err := l.Reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return l.Pos, EOF, ""
			}

			panic(err)
		}

		l.Pos.Column++
		switch r {
		case '\n':
			l.resetPosition()
		case ';':
			startPos := l.Pos
			startPos.Column++
			lit := l.lexIdent()
			return startPos, COMMENT, ";" + lit
		default:
			if unicode.IsSpace(r) {
				continue
			} else if unicode.IsDigit(r) {
				// backup and let lexInt rescan the beginning of the int
				startPos := l.Pos
				l.backup()
				lit := l.lexInt()
				return startPos, INT, lit
			} else if unicode.IsLetter(r) {
				// backup and let lexIdent rescan the beginning of the ident
				startPos := l.Pos
				l.backup()
				lit := l.lexIdent()
				return startPos, OP, lit
			} else {
				return l.Pos, ILLEGAL, string(r)
			}
		}
	}
}*/

/*
func (l *Lexer) resetPosition() {
	l.Pos.Line++
	l.Pos.Column = 0
}

func (l *Lexer) backup() {
	if err := l.Reader.UnreadRune(); err != nil {
		panic(err)
	}
	l.Pos.Column--
}

func (l *Lexer) lexInt() string {
	var lit string
	for {
		r, _, err := l.Reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				// at the end of the int
				return lit
			}
		}

		l.Pos.Column++
		if unicode.IsDigit(r) {
			lit = lit + string(r)
		} else {
			// scanned something not in the integer
			l.backup()
			return lit
		}
	}
}

func (l *Lexer) lexIdent() string {
	var lit string
	for {
		r, _, err := l.Reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return lit
			}
		}

		l.Pos.Column++
		if unicode.IsLetter(r) {
			lit = lit + string(r)
		} else {
			// something not in the identifier
			l.backup()
			return lit
		}
	}
}
*/
