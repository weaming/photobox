package fs

import (
	"os"
	"path/filepath"
)

func Exist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func CreateDirIfNotExist(dir string, force bool) error {
	fi, err := os.Stat(dir)
	if os.IsNotExist(err) {
		// prepare dir
		if err := os.MkdirAll(dir, 0777); err != nil {
			return err
		}
	} else {
		// if is normal file
		if force && fi.Mode().IsRegular() {
			if err := os.Remove(dir); err != nil {
				return err
			}
		}
	}
	return nil
}

func OpenFile(path string) (*os.File, error) {
	path = ExpandUser(path)
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0777)
		if err != nil {
			return nil, err
		}
	}
	return os.Create(path)
}
