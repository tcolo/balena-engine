// TODO: do we need to handle .wh..wh.plnk layer hardlinks?
package aufsutil

import (
	"bufio"
	"errors"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/docker/docker/cmd/a2o-migrate/osutil"
	"github.com/docker/docker/pkg/archive"
)

var (
	// ErrAuFSRootNotExists indicates the aufs root directory wasn't found
	ErrAuFSRootNotExists = errors.New("Aufs root doesn't exists")
)

// CheckRootExists checks for the aufs storage root directory
func CheckRootExists(engineDir string) error {
	root := filepath.Join(engineDir, "aufs")
	logrus.WithField("aufs_root", root).Debug("checking if aufs root exists")
	ok, err := osutil.Exists(root, true)
	if err != nil {
		return err
	}
	if !ok {
		return ErrAuFSRootNotExists
	}
	return nil
}

func IsWhiteout(filename string) bool {
	return strings.HasPrefix(filename, archive.WhiteoutPrefix)
}

func IsWhiteoutMeta(filename string) bool {
	return strings.HasPrefix(filename, archive.WhiteoutMetaPrefix)
}

func IsOpaqueParentDir(filename string) bool {
	return filename == archive.WhiteoutOpaqueDir
}

func StripWhiteoutPrefix(filename string) string {
	out := filename
	for IsWhiteout(out) && !IsWhiteoutMeta(out) {
		out = strings.TrimPrefix(out, archive.WhiteoutPrefix)
	}
	return out
}

// Read the layers file for the current id and return all the
// layers represented by new lines in the file
//
// If there are no lines in the file then the id has no parent
// and an empty slice is returned.
//
// from daemon/graphdriver/aufs/dirs.go
func GetParentIDs(root, id string) ([]string, error) {
	f, err := os.Open(path.Join(root, "layers", id))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var out []string
	s := bufio.NewScanner(f)

	for s.Scan() {
		if t := s.Text(); t != "" {
			out = append(out, s.Text())
		}
	}
	return out, s.Err()
}
