package bu

type targetType func(*node) target

type target interface {
	Name() string
	Deps() []dependency
	Run()
}

type dependency interface {
	isDone(w *worker) bool
	runnable(s *script) target
}

type targetDependency struct {
	name   string
}

func (d *targetDependency) isDone(w *worker) bool {
	return w.q.hasDone(d.name)
}

func (d *targetDependency) runnable(s *script) target {
	return s.Target(d.name)
}

type fileDependency struct {
	name string
}

func (d *fileDependency) isDone(w *worker) bool {
	return true
}

func (d *fileDependency) runnable(s *script) target {
	return nil
}

type shell interface {
	execute(cmd string)
}

type bash struct{}

func execute(cmd string) {
}


type target2 interface {
	Run() error
	CanRun() bool
}

type questionTarget struct {
	name  string
	dflt  string
	usage string
}

func (t *questionTarget) Name() string {
	return t.name
}

func (t *questionTarget) Deps() []dependency {
	return nil
}

type shellTarget struct {
	name      string
	body      string
	depsNames []string
	deps      []dependency
	shell     string
	args      []string
}

func (t *shellTarget) Name() string {
	return t.name
}

func (t *shellTarget) Deps() []dependency {
	return t.deps
}

func (t *shellTarget) Shell() {
}

func newShellTarget(n *node) target {
	t := &shellTarget{name: n.key}
	return t
}

var targetTypes = map[tokenType]targetType{
	tokenColon: newShellTarget,
}
