package config

import (
	"errors"
	"path/filepath"

	"github.com/adrg/xdg"
)

// expandPath expands the path to include the home directory if the path
// is prefixed with `~`. If it isn't prefixed with `~`, the path is
// returned as-is.
//
// https://github.com/mitchellh/go-homedir/blob/af06845cf3004701891bf4fdb884bfe4920b3727/homedir.go#L55-L77
func expandPath(path string) (string, error) {
	if len(path) == 0 {
		return path, nil
	}

	if path[0] != '~' {
		return path, nil
	}

	if len(path) > 1 && path[1] != '/' && path[1] != '\\' {
		return "", errors.New("cannot expand user-specific home dir")
	}

	return filepath.Join(xdg.Home, path[1:]), nil
}
