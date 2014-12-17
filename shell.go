package bu

import (
	"os"
	"os/exec"
)

var shells = map[string]string{
	"sh": "bash",
	"py": "python",
}

func (t *target) cmd(r *runtime) *exec.Cmd {
	args := append([]string{"-c", t.body, t.name}, r.argv...)
	cmd := exec.Command(shells[t.shell], args...)
	cmd.Stderr = os.Stderr
	cmd.Env = r.env
	return cmd
}
