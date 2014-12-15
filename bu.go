package bu

import (
	"github.com/aliafshar/toylog"
)

func Run(bufile, targetName string, args ...string) {
	s, _ := Load(bufile, args...)
	r := newRuntime(s)
	t := s.Target(targetName)
	if t == nil {
		toylog.Fatalf("target not found %q", targetName)
	}
	r.build(t)
	r.start()
}
