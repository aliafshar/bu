package bu

import (
	"github.com/aliafshar/toylog"
)

func parseError(msg string, t *token) {
	toylog.Errorf("Error parsing. "+msg+" [%q]", t.value())
}

func dependencyError(depName string, t target) {
	toylog.Errorf("Missing dependency %q in %q.", depName, t.Name())
}

func logicalError(msg string, n *node) {
	toylog.Errorf("Error understanding. "+msg+" [%+v]", n)
}
