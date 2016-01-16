package shooter

import (
	"crypto/md5"
	"fmt"
	"os"
	"strings"
)

func FileHash(name string) (string, error) {
	file, err := os.Open(name)
	if err != nil {
		return "", err
	}
	defer file.Close()
	stat, err := file.Stat()
	if err != nil {
		return "", err
	}
	size := stat.Size()
	if size < 4096*4 {
		return "", fmt.Errorf("File size too small.")
	}
	b := make([]byte, 4096)
	ss := make([]string, 0, 4)
	var readChunk = func(off int64) error {
		if n, err := file.ReadAt(b, off); err != nil {
			return err
		} else if n != 4096 {
			return fmt.Errorf("Partial read: %d", n)
		}
		ss = append(ss, fmt.Sprintf("%x", md5.Sum(b)))
		return nil
	}

	if err = readChunk(4096); err != nil {
		return "", err
	}
	if err = readChunk(size / 3 * 2); err != nil {
		return "", err
	}
	if err = readChunk(size / 3); err != nil {
		return "", err
	}
	if err = readChunk(size - 8192); err != nil {
		return "", err
	}
	return strings.Join(ss, ";"), nil
}
