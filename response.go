package main

import (
	"fmt"
	"log"
	"net/url"
	"path"
	"time"

	"github.com/weaming/golib/utils"
	"github.com/weaming/photobox/imageupload"
)

type UploadResponse struct {
	Image *imageupload.Image `json:"image"`
	Thumb *imageupload.Image `json:"thumb"`
	PU    *PathUrl           `json:"data"`
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
		joinURI(ImageDomain, o),
		joinURI(ImageDomain, t),
	}
}

func joinURI(base, uri string) string {
	u, err := url.Parse(uri)
	if err != nil {
		log.Fatal(err)
	}
	b, err := url.Parse(base)
	if err != nil {
		log.Fatal(err)
	}
	return b.ResolveReference(u).String()
}
