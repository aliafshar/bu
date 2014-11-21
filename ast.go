package bu

import (
	"fmt"
	"strings"
)

type nodeType string

const (
	nodeRoot     = "ROOT"
	nodeBlock    = "BLOCK"
	nodeModule   = "MODULE"
	nodeImport   = "IMPORT"
	nodeOperator = "OPERATOR"
)

type opType string

const (
	opUnnamed  opType = "UNNAMED"
	opShell    opType = "SHELL"
	opPipe     opType = "PIPE"
	opRedirect opType = "REDIRECT"
	opTarget   opType = "TARGET"
	opQuestion opType = "QUESTION"
	opSetvar   opType = "SETVAR"
	opComment  opType = "COMMENT"
	opImport   opType = "IMPORT"
)

type node struct {
	nodeType nodeType
	key      string
	nodes    []*node
	lines    []string
	op       opType
	tokens   []*token
}

func (n *node) extend(line string) {
	n.lines = append(n.lines, line)
}

func (n *node) String() string {
	return fmt.Sprintf("<%v, %q, %q, %+v, %q>", n.nodeType, n.op, n.key, n.nodes, n.body())
}

func (n *node) body() string {
	if len(n.lines) == 0 {
		return ""
	}
	i := 0
	found := false
	for !found {
		r := string(n.lines[0][i])
		switch r {
		case " ", "\t":
			i = i + 1
		default:
			found = true
		}
	}
	ls := []string{}
	for _, line := range n.lines {
		ls = append(ls, line[i:])
	}
	return strings.Join(ls, "\n")
}

func (n *node) newNode() *node {
	m := &node{}
	n.nodes = append(n.nodes, m)
	return m
}

type builder struct {
	root   *node
	module *node
	block  *node
	line   []*token
}

func newBuilder() *builder {
	return &builder{root: &node{nodeType: nodeModule}}
}

func (p *builder) parseLine() {
	if len(p.line) == 0 {
		return
	}
	if len(p.line) == 1 {
		return // This is an error
	}
	fst := p.line[0]
	snd := p.line[1]
	// First handle indented lines
	if fst.tokenType == tokenIndent {
		p.block.extend(snd.value())
		return
	}
	// Now blocks
	p.block = p.module.newNode()
	p.block.tokens = p.line
	p.block.nodeType = nodeBlock
	var opIndex int
	switch p.line[0].tokenType {
	case tokenName:
		opIndex = 1
		p.block.key = p.line[0].value()
	case tokenLeft, tokenComment:
		opIndex = 0
	default:
		panic("not here!")
	}
	p.block.op = opTypes[p.line[opIndex].tokenType]
	p.parseCommand(p.line[opIndex+1:])
}

func (p *builder) parseCommand(cmd []*token) {
	currentOp := opUnnamed
	var n *node
	for _, t := range cmd {
		op, isOp := opTypes[t.tokenType]
		if isOp {
			currentOp = op
		}
		switch t.tokenType {
		case tokenName:
			n = p.block.newNode()
			n.nodeType = nodeOperator
			n.key = t.value()
			n.op = currentOp
		case tokenRaw:
			p.block.extend(t.value())
		}
		if t.tokenType == tokenName {
		}
	}
}

var nodeTypes = map[tokenType]nodeType{
	tokenColon:    nodeBlock,
	tokenEquals:   nodeBlock,
	tokenQuestion: nodeBlock,
	tokenPipe:     nodeBlock,
	tokenRight:    nodeBlock,
	tokenLeft:     nodeImport,
}

var opTypes = map[tokenType]opType{
	tokenRight:    opRedirect,
	tokenPipe:     opPipe,
	tokenPling:    opShell,
	tokenColon:    opShell,
	tokenQuestion: opQuestion,
	tokenEquals:   opSetvar,
	tokenComment:  opComment,
	tokenLeft:     opImport,
}

func (p *builder) feed(t *token) {
	switch t.tokenType {
	case tokenSof:
		p.module = p.root // TODO wtf
	case tokenSol:
		p.line = []*token{}
	case tokenEol:
		p.parseLine()
	default:
		p.line = append(p.line, t)
	}
}