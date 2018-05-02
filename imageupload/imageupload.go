package imageupload

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"sync"
)

const thumbTempName = "thumbnail.jpg"

var hashPathMap = map[string]*Image{}
var mapLock = sync.Mutex{}

type Image struct {
	Filename    string `json:"filename"`
	ContentType string `json:"content-type"`
	Format      string `json:"format"`
	Size        int    `json:"size"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	Data        []byte `json:"-"`
	Path        string `json:"-"`
	Sha256      string `json:"sha256"`
}

// Save image to file.
func (i *Image) Save(filename string) (*Image, error) {
	i.Path, _ = filepath.Abs(filename)

	/*
			if old, ok := hashPathMap[hash]; ok {
				// load old data
				if ExistFile(old.Path) {
					dat, err := ioutil.ReadFile(old.Path)
					if err != nil {
						goto SAVE
					}
					old.Data = dat
					return old, nil
				}
			}

		SAVE:
	*/
	err := ioutil.WriteFile(filename, i.Data, 0644)
	if err == nil {
		// update the temp name
		if i.Filename == thumbTempName {
			i.Filename = path.Base(filename)
		}

	}

	return i, err
}

// Convert image to base64 data uri.
func (i *Image) DataURI() string {
	return fmt.Sprintf("data:%s;base64,%s", i.ContentType, base64.StdEncoding.EncodeToString(i.Data))
}

// Write image to HTTP response.
func (i *Image) Write(w http.ResponseWriter) {
	w.Header().Set("Content-Type", i.ContentType)
	w.Header().Set("Content-Length", strconv.Itoa(i.Size))
	w.Write(i.Data)
}

// Create JPEG thumbnail from image.
func (i *Image) ThumbnailJPEG(width int, height int, quality int) (*Image, error) {
	return ThumbnailJPEG(i, width, height, quality)
}

// Create PNG thumbnail from image.
func (i *Image) ThumbnailPNG(width int, height int) (*Image, error) {
	return ThumbnailPNG(i, width, height)
}

// Limit the size of uploaded files, put this before imageupload.Process.
func LimitFileSize(maxSize int64, w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxSize)
}

func okContentType(contentType string) bool {
	return contentType == "image/png" || contentType == "image/jpeg" || contentType == "image/gif"
}

// Process uploaded file into an image.
func Process(r *http.Request, field string) (*Image, error) {
	file, info, err := r.FormFile(field)

	if err != nil {
		return nil, err
	}

	contentType := info.Header.Get("Content-Type")

	if !okContentType(contentType) {
		return nil, errors.New(fmt.Sprintf("Wrong content type: %s", contentType))
	}

	bs, err := ioutil.ReadAll(file)

	if err != nil {
		return nil, err
	}

	img, format, err := image.Decode(bytes.NewReader(bs))

	if err != nil {
		return nil, err
	}

	i := &Image{
		Filename:    info.Filename,
		ContentType: contentType,
		Format:      format,
		Data:        bs,
		Size:        len(bs),
		Width:       img.Bounds().Max.X,
		Height:      img.Bounds().Max.Y,
		Sha256:      Sha256(bs),
	}

	return i, nil
}

func ExistFile(fp string) bool {
	if _, err := os.Stat(fp); err == nil {
		return true
	}
	return false
}

func Sha256(content []byte) string {
	return fmt.Sprintf("%x", sha256.Sum256(content))
}
