package main

import (
	"fmt"
	"net/http"
	"path"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	libfs "github.com/weaming/golib/fs"
	"github.com/weaming/golib/utils"
	"github.com/weaming/photobox/imageupload"
)

type UploadResponse struct {
	Image *imageupload.Image `json:"image"`
	Thumb *imageupload.Image `json:"thumb"`
	PU    *PathUrl           `json:"data"`
}

func Upload(c *gin.Context) {
	img, err := imageupload.Process(c.Request, "file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	width, height, quality := getThumbParams(c)
	t, err := imageupload.Thumbnail(img, width, height, quality)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pu := generateFilePath()

	fp := path.Join(DataDir, pu.OriginPath)
	err = saveImg(fp, img, c)
	if err != nil {
		return
	}

	fp = path.Join(DataDir, pu.ThumbPath)
	err = saveImg(fp, t, c)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, UploadResponse{img, t, &pu})
}

func saveImg(fp string, img *imageupload.Image, c *gin.Context) error {
	err := libfs.CreateDirIfNotExist(path.Dir(fp), false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return err
	}

	err = img.Save(fp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return err
	}
	return nil
}

func Thumbnail(c *gin.Context) {
	img, err := imageupload.Process(c.Request, "file")
	panicErr(err)

	width, height, quality := getThumbParams(c)
	t, err := imageupload.Thumbnail(img, width, height, quality)
	t.Write(c.Writer)
}

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}

type PathUrl struct {
	OriginPath string `json:"path"`
	ThumbPath  string `json:"thumb_path"`
	OriginURL  string `json:"url"`
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
