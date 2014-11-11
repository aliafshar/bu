package bu

import (
	"github.com/aliafshar/toylog"
	"strings"
)

type parser struct {
	module *module
	stack  []*token
	line   int
  node  lined
}

func parseError(msg string, t *token) {
	toylog.Errorf("Error parsing. "+msg+" [line %v %q]", t.line, t.value())
}

func newParser() *parser {
	return &parser{module: &module{targetIndex: map[string]target{}}, stack: []*token{}, line: 1}
}

func (p *parser) parse(l *lexer) {
	go l.lex()
	for l.lastToken.typ != tokenEof {
		t := <-l.out
		if t.typ != tokenWhitespace {
			p.feed(t)
		}
	}
	p.finalize()
}

func (p *parser) handleRaw(ts []*token) {
	if len(ts) < 2 {
		// Trailing whitespace.
		return
	}
	if p.node == nil {
		parseError("Target body outside target.", ts[1])
	} else {
		p.node.AppendBody(ts[1].value())
	}
}

func (p *parser) handleTarget(ts []*token) {
	tg := &shellTarget{
		name: string(ts[0].val),
		typ:  "sh",
	}
	p.module.targets = append(p.module.targets, tg)
	finishedDeps := false
	for _, t := range ts[2:] {
		switch t.typ {
		case tokenName:
			if finishedDeps {
				tg.typ = string(t.val)
			} else {
				tg.depsNames = append(tg.depsNames, string(t.val))
			}
		case tokenPling:
			finishedDeps = true
		}
	}
	p.node = tg
}

func (p *parser) handleQuestion(ts []*token) {
	tg := &questionTarget{name: string(ts[0].val), bodyLines: []string{}}
	p.module.targets = append(p.module.targets, tg)
	if len(ts) > 2 {
		tg.dflt = strings.Trim(ts[2].value(), " ")
	}
	p.node = tg
}

func (p *parser) handleSetvar(ts []*token) {
  tg := &setvar{key: ts[0].value()}
  p.module.setvars = append(p.module.setvars, tg)
  if len(ts) < 3 {
    return
  }
  if ts[2].typ == tokenRaw {
    tg.AppendBody(ts[2].value())
    return
  }
  p.node = tg
}

func (p *parser) handleNamed(ts []*token) {
	if len(ts) == 1 {
		parseError("Identifier doing nothing.", ts[0])
		return
	}
	switch ts[1].typ {
	case tokenEquals:
		p.handleSetvar(ts)
	case tokenQuestion:
		p.handleQuestion(ts)
	case tokenColon:
		p.handleTarget(p.stack)
	}
}

func (p *parser) handleImport(ts []*token) {
	if len(ts) < 2 {
		parseError("Import not provided.", p.stack[0])
		return
	}
	for _, t := range ts[1:] {
		p.module.imports = append(p.module.imports, t.value())
	}
}

func (p *parser) examine() {
	switch p.stack[0].typ {
	case tokenSof, tokenComment, tokenEof:
	case tokenIndent:
		p.handleRaw(p.stack)
	case tokenLessthan:
		p.handleImport(p.stack)
	case tokenName:
		p.handleNamed(p.stack)
	default:
		parseError("Unknown token.", p.stack[0])
	}
}

func (p *parser) feed(t *token) {
	if t.typ == tokenNewline {
		p.line = p.line + len(t.val)
		p.examine()
		p.stack = []*token{}
	} else {
		t.line = p.line
		p.stack = append(p.stack, t)
	}
}

func trimJoinBody(lines []string) string {
  if len(lines) == 0 {
    return ""
  }
	i := 0
	found := false
	for !found {
		r := string(lines[0][i])
		switch r {
		case " ", "\t":
			i = i + 1
		default:
			found = true
		}
	}
	ls := []string{}
	for _, line := range lines {
		ls = append(ls, line[i:])
	}
	return strings.Join(ls, "\n")
}

func (p *parser) finalize() {
	for _, t := range p.module.targets {
		p.module.targetIndex[t.Name()] = t
	}
}