package fs

import (
	"os"
	"strconv"
)

func ArchiveBigFile(fp string, maxSize int64, n uint) (uint, error) {
	// n must be given 1 when call manually
	maxBackup := n
	if fi, err := os.Stat(fp); os.IsNotExist(err) {
		return 0, nil
	} else {
		// check and call recursively
		if fi.Size() >= maxSize {
			newFilePath := fp + "." + strconv.Itoa(int(n))
			_, err := os.Stat(newFilePath)
			if !os.IsNotExist(err) {
				// call self recursively
				maxBackup, err = ArchiveBigFile(fp, maxSize, n+1)
				if err != nil {
					// todo: test error
					return maxBackup, err
				}
			}

			// rename to new name
			realLastFilePath := fp
			if n > 1 {
				realLastFilePath = fp + "." + strconv.Itoa(int(n-1))
			}
			err = os.Rename(realLastFilePath, newFilePath)
			if err != nil {
				return maxBackup, err
			}
		}
	}

	return maxBackup, nil
}
