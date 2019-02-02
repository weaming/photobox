package storage

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/weaming/photobox/imageupload"
)

type S3Storage struct {
	Img *imageupload.Image
}

func (s *S3Storage) Save(fp string) (err error) {
	err = errors.New("not uploaded yet")
	count := 3
	bucket := os.Getenv("PHOTOBOX_BUCKET")
	if bucket == "" {
		bucket = "photobox-develop"
	}
	for {
		err = S3Upload(bucket, fp, s.Img.Data)
		count -= 1
		if err == nil || count <= 0 {
			break
		}
	}
	return
}

func (s *S3Storage) Read(fp string) ([]byte, error) {
	return ioutil.ReadFile(fp)
}

func S3Upload(bucket, key string, data []byte) error {
	region := os.Getenv("AWS_DEFAULT_REGION")
	if region == "" {
		region = "us-east-2"
	}
	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String(region)}))
	uploader := s3manager.NewUploader(sess)

	buf := bytes.NewBuffer(data)
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   buf,
		// extra options
		ContentType: aws.String("image/png"),
	})
	if err == nil {
		log.Printf("file uploaded to %s\n", aws.StringValue(&result.Location))
	}
	return err
}
