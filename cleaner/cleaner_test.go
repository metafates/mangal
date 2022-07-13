package cleaner

import (
	"fmt"
	"github.com/metafates/mangal/common"
	"github.com/metafates/mangal/filesystem"
	"github.com/spf13/afero"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func init() {
	filesystem.Set(afero.NewMemMapFs())
}

// generateFiles at the given directory with prefix appended to each file
func generateFiles(t *testing.T, count int, dir, prefix string) ([]string, error) {
	t.Helper()

	var files []string

	if err := filesystem.Get().MkdirAll(dir, 0777); err != nil {
		return nil, err
	}

	for i := 0; i < count; i++ {
		path := filepath.Join(dir, prefix+strconv.Itoa(rand.Intn(100000)))

		if _, err := filesystem.Get().Create(path); err != nil {
			return nil, err
		}
		files = append(files, path)
	}

	return files, nil
}

// existAll checks if all given files exist
func existAll(t *testing.T, files []string) (bool, error) {
	t.Helper()
	for _, file := range files {
		if exists, err := afero.Exists(filesystem.Get(), file); err != nil {
			return false, err
		} else if !exists {
			return false, nil
		}
	}

	return true, nil
}

func TestRemoveTemp(t *testing.T) {
	const count = 13

	files, err := generateFiles(t, count, os.TempDir(), common.TempPrefix)
	if err != nil {
		t.Fatal("could not create files for testing")
	}

	if exist, err := existAll(t, files); err != nil {
		t.Fatal("error while checking for files existence")
	} else if !exist {
		t.Fatal("files was not created")
	}

	removedCount, _ := RemoveTemp()

	if removedCount != count {
		fmt.Println(removedCount, count)
		t.Error("removed files count does not match expected count")
	}

	removedCount, _ = RemoveTemp()

	if removedCount != 0 {
		t.Error("repeated call expected to remove 0 files")
	}
}

func TestRemoveCache(t *testing.T) {
	RemoveCache()

	cacheDir, err := os.UserCacheDir()

	if err != nil {
		t.Fatal(err)
	}

	const count = 13

	files, err := generateFiles(t, count, cacheDir, "")
	if err != nil {
		t.Fatal("could not create files for testing", err)
	}

	if exist, err := existAll(t, files); err != nil {
		t.Fatal("error while checking for files existence")
	} else if !exist {
		t.Fatal("files was not created")
	}

	removedCount, _ := RemoveCache()

	if removedCount != 0 {
		t.Error("repeated call expected to remove 0 files")
	}
}
