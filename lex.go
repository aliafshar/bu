package bu

import (
	"bufio"
	"fmt"
	"io"
	"unicode/utf8"
)

type tokenType string

type pos struct {
	start int
	end   int
}

const (
	tokenName       tokenType = "NAME"
	tokenColon      tokenType = "COLON"
	tokenEquals     tokenType = "EQUALS"
	tokenPling      tokenType = "PLING"
	tokenNewline    tokenType = "NEWLINE"
	tokenWhitespace tokenType = "WHITESPACE"
  tokenIndent     tokenType = "INDENT"
	tokenEof        tokenType = "EOF"
	tokenSof        tokenType = "SOF"
	tokenRaw        tokenType = "RAW"
	tokenComment    tokenType = "COMMENT"
  tokenLessthan tokenType = "LESSTHAN"
	eof                       = 4
  lf = 10
)

var keywords = map[string]tokenType{
	"!":  tokenPling,
	"\n": tokenNewline,
	":":  tokenColon,
	"=":  tokenEquals,
	" ":  tokenWhitespace,
	"\t": tokenWhitespace,
	"#":  tokenComment,
  "<": tokenLessthan,
}

type token struct {
	typ tokenType
	pos pos
  line int
	val []rune
}

func (i *token) String() string {
	return fmt.Sprintf("Token: %v, %q", i.typ, string(i.val))
}

func (t *token) value() string {
	return string(t.val)
}

type lexer struct {
	tokens           []*token
	lastToken        *token
	last2Token       *token
	currentTokenType tokenType
	currentRune      rune
	r                io.ReadCloser
	out              chan *token
}

func NewLexer(r io.ReadCloser) *lexer {
	s := &token{typ: tokenSof, pos: pos{start: -1}}
	return &lexer{tokens: []*token{s}, lastToken: s, last2Token: s, out: make(chan *token), r: r}
}

func (l *lexer) isRaw() bool {
	return l.isLast(tokenRaw)
}

func (l *lexer) isLast(typ tokenType) bool {
	return l.lastToken.typ == typ
}

func (l *lexer) isSecondLast(typ tokenType) bool {
	return l.last2Token.typ == typ
}

func (l *lexer) isCurrentEqualsLast() bool {
	return l.currentTokenType == l.lastToken.typ
}

func (l *lexer) isCurrent(typ tokenType) bool {
	return l.currentTokenType == typ
}

func (l *lexer) shouldRaw() bool {
	return l.isLast(tokenEquals) || l.isLast(tokenComment) || l.isLast(tokenIndent)
}

func (l *lexer) doRaw(r rune, char int) {
  if l.isLast(tokenEquals) {
    l.newToken(tokenRaw, r, char)
  } else {
		l.newToken(tokenRaw, l.lastToken.val[0], char)
		l.extendToken(r)
  }
}

func (l *lexer) shouldIndent() bool {
  return l.isCurrent(tokenWhitespace) && (l.isLast(tokenNewline) || l.isLast(tokenSof))
}

func (l *lexer) doIndent(r rune, char int) {
  l.newToken(tokenIndent, r, char)
}


func (l *lexer) isExtendable() bool {
	return l.isCurrentEqualsLast() &&
		(l.isCurrent(tokenName) || l.isCurrent(tokenWhitespace) || l.isCurrent(tokenNewline))
}

func (l *lexer) newToken(typ tokenType, r rune, char int) {
  t := &token{typ: typ, val: []rune{r}, pos: pos{start: char, end: char + 1}}
	l.out <- l.lastToken
	l.last2Token = l.lastToken
	l.lastToken = t
	l.tokens = append(l.tokens, t)
}

func (l *lexer) extendToken(r rune) {
	l.lastToken.val = append(l.lastToken.val, r)
  l.lastToken.pos.end++
}

func (l *lexer) eof() {
	l.newToken(tokenEof, 4, -1)
	l.out <- l.lastToken
	close(l.out)
}

func (l *lexer) feed(r rune, char int) {
	typ, isKw := keywords[string(r)]
	l.currentTokenType = typ
	if l.isRaw() && !l.isCurrent(tokenNewline) {
		l.extendToken(r)
		return // otherwise drop out of raw mode
	}
	if l.shouldRaw() {
    l.doRaw(r, char)
		return
	}
  if l.shouldIndent() {
    l.doIndent(r, char)
    return
  }
	if !isKw {
		l.currentTokenType = tokenName
	}
	if l.isExtendable() {
		l.extendToken(r)
		return
	}
	l.newToken(l.currentTokenType, r, char)
}

func (l *lexer) lex() {
	scanner := bufio.NewScanner(l.r)
	scanner.Split(bufio.ScanRunes)
	i := 0
	for scanner.Scan() {
		r, _ := utf8.DecodeLastRuneInString(scanner.Text())
		l.feed(r, i)
		i++
	}
	l.eof()
}
