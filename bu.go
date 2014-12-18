package bu

import (
	"github.com/aliafshar/toylog"
	"sort"
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

func List(bufile string) []string {
	s, _ := Load(bufile)
	var ts []string
	for k := range s.targets {
		ts = append(ts, k)
	}
	sort.Strings(ts)
	return ts
}
