package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/zyxar/shooter"
)

var dirname string

func init() {
	flag.StringVar(&dirname, "dir", "", "set target directory")
	flag.Usage = func() {
		fmt.Println("Usage: shooter [OPTION] film_file")
		flag.PrintDefaults()
	}
}

type message struct {
	body io.ReadCloser
	fn   string
}

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}
	for _, fullpath := range flag.Args() {
		filehash, err := shooter.FileHash(fullpath)
		if err != nil {
			fmt.Println("[ERROR]", err)
			continue
		}
		var filename string
		if dirname == "" {
			fullpath, err = filepath.Abs(fullpath)
			if err != nil {
				fmt.Println("[ERROR]", err)
				continue
			}
			dirname, filename = filepath.Split(fullpath)
		} else {
			filename = filepath.Base(fullpath)
		}
		files, err := shooter.Query(filehash, filename)
		if err != nil {
			fmt.Println("[ERROR]", err)
			continue
		}
		filesNum := len(files)
		fmt.Printf("Found %d subtitles for %s\n", filesNum, filename)
		msgChan := make(chan *message, filesNum)
		for i := range files {
			go func(i int) {
				body, fn, err := files[i].FetchContent()
				if err != nil {
					fmt.Printf("[ERROR] %s-%d %v\n", fn, i, err)
					msgChan <- nil
				} else {
					msgChan <- &message{body, fn}
				}
			}(i)
		}
		j := 1
		var saveFile = func(body io.ReadCloser, fn string) {
			ext := filepath.Ext(fn)
			if dirname != "" {
				fn = filepath.Join(dirname, fn)
			}
			filename := fn
			fn = fn[:len(fn)-len(ext)]
			var file *os.File
			var err error
		retry:
			if _, err = os.Lstat(filename); os.IsNotExist(err) {
				if file, err = os.Create(filename); err == nil {
					_, err = io.Copy(file, body)
					body.Close()
					file.Close()
					fmt.Printf("[DONE] %s\n", filename)
				} else {
					fmt.Println("[ERROR]", err)
				}
			} else {
				filename = fmt.Sprintf("%s-%d%s", fn, j, ext)
				j++
				goto retry
			}
			return
		}
		for i := 0; i < filesNum; i++ {
			msg := <-msgChan
			if msg != nil {
				saveFile(msg.body, msg.fn)
			}
		}
	}
}
