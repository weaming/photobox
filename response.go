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
	Data  *ImagePathUrl      `json:"data"`
}

type ImagePathUrl struct {
	OriginPath string `json:"path"`
	ThumbPath  string `json:"thumb_path"`
	OriginURL  string `json:"url"`
	ThumbURL   string `json:"thumb_url"`
}

func generateImagePathUrl(hash, format string) ImagePathUrl {
	now := time.Now()
	if format == "jpeg" {
		format = "jpg"
	}
	suffix := path.Join(utils.FormatDateTime("%Y/%m/%d", now), fmt.Sprintf("%d_%v.%v", now.Unix(), hash, format))
	o, t := path.Join("origin", suffix), path.Join("thumb", suffix)
	return ImagePathUrl{
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
