package main

import (
	"flag"
	"fmt"
	"path/filepath"

	"github.com/zyxar/shooter"
)

func main() {
	flag.Parse()
	for _, fullpath := range flag.Args() {
		filehash, err := shooter.FileHash(fullpath)
		if err != nil {
			fmt.Println("[ERROR]", err)
			continue
		}
		fullpath, err = filepath.Abs(fullpath)
		if err != nil {
			fmt.Println("[ERROR]", err)
			continue
		}
		dirname, filename := filepath.Split(fullpath)
		files, err := shooter.Query(filehash, filename)
		if err != nil {
			fmt.Println("[ERROR]", err)
			continue
		}
		filesNum := len(files)
		fmt.Printf("Found %d subtitles for %s\n", filesNum, filename)
		chs := make(chan error, filesNum)
		for i := range files {
			go func(i int) {
				fn, err := files[i].Fetch(dirname)
				if err != nil {
					fmt.Printf("[ERROR] %s %v\n", fn, err)
				} else {
					fmt.Printf("[DONE] %s\n", fn)
				}
				chs <- err
			}(i)
		}
		for i := 0; i < filesNum; i++ {
			err = <-chs
		}
	}
}
