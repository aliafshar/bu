package bu

import (
	"io"
  "github.com/aliafshar/toylog"
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
		if t.is(tokenWhitespace) || t.is(tokenQuote) {
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

func (p *parser) createTarget(n *node) *target {
	t := &target{
		name:     n.key,
		body:     n.body(),
		shell:    "sh",
		redirect: &redirect{},
	}
	for _, o := range n.nodes {
		switch o.op {
		case opUnnamed:
			t.deps = append(t.deps, &targetDependency{name: o.key})
		case opShell:
			t.shell = o.key
		case opQuestion:
			t.deps = append(t.deps, &fileDependency{filename: o.key})
		case opPipe:
			t.pipe = append(t.pipe, &targetDependency{name: o.key})
			t.deps = append(t.deps, &pipeDependency{name: o.key})
		case opAt:
			t.deps = append(t.deps, &webDependency{uri: o.key})
		case opRedirect:
			t.redirect.ofile = o.key
		case opImport:
			t.redirect.ifile = o.key
    case opCaret:
      t.watch = o.key
    default:
      toylog.Debugf("unknown operator %+v\n", o)
		}
	}
	return t
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
		case opQuestion, opShell, opPipe:
			t := p.createTarget(bNode)
			m.targets = append(m.targets, t)
		case opSetvar:
			p.createSetvar(m, bNode)
		default:
			logicalError("Unknown node type", bNode)
		}
	}
	return m, nil
}
