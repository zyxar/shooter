package shooter

import (
	"crypto/md5"
	"fmt"
	"os"
	"strings"
)

// FileHash reads a file and calculates its hash.
// The file should be larger than 16 KB (reasonable for a normal film file).
// The hash can be used to query the subtitles in SHOOTER later.
func FileHash(name string) (str string, err error) {
	var file *os.File
	if file, err = os.Open(name); err != nil {
		return
	}
	defer file.Close()
	var stat os.FileInfo
	if stat, err = file.Stat(); err != nil {
		return
	}
	size := stat.Size()
	if size < 4096*4 {
		err = errFileSizeInsufficient
		return
	}
	b := make([]byte, 4096)
	ss := make([]string, 0, 4)
	var readChunk = func(off int64) error {
		if n, err := file.ReadAt(b, off); err != nil {
			return err
		} else if n != 4096 {
			return errPartialRead
		}
		ss = append(ss, fmt.Sprintf("%x", md5.Sum(b)))
		return nil
	}

	if err = readChunk(4096); err != nil {
		return
	}
	if err = readChunk(size / 3 * 2); err != nil {
		return
	}
	if err = readChunk(size / 3); err != nil {
		return
	}
	if err = readChunk(size - 8192); err != nil {
		return
	}
	str = strings.Join(ss, ";")
	return
}
