package bu

import (
	"github.com/aliafshar/toylog"
	"os"
	"os/user"
	"path/filepath"
)

type module struct {
	targets     []target
	targetIndex map[string]target
	setvars     []*setvar
	imports     []string
}

type setvar struct {
	key       string
	bodyLines []string
}

func (t *setvar) AppendBody(s string) {
	t.bodyLines = append(t.bodyLines, s)
}

func (t *setvar) value() string {
	return trimJoinBody(t.bodyLines)
}

type script struct {
	modules []*module
	setvars []*setvar
	path    []string
	args    []string
	targets map[string]target
}

func (s *script) Target(name string) target {
	return s.targets[name]
}

func homeFilename(filename string) (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(usr.HomeDir, filename), nil
}

func defaultPath() []string {
	path := []string{"."}
	home, err := homeFilename(".bu")
	if err == nil {
		path = append(path, home)
	}
	return path
}

func (s *script) finalize() {
	mro := append(s.modules[1:], s.modules[0])
	for _, m := range mro {
		s.setvars = append(s.setvars, m.setvars...)
		for _, t := range m.targets {
			s.targets[t.Name()] = t
		}
	}
	for _, t := range s.targets {
		t.Finalize(s)
	}
}

func (s *script) loadModule(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		toylog.Errorln("File not loaded.", filename, err)
		return
	}
	defer f.Close()
	l := NewLexer(f)
	p := newParser()
	p.parse(l)
	s.modules = append(s.modules, p.module)
	for _, i := range p.module.imports {
		filename := s.resolveModule(i)
		if filename == "" {
			toylog.Errorln("Unable to find module", i)
			continue
		}
		s.loadModule(filename)
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

func NewScript(filename string, args []string) *script {
	s := &script{path: defaultPath(), args: args, targets: make(map[string]target)}
	s.loadModule(filename)
	s.finalize()
	return s
}
