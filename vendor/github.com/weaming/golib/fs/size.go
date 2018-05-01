package fs

import (
	"fmt"
	"os"
)

func HumanSize(size int64) string {
	const ratio = 1024
	size_float := float64(size)
	units := []string{"B", "KB", "MB", "GB", "TB", "EB"}

	index := 0
	for ; size_float > ratio; index += 1 {
		size_float /= ratio
	}
	return fmt.Sprintf("%.2f %s", size_float, units[index])
}

func FileSize(path string) int64 {
	fileInfo, err := os.Lstat(path)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	return fileInfo.Size()
}

func FilesSize(files []string) (total int64) {
	for _, path := range files {
		total += FileSize(path)
	}
	return
}

func FileHumanSize(path string) string {
	return HumanSize(FileSize(path))
}
