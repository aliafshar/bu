package bu

import (
	"net/http"
	"os"
	"strings"

	"github.com/aliafshar/toylog"
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

func (d *webDependency) can(r *runtime) bool {
	uri := d.uri
	if !strings.HasPrefix(uri, "http") {
		uri = "http://" + uri
	}
	resp, err := http.Get(uri)
	if err != nil {
		toylog.Errorln("error fetching", d.uri, err)
		return false
	}
	return resp.StatusCode == 200
}

func (d *webDependency) resolve(r *runtime) *target {
	return nil
}
