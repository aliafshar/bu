package bu

import (
  "io"
	"github.com/aliafshar/toylog"
	"os"
	"path/filepath"
)

type script struct {
	parser   *parser
	modules  []*module
	module   *module
	setvars  []*setvar
	path     []string
	args     []string
	targets  map[string]*target
	filename string
}

type module struct {
	targets []*target
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
	return &script{
		parser:  newParser(),
		targets: make(map[string]*target),
		// TODO move out of script
		path: defaultPath(),
	}
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
	return s.load(f)
}

func (s *script) load(r io.Reader) (*module, error) {
	return s.parser.parse(s, r)
}

func (s *script) Target(name string) *target {
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
			s.targets[t.name] = t
		}
		s.setvars = append(s.setvars, m.setvars...)
	}
}


func Load(r io.Reader, filename string, args ...string) (*script, error) {
	s := newScript()
	s.filename = filename
	s.args = args
	m, err := s.load(r)
	if err != nil {
		return nil, err
	}
	for _, imp := range m.imports {
		path := s.resolveModule(imp.key)
		if path == "" {
			toylog.Errorf("unable to resolve module %q", imp.key)
			continue
		}
		s.loadFile(path)
	}
	s.aggregate()
	return s, nil
}
