package bu

import (
	"os/exec"
)

var shells = map[string]string{
	"sh": "bash",
	"py": "python",
}

func (t *target) cmd(r *runtime) *exec.Cmd {
	args := append([]string{"-c", t.body, t.name}, r.argv...)
	return exec.Command(shells[t.shell], args...)
}
