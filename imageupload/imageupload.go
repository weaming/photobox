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
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
)

const thumbTempName = "thumbnail.jpg"

type Image struct {
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
	Format      string `json:"format"`
	Size        int    `json:"size"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	Data        []byte `json:"-"`
	Path        string `json:"file"`
	Sha256      string `json:"sha256"`
}

// Save image to file.
func (i *Image) Save(filename string) error {
	absPath, err := filepath.Abs(filename)
	i.Path = absPath
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, i.Data, 0644)
	if err == nil {
		// update the temp name
		if i.Filename == thumbTempName {
			i.Filename = path.Base(filename)
		}
	}
	return err
}

// Convert image to base64 data uri.
func (i *Image) DataURI() string {
	return fmt.Sprintf("data:%s;base64,%s", i.ContentType, base64.StdEncoding.EncodeToString(i.Data))
}

// Write image to HTTP response.
func (i *Image) WriteResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", i.ContentType)
	w.Header().Set("Content-Length", strconv.Itoa(i.Size))
	_, err := w.Write(i.Data)
	if err != nil {
		log.Println(err)
	}
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

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	img, format, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	i := &Image{
		Filename:    info.Filename,
		ContentType: contentType,
		Format:      format,
		Data:        data,
		Size:        len(data),
		Width:       img.Bounds().Max.X,
		Height:      img.Bounds().Max.Y,
		Sha256:      Sha256(data),
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
