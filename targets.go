package bu

type target interface {
	Name() string
	Deps() []target
	Run()
	Finalize(*script)
}

type questionTarget struct {
	name      string
	dflt      string
	usage     string
	bodyLines []string
}

func (t *questionTarget) Name() string {
	return t.name
}

func (t *questionTarget) Deps() []target {
	return nil
}

func (t *questionTarget) Finalize(s *script) {
	t.usage = trimJoinBody(t.bodyLines)
}

func (t *questionTarget) AppendBody(s string) {
	t.bodyLines = append(t.bodyLines, s)
}

type shellTarget struct {
	name      string
	body      string
	bodyLines []string
	depsNames []string
	deps      []target
	typ       string
	args      []string
}

func (t *shellTarget) Name() string {
	return t.name
}

func (t *shellTarget) Deps() []target {
	return t.deps
}

func (t *shellTarget) AppendBody(s string) {
	t.bodyLines = append(t.bodyLines, s)
}

func (t *shellTarget) Finalize(s *script) {
	t.body = trimJoinBody(t.bodyLines)
	t.args = s.args
	for _, dn := range t.depsNames {
		if d := s.Target(dn); d != nil {
			t.deps = append(t.deps, d)
		} else {
			dependencyError(dn, t)
			return
		}
	}
}
