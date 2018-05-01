package fs

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Dir struct {
	Root string

	Dirs   []string
	Files  []string
	Images []string

	AbsDirs   []string
	AbsFiles  []string
	AbsImages []string
}

func NewDir(path string) *Dir {
	fi, err := os.Stat(path)
	if err != nil {
		panic(err)
	}

	if !fi.IsDir() {
		return nil
	} else {
		dir := Dir{Root: path}
		files, _ := ioutil.ReadDir(path)
		for _, fi := range files {
			name := fi.Name()
			absPath, err := filepath.Abs(fi.Name())
			if err != nil {
				panic(err)
			}

			if fi.IsDir() {
				dir.Dirs = append(dir.Dirs, name)
				dir.AbsDirs = append(dir.AbsDirs, absPath)
			} else {
				dir.Files = append(dir.Files, name)
				dir.AbsFiles = append(dir.AbsFiles, absPath)

				// list all photos
				switch strings.ToLower(filepath.Ext(name)) {
				case ".jpg", ".jpeg", ".png", ".gif", ".bmp":
					dir.Images = append(dir.Images, name)
					dir.AbsImages = append(dir.AbsImages, absPath)
				default:
				}
			}
		}
		return &dir
	}
}

func ContainsPhoto(path string) bool {
	dir := NewDir(path)
	if len(dir.Images) > 0 {
		return true
	} else {
		for _, subpath := range dir.AbsDirs {
			if ContainsPhoto(subpath) {
				return true
			}
		}
	}
	return false
}

func DirFilesSize(dirs []string) (total int64) {
	for _, path := range dirs {
		tmp := NewDir(path)
		total = total + FilesSize(tmp.AbsImages) + DirFilesSize(tmp.AbsDirs)
	}
	return
}
