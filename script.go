package bu

import (
	"os"
  "fmt"
  "path/filepath"
  "github.com/aliafshar/toylog"
  "os/user"
)

type module struct {
	targets     []*target
	targetIndex map[string]*target
	setvars     []*setvar
  imports []string
}

func (s *module) Target(name string) *target {
	if name == "" {
		return s.targets[0]
	}
	return s.targetIndex[name]
}

type target struct {
	name      string
	body      string
	bodyLines []string
	depsNames []string
	deps      []*target
	typ       string
	module    *module
  tokens []*token
}

type setvar struct {
	key   string
	value string
}

type script struct {
  modules []*module
  setvars []*setvar
  path []string
}

func (s *script) Target(name string) *target {
  for _, m := range s.modules {
    if t, ok := m.targetIndex[name]; ok {
      return t
    }
  }
  return nil
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
    for _, t := range m.targets {
      s.setvars = append(s.setvars, m.setvars...)
      for _, dn := range t.depsNames {
        if d := s.Target(dn); d != nil {
          t.deps = append(t.deps, d)
        } else {
          parseError(fmt.Sprintf("Missing dependency %q.", t.name), t.tokens[0])
          return
        }
      }
    }
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
  for _, path := range(s.path) {
    filename := filepath.Join(path, name)
    if _, err := os.Stat(filename); err == nil {
      return filename
    }
  }
  return ""
}

func NewScript(filename string) *script {
  s := &script{modules: []*module{}, setvars: []*setvar{}, path: defaultPath()}
  s.loadModule(filename)
  s.finalize()
  return s
}
