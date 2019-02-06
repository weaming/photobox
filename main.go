package main

import (
	"flag"
	"log"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/weaming/photobox/storage"
)

var (
	DataDir          = "photos"
	ImageDomain      = "http://photobox.bitsflow.org"
	listen           = ":5000"
	redisHostPort    = "127.0.0.1:6379"
	redisDB          = 8
	defaultMaxWidth  = 1280
	defaultMaxHeight = 1280
	defaultQuality   = 80
)

func init() {
	flag.StringVar(&DataDir, "data", DataDir, "directory to save photos")
	flag.StringVar(&ImageDomain, "domain", ImageDomain, "photos storage domain name")
	flag.StringVar(&listen, "listen", listen, "bind [<host>]:<port>")
	flag.StringVar(&redisHostPort, "redis", redisHostPort, "redis <host>:<port>")
	flag.IntVar(&redisDB, "db", redisDB, "redis database")
	flag.Parse()
	if !strings.HasPrefix(ImageDomain, "http") {
		ImageDomain = "http://" + ImageDomain
	}

	initCodec()
}

func main() {
	r := gin.Default()
	r.StaticFS("/origin", storage.NewLocalS3FS(path.Join(DataDir, "origin"), "/origin"))
	r.StaticFS("/thumb", storage.NewLocalS3FS(path.Join(DataDir, "thumb"), "/thumb"))

	r.GET("/", func(c *gin.Context) {
		c.File("index.html")
	})

	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"success": false, "message": "resource not found"})
	})

	r.POST("/upload", APIUpload)
	r.POST("/thumbnail", APIThumbnail)

	log.Fatal(r.Run(listen))
}
