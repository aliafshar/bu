package bu

import (
	"io"
)

type parser struct {
	tokenStream chan *token
	root        *node
}

func newParser() *parser {
	return &parser{tokenStream: make(chan *token)}
}

func (p *parser) buildAst(r io.Reader) (*node, error) {
	b := newBuilder()
	go lex(r, p.tokenStream)
	for {
		t := <-p.tokenStream
		if t.is(tokenWhitespace) {
			// We just ignore whitespace from now on
			continue
		}
		if t.is(tokenEof) {
			break
		}
		b.feed(t)
	}
	return b.root, nil
}

var targetBuilders = map[opType]func(*node) target{
	opShell: func(n *node) target {
		return &shellTarget{name: n.key, body: n.body()}
	},
	opQuestion: func(n *node) target {
		return &questionTarget{name: n.key}
	},
}

func (p *parser) createShellTarget(n *node) target {
	t := &shellTarget{name: n.key, body: n.body(), shell: shells["sh"]}
	for _, o := range n.nodes {
		switch o.op {
		case opUnnamed:
			t.deps = append(t.deps, &targetDependency{name: o.key})
			t.depsNames = append(t.depsNames, o.key)
		case opShell:
			t.shell = shells[o.key]
		case opQuestion:
			t.deps = append(t.deps, &fileDependency{filename: o.key})
		case opRedirect:
			t.outfile = o.key
		case opImport:
			t.infile = o.key
		}
	}
	return t
}

func (p *parser) createQuestionTarget(n *node) target {
	t := &questionTarget{name: n.key, usage: n.body()}
	for _, o := range n.nodes {
		switch o.op {
		case opUnnamed:
			t.dflt = o.key
			break
		}
	}

	return t
}

func (p *parser) createTarget(m *module, n *node) error {
	var t target
	switch n.op {
	case opQuestion:
		t = p.createQuestionTarget(n)
	case opShell:
		t = p.createShellTarget(n)
	}
	m.targets = append(m.targets, t)
	return nil
}

func (p *parser) createImport(m *module, n *node) error {
	for _, o := range n.nodes {
		if o.op == opUnnamed {
			i := &imports{key: o.key}
			m.imports = append(m.imports, i)
		}
	}
	return nil
}

func (p *parser) createSetvar(m *module, n *node) error {
	s := &setvar{key: n.key, value: n.body()}
	m.setvars = append(m.setvars, s)
	return nil
}

func (p *parser) parse(s *script, r io.Reader) (*module, error) {
	root, err := p.buildAst(r)
	p.root = root
	if err != nil {
		return nil, err
	}
	m := &module{}
	s.modules = append(s.modules, m)
	if s.module == nil {
		s.module = m
	}
	for _, bNode := range p.root.nodes {
		switch bNode.op {
		case opComment:
			// Ignore comments here.
			continue
		case opImport:
			p.createImport(m, bNode)
		case opQuestion, opShell:
			p.createTarget(m, bNode)
		case opSetvar:
			p.createSetvar(m, bNode)
		default:
			logicalError("Unknown node type", bNode)
		}
	}
	return m, nil
}
