package bu

import (
	"fmt"
	"io"
)

type parser struct {
	tokenStream chan *token
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

var targetBuilders = map[opType]func(*node) target {
	opShell: func(n *node) target {
		return &shellTarget{name: n.key, body: n.body()}
	},
	opQuestion: func(n *node) target {
		return &questionTarget{name: n.key}
	},
}

func (p *parser) createShellTarget(n *node) target {
		t := &shellTarget{name: n.key, body: n.body(), shell: "sh"}
		for _, o := range n.nodes {
			switch o.op {
				case opUnnamed:
					t.deps = append(t.deps, &targetDependency{name: o.key})
					t.depsNames = append(t.depsNames, o.key)
				case opShell:
					t.shell = o.key
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
	fmt.Printf("%+v\n", t)
	return nil
}

func (p *parser) createImport(m *module, n *node) error {
	for _, o := range n.nodes {
		fmt.Println(o)
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
	root, _ := p.buildAst(r)
	fmt.Println(root)
	m := &module{}
	s.modules = append(s.modules, m)
	if s.module == nil {
		s.module = m
	}
	for _, bNode := range root.nodes {
		fmt.Println(bNode)
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
			panic(bNode)
		}
	}
	return m, nil
}
