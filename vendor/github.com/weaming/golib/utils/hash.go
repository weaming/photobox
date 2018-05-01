package utils

import (
	"crypto/md5"
	"fmt"
)

func MD5(input []byte) string {
	digest := md5.New()
	digest.Write(input)
	sum := digest.Sum(nil)
	return fmt.Sprintf("%X", sum)
}
