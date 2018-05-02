package imageupload

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"net/http"
	"path"
	"strconv"
	"strings"

	"github.com/nfnt/resize"
)

const thumbTempName = "thumbnail.jpg"

type Image struct {
	Filename    string `json:"filename"`
	ContentType string `json:"content-type"`
	Format      string `json:"format"`
	Size        int    `json:"size"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	Data        []byte `json:"-"`
}

// Save image to file.
func (i *Image) Save(filename string) error {
	err := ioutil.WriteFile(filename, i.Data, 0644)
	if err == nil && i.Filename == thumbTempName {
		i.Filename = path.Base(filename)
	}
	return err
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
	}

	return i, nil
}

// Create JPEG thumbnail.
func ThumbnailJPEG(i *Image, width int, height int, quality int) (*Image, error) {
	img, format, err := image.Decode(bytes.NewReader(i.Data))

	thumbnail := resize.Thumbnail(uint(width), uint(height), img, resize.Lanczos3)

	data := new(bytes.Buffer)
	err = jpeg.Encode(data, thumbnail, &jpeg.Options{
		Quality: quality,
	})

	if err != nil {
		return nil, err
	}

	bs := data.Bytes()

	t := &Image{
		Filename:    thumbTempName,
		ContentType: "image/jpeg",
		Format:      format,
		Data:        bs,
		Size:        len(bs),
		Width:       thumbnail.Bounds().Max.X,
		Height:      thumbnail.Bounds().Max.Y,
	}

	return t, nil
}

// Create PNG thumbnail.
func ThumbnailPNG(i *Image, width int, height int) (*Image, error) {
	img, format, err := image.Decode(bytes.NewReader(i.Data))

	thumbnail := resize.Thumbnail(uint(width), uint(height), img, resize.Lanczos3)

	data := new(bytes.Buffer)
	err = png.Encode(data, thumbnail)

	if err != nil {
		return nil, err
	}

	bs := data.Bytes()

	t := &Image{
		Filename:    thumbTempName,
		ContentType: "image/png",
		Format:      format,
		Data:        bs,
		Size:        len(bs),
		Width:       thumbnail.Bounds().Max.X,
		Height:      thumbnail.Bounds().Max.Y,
	}

	return t, nil
}

func Thumbnail(img *Image, width, height, quality int) (*Image, error) {
	if strings.HasSuffix(strings.ToLower(img.Filename), ".png") {
		return ThumbnailPNG(img, width, height)
	} else {
		return ThumbnailJPEG(img, width, height, quality)
	}
}
