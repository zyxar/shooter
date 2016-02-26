package shooter

import (
	"os"
	"testing"
)

func TestRequest(t *testing.T) {
	filehash := "66781fc73341bf357500505ad7de1ede;e75e7b0e54b37e3ca511523314e3f1e2;454a5bcb53654a08b2345a606e0cafe6;3f2af6eab10caa9909c0a54652400dd9"
	filename := "Eva.2011.720p.BluRay.x264-DON.mkv"
	files, err := Query(filehash, filename)
	if err != nil {
		t.Error(err)
	}
	chs := make(chan error, len(files))
	for i := range files {
		go func(i int) {
			var err error
			var name string
			name, err = files[i].Fetch(".")
			os.Remove(name)
			chs <- err
		}(i)
	}
	for i := 0; i < len(files); i++ {
		err = <-chs
		if err != nil {
			t.Error(err)
		}
	}
}
