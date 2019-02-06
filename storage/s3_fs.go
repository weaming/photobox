package storage

import (
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type LocalS3FS struct {
	fs        http.FileSystem
	localRoot string
	bucket    string
	keyPrefix string
}

func PrepareDir(filePath string, force bool) {
	if !strings.HasSuffix(filePath, "/") || force {
		filePath = path.Dir(filePath)
	}
	err := os.MkdirAll(filePath, os.FileMode(0755))
	if err != nil {
		log.Fatal(err)
	}
}

func (r *LocalS3FS) Open(name string) (http.File, error) {
	log.Printf("reading file with name: %v\n", name)
	f, err := r.fs.Open(name)

	// file not found on disk, try get from S3
	if err != nil {
		key := path.Join("/", r.keyPrefix, name)
		log.Printf("reading s3 key: %v\n", key)
		data, err := S3Read(r.bucket, key)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		// save to local disk
		filePath := path.Join(r.localRoot, name)
		PrepareDir(filePath, false)
		err = ioutil.WriteFile(filePath, data, 0644)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		log.Printf("wrote local file from S3: %v\n", filePath)

		return &S3File{
			data:   data,
			bucket: r.bucket,
			key:    key,
		}, nil
	}
	return NeuteredReaddirFile{f}, nil
}

func NewLocalS3FS(root, keyPrefix string) *LocalS3FS {
	bucket := GetBucketName()
	fs := gin.Dir(root, false)
	return &LocalS3FS{
		fs:        fs,
		localRoot: root,
		bucket:    bucket,
		keyPrefix: keyPrefix,
	}
}

type S3File struct {
	bucket string
	key    string
	data   []byte
	offset int64

	http.File
	// type File interface {
	// 	io.Closer
	// 	io.Reader
	// 	io.Seeker
	// 	Readdir(count int) ([]os.FileInfo, error)
	// 	Stat() (os.FileInfo, error)
	// }
}

func (r *S3File) Read(p []byte) (n int, err error) {
	// log.Printf("reading %v %v %v\n", len(p), r.Size(), r.offset)
	n = copy(p, r.data[r.offset:len(r.data)])
	r.offset += int64(n)
	if n != len(p) {
		return n, errors.New("not fully copied")
	}
	return n, nil
}

var errWhence = errors.New("Seek: invalid whence")
var errOffset = errors.New("Seek: invalid offset")

func (r *S3File) Seek(offset int64, whence int) (int64, error) {
	var off int64
	switch whence {
	case io.SeekStart:
		off = offset
	case io.SeekCurrent:
		off = r.offset + offset
	case io.SeekEnd:
		off = int64(len(r.data)) + offset
	default:
		return 0, errWhence
	}
	if off < 0 {
		return 0, errOffset
	}
	r.offset = off
	// log.Printf("seeking %v %v %v\n", offset, whence, r.offset)
	return r.offset, nil
}

func (r *S3File) Stat() (os.FileInfo, error) {
	// type FileInfo interface {
	// 	Name() string       // base name of the file
	// 	Size() int64        // length in bytes for regular files; system-dependent for others
	// 	Mode() FileMode     // file mode bits
	// 	ModTime() time.Time // modification time
	// 	IsDir() bool        // abbreviation for Mode().IsDir()
	// 	Sys() interface{}   // underlying data source (can return nil)
	// }
	return r, nil
}

func (r *S3File) IsDir() bool        { return false }
func (r *S3File) Name() string       { return path.Base(r.key) }
func (r *S3File) Size() int64        { return int64(len(r.data)) }
func (r *S3File) Mode() os.FileMode  { return os.ModePerm }
func (r *S3File) ModTime() time.Time { return time.Now() }
func (r *S3File) Sys() interface{}   { return nil }

func (r *S3File) Close() error {
	return nil
}

type NeuteredReaddirFile struct {
	http.File
}

// Overrides the http.File default implementation
func (f NeuteredReaddirFile) Readdir(count int) ([]os.FileInfo, error) {
	// this disables directory listing
	return nil, nil
}
