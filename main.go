package main

import (
	"flag"
	"path"

	"github.com/gin-gonic/gin"
)

var DataDir = "photos"

func init() {
	flag.StringVar(&DataDir, "data", DataDir, "directory to save photos")
	flag.Parse()
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
