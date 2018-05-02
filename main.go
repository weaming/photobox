package main

import (
	"flag"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	DataDir     = "photos"
	ImageDomain = "http://photobox.bitsflow.org"
	listen      = ":5000"
)

func init() {
	flag.StringVar(&DataDir, "data", DataDir, "directory to save photos")
	flag.StringVar(&ImageDomain, "domain", ImageDomain, "photos storage domain name")
	flag.StringVar(&listen, "listen", listen, "bind [<host>]:<port>")
	flag.Parse()
	if !strings.HasPrefix(ImageDomain, "http") {
		ImageDomain = "http://" + ImageDomain
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

	r.Run(listen)
}
