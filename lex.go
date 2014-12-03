package bu

import (
	"bufio"
	"fmt"
	"io"
)

type tokenType int

const (
	tokenName       tokenType = iota
	tokenColon
	tokenEquals     tokenType = iota
	tokenPling      tokenType = iota
	tokenNewline    tokenType = iota
	tokenWhitespace tokenType = iota
	tokenIndent     tokenType = iota
	tokenEof        tokenType = iota
	tokenSof        tokenType = iota
	tokenRaw        tokenType = iota
	tokenComment    tokenType = iota
	tokenLeft       tokenType = iota
	tokenRight      tokenType = iota
	tokenQuestion   tokenType = iota
	tokenPipe       tokenType = iota
	tokenSol        tokenType = iota
	tokenEol        tokenType = iota
  tokenQuote      tokenType = iota
)

var operators = map[string]tokenType{
	"!":  tokenPling,
	"\n": tokenNewline,
	":":  tokenColon,
	"=":  tokenEquals,
	" ":  tokenWhitespace,
	"\t": tokenWhitespace,
	"#":  tokenComment,
	"<":  tokenLeft,
	"?":  tokenQuestion,
	"|":  tokenPipe,
	">":  tokenRight,
  "'":  tokenQuote,
}

var names = map[tokenType]string{
	tokenName:       "NAME",
	tokenColon:      "COLON",
	tokenEquals:     "EQUALS",
	tokenPling:      "PLING",
	tokenNewline:    "NEWLINE",
	tokenWhitespace: "WHITESPACE",
	tokenIndent:     "INDENT",
	tokenEof:        "EOF",
	tokenSof:        "SOF",
	tokenRaw:        "RAW",
	tokenComment:    "COMMENT",
	tokenLeft:       "LESSTHAN",
	tokenRight:      "GREATERTHAN",
	tokenQuestion:   "QUESTION",
	tokenPipe:       "PIPE",
	tokenSol:        "SOL",
	tokenEol:        "EOL",
  tokenQuote: "QUOTE",
}

type token struct {
	tokenType tokenType
	val       []rune
}

func (t *token) extend(c rune) {
	t.val = append(t.val, c)
}

func (t *token) typeName() string {
	return names[t.tokenType]
}

func (t *token) String() string {
	return fmt.Sprintf("<%v, %q>", t.typeName(), t.value())
}

func (t *token) value() string {
	return string(t.val)
}

func (t *token) is(typ tokenType) bool {
	return t.tokenType == typ
}

type lineLexer struct {
	tokens []*token
  isQuoted bool
}

func (l *lineLexer) String() string {
	return fmt.Sprintf("%+v", l.tokens)
}

func newSingleLine(typ tokenType) *lineLexer {
	return &lineLexer{tokens: []*token{&token{tokenType: typ}}}
}

func newLine() *lineLexer {
	return newSingleLine(tokenSol)
}

func (l *lineLexer) extend(t *token) {
	l.tokens = append(l.tokens, t)
}

func (l *lineLexer) newToken(typ tokenType, c rune) *token {
	t := &token{tokenType: typ, val: []rune{c}}
	l.extend(t)
	return t
}

func (l *lineLexer) newEmptyToken(typ tokenType) *token {
	t := &token{tokenType: typ}
	l.extend(t)
	return t
}

func (l *lineLexer) eol() {
	l.extend(&token{tokenType: tokenEol})
}

func (l *lineLexer) last() *token {
	return l.tokens[len(l.tokens)-1]
}

func (l *lineLexer) is(typ tokenType) bool {
	return l.last().tokenType == typ
}

func (l *lineLexer) lex(text string) {
	for _, c := range text {
		typ, isOp := operators[string(c)]
		if !isOp {
			typ = tokenName
		}
		switch l.last().tokenType {
		case tokenSol:
			if typ == tokenWhitespace {
				l.newEmptyToken(tokenIndent)
				l.newToken(tokenRaw, c)
				continue
			}
		case tokenComment, tokenEquals:
			l.newToken(tokenRaw, c)
			continue
    case tokenQuote:
      if !l.isQuoted {
        l.newToken(tokenName, c)
        l.isQuoted = true
        continue
      }
		case tokenRaw:
			l.last().extend(c)
			continue
		case tokenName:
			if typ == tokenName || (l.isQuoted && typ != tokenQuote) {
				l.last().extend(c)
				continue
			}
		}
		l.newToken(typ, c)
	}
	l.eol()
}

func (l *lineLexer) send(out chan *token) {
	for _, t := range l.tokens {
		out <- t
	}
}

func lex(r io.Reader, out chan *token) {
	out <- &token{tokenType: tokenSof}
	s := bufio.NewScanner(r)
	for s.Scan() {
		l := newLine()
		l.lex(s.Text())
		l.send(out)
	}
	out <- &token{tokenType: tokenEof}
}
