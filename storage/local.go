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

func (s *LocalStorage) Save(fp string) error {
	err := libfs.CreateDirIfNotExist(path.Dir(fp), false)
	if err != nil {
		return err
	}

	return s.Img.Save(fp)
}
func (s *LocalStorage) Read(fp string) ([]byte, error) {
	return ioutil.ReadFile(fp)
}
