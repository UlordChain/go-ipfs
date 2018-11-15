package fsrepo

import (
	"os"

	homedir "github.com/udfs/go-udfs/Godeps/_workspace/src/github.com/mitchellh/go-homedir"
	"github.com/udfs/go-udfs/repo/config"
)

// BestKnownPath returns the best known fsrepo path. If the ENV override is
// present, this function returns that value. Otherwise, it returns the default
// repo path.
func BestKnownPath() (string, error) {
	udfsPath := config.DefaultPathRoot
	if os.Getenv(config.EnvDir) != "" {
		udfsPath = os.Getenv(config.EnvDir)
	}
	udfsPath, err := homedir.Expand(udfsPath)
	if err != nil {
		return "", err
	}
	return udfsPath, nil
}
