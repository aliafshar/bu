package bu

import (
	"os"
)

type dependency interface {
	can(rt *runtime) bool
	resolve(rt *runtime) *target
}

type targetDependency struct {
	name string
}

func (d *targetDependency) can(rt *runtime) bool {
	return rt.history.done(d.name)
}

func (d *targetDependency) resolve(r *runtime) *target {
	return r.script.Target(d.name)
}

type fileDependency struct {
	filename string
}

func (d *fileDependency) can(rt *runtime) bool {
	_, err := os.Stat(d.filename)
	return err == nil
}

func (d *fileDependency) resolve(r *runtime) *target {
	return nil
}

type webDependency struct {
	uri string
}
