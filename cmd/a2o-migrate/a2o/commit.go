package a2o

import (
	"path/filepath"

	"github.com/sirupsen/logrus"
)

// Commit finalises the migration by deleting leftover data.
func Commit() error {
	logrus.Info("committing changes")

	// remove aufs layer data
	err := removeDirIfExists(aufsRoot())
	if err != nil {
		return err
	}

	// remove images
	aufsImageDir := filepath.Join(StorageRoot, "image", "aufs")
	err = removeDirIfExists(aufsImageDir)
	if err != nil {
		return err
	}

	return nil
}
