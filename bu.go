package bu

import (
	"github.com/aliafshar/toylog"
	"io"
	"sort"
)

func Run(r io.Reader, filename, targetName string, args ...string) {
	s, _ := Load(r, filename, args...)
	ru := newRuntime(s)
	t := s.Target(targetName)
	toylog.Debugf("target is: %+v\n", t)
	if t == nil {
		toylog.Fatalf("target not found %q", targetName)
	}
	ru.run(t)
}

func List(r io.Reader, filename string) []string {
	s, _ := Load(r, filename)
	var ts []string
	for k := range s.targets {
		ts = append(ts, k)
	}
	sort.Strings(ts)
	return ts
}
