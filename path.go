package bu

import (
	"os/user"
	"path/filepath"
)

func homeFilename(filename string) (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(usr.HomeDir, filename), nil
}

func defaultPath() []string {
	path := []string{"."}
	home, err := homeFilename(".bu")
	if err == nil {
		path = append(path, home)
	}
	return path
}
