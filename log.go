package bu

import (
	"github.com/aliafshar/toylog"
)

func parseError(msg string, t *token) {
	toylog.Errorf("Error parsing. "+msg+" [line %v %q]", t.line, t.value())
}

func dependencyError(depName string, t target) {
	toylog.Errorf("Missing dependency %q in %q.", depName, t.Name())
}
