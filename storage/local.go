package storage

import (
	"io/ioutil"
	"path"

	libfs "github.com/weaming/golib/fs"
	"github.com/weaming/photobox/imageupload"
)

type LocalStorage struct {
	Img *imageupload.Image
}

func (l *LocalStorage) Save(fp string) error {
	err := libfs.CreateDirIfNotExist(path.Dir(fp), false)
	if err != nil {
		return err
	}

	return l.Img.Save(fp)
}
func (l *LocalStorage) Read(fp string) ([]byte, error) {
	return ioutil.ReadFile(fp)
}
