package bu

import (
  "os"
)


type dependency interface {
	isDone(w *worker) bool
	runnable(s *script) target
}

type targetDependency struct {
	name string
}

func (d *targetDependency) isDone(w *worker) bool {
	return w.q.hasDone(d.name)
}

func (d *targetDependency) runnable(s *script) target {
	return s.Target(d.name)
}

type fileDependency struct {
	filename string
}

func (d *fileDependency) isDone(w *worker) bool {
  _, err := os.Stat(d.filename)
  return err == nil
}

func (d *fileDependency) runnable(s *script) target {
	return nil
}
