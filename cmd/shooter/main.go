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
		filename := filepath.Base(fullpath)
		files, err := shooter.Query(filehash, filename)
		if err != nil {
			fmt.Println("[ERROR]", err)
			continue
		}
		chs := make(chan error, len(files))
		for i := range files {
			go func(i int) {
				chs <- files[i].Fetch()
			}(i)
		}
		for i := 0; i < len(files); i++ {
			err = <-chs
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
