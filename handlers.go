package main

import (
	"log"
	"net/http"
	"path"
	"strconv"

	"github.com/gin-gonic/gin"
	libfs "github.com/weaming/golib/fs"
	"github.com/weaming/photobox/imageupload"
	"github.com/weaming/photobox/storage"
)

func Upload(c *gin.Context) {
	img, err := imageupload.Process(c.Request, "file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// check cache
	cacheRes := UploadResponse{}
	err = CacheGet(img.Sha256, &cacheRes)
	if err == nil {
		if libfs.Exist(cacheRes.Image.Path) && libfs.Exist(cacheRes.Thumb.Path) {
			log.Printf("hit cache %v", img.Sha256)
			c.JSON(http.StatusOK, cacheRes)
			return
		}
	}

	width, height, quality := getThumbParams(c)
	t, err := imageupload.Thumbnail(img, width, height, quality)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pu := generateFilePath(img.Sha256)

	fp := path.Join(DataDir, pu.OriginPath)
	err = saveImage(fp, img)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fp = path.Join(DataDir, pu.ThumbPath)
	err = saveImage(fp, t)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// cache and return
	res := UploadResponse{img, t, &pu}
	CacheSet(img.Sha256, &res)
	c.JSON(http.StatusOK, res)
}

func saveImage(fp string, img *imageupload.Image) error {
	local := storage.LocalStorage{Img: img}
	return storage.SaveTo(&local, fp)
}

func Thumbnail(c *gin.Context) {
	img, err := imageupload.Process(c.Request, "file")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	width, height, quality := getThumbParams(c)
	t, err := imageupload.Thumbnail(img, width, height, quality)
	t.Write(c.Writer)
}

func getThumbParams(c *gin.Context) (int, int, int) {
	width := defaultMaxWidth
	height := defaultMaxHeight
	quality := defaultQuality
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
