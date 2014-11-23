package bu

import (
	"fmt"
	"os"
	"path/filepath"
)

type script struct {
	parser  *parser
	modules []*module
	module  *module
	setvars []*setvar
	path    []string
	args    []string
	targets map[string]target
}

type module struct {
	targets []target
	setvars []*setvar
	imports []*imports
}

type setvar struct {
	key   string
	value string
}

type imports struct {
	key string
}

func newScript() *script {
	return &script{parser: newParser(), targets: make(map[string]target), path: defaultPath()}
}

func (s *script) resolveModule(name string) string {
	for _, path := range s.path {
		filename := filepath.Join(path, name)
		if _, err := os.Stat(filename); err == nil {
			return filename
		}
	}
	return ""
}

func (s *script) loadFile(filename string) (*module, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return s.parser.parse(s, f)
}

func (s *script) Target(name string) target {
	if name == "" {
		return s.module.targets[0]
	}
	return s.targets[name]
}

func (s *script) mro() []*module {
	return append(s.modules[1:], s.module)
}

func (s *script) aggregate() {
	for _, m := range s.mro() {
		for _, t := range m.targets {
			s.targets[t.Name()] = t
		}
		s.setvars = append(s.setvars, m.setvars...)
	}
}

func Load(filename string, args ...string) (*script, error) {
	s := newScript()
	s.args = args
	m, err := s.loadFile(filename)
	if err != nil {
		fmt.Println(err)
	}
	for _, imp := range m.imports {
		path := s.resolveModule(imp.key)
		if path == "" {
			fmt.Println("Error, unable to resolve module")
			continue
		}
		s.loadFile(path)
	}
	s.aggregate()
	return s, nil
}
