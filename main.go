package main

import (
	"flag"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	DataDir        = "photos"
	ImageURLPrefix = "http://photobox.bitsflow.org"
)

func init() {
	flag.StringVar(&DataDir, "data", DataDir, "directory to save photos")
	flag.StringVar(&ImageURLPrefix, "host", ImageURLPrefix, "photos storage host name")
	flag.Parse()
	if !strings.HasPrefix(ImageURLPrefix, "http") {
		ImageURLPrefix = "http://" + ImageURLPrefix
	}
}

func main() {
	r := gin.Default()
	r.Static("/origin", path.Join(DataDir, "origin"))
	r.Static("/thumb", path.Join(DataDir, "thumb"))

	r.GET("/", func(c *gin.Context) {
		c.File("index.html")
	})

	r.POST("/upload", Upload)
	r.POST("/thumbnail", Thumbnail)

	r.Run(":5000")
}
