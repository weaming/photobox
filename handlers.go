package main

import (
	"fmt"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	libfs "github.com/weaming/golib/fs"
	"github.com/weaming/golib/utils"
	"github.com/weaming/photobox/imageupload"
)

var (
	ImageURLPrefix = "http://photobox.bitsflow.org"
)

type UploadResponse struct {
	Image *imageupload.Image `json:"data"`
	PU    PathUrl            `json:"url"`
}

func Upload(c *gin.Context) {
	img, err := imageupload.Process(c.Request, "file")
	panicErr(err)

	width, height, quality := getThumbParams(c)
	t, err := CommonThumb(img, width, height, quality)
	panicErr(err)

	pu := generateFilePath()

	fp := path.Join(DataDir, pu.OriginPath)
	panicErr(libfs.CreateDirIfNotExist(path.Dir(fp), false))
	panicErr(img.Save(fp))

	fp = path.Join(DataDir, pu.ThumbPath)
	panicErr(libfs.CreateDirIfNotExist(path.Dir(fp), false))
	panicErr(t.Save(fp))

	c.JSON(http.StatusOK, UploadResponse{img, pu})
}

func CommonThumb(img *imageupload.Image, width, height, quality int) (*imageupload.Image, error) {
	if strings.HasSuffix(strings.ToLower(img.Filename), ".png") {
		return imageupload.ThumbnailPNG(img, width, height)
	} else {
		return imageupload.ThumbnailJPEG(img, width, height, quality)
	}
}

func Thumbnail(c *gin.Context) {
	img, err := imageupload.Process(c.Request, "file")
	panicErr(err)

	width, height, quality := getThumbParams(c)
	t, err := CommonThumb(img, width, height, quality)
	t.Write(c.Writer)
}

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}

type PathUrl struct {
	OriginPath string `json:"origin_path"`
	ThumbPath  string `json:"thumb_path"`
	OriginURL  string `json:"origin_url"`
	ThumbURL   string `json:"thumb_url"`
}

func generateFilePath() PathUrl {
	now := time.Now()
	suffix := path.Join(utils.FormatDateTime("%Y/%m/%d", now), fmt.Sprintf("%d.png", now.Unix()))
	o, t := path.Join("origin", suffix), path.Join("thumb", suffix)
	return PathUrl{
		o, t,
		path.Join(ImageURLPrefix, o),
		path.Join(ImageURLPrefix, t),
	}
}

func getThumbParams(c *gin.Context) (int, int, int) {
	width := 300
	height := 300
	quality := 80
	if str, ok := c.GetQuery("width"); ok {
		if v, e := strconv.Atoi(str); e == nil {
			width = v
		}
	}
	if str, ok := c.GetQuery("height"); ok {
		if v, e := strconv.Atoi(str); e == nil {
			height = v
		}
	}
	if str, ok := c.GetQuery("quality"); ok {
		if v, e := strconv.Atoi(str); e == nil {
			quality = v
		}
	}
	return width, height, quality
}
