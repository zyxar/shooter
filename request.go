package shooter

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

const (
	userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.111 Safari/537.36"
)

type subtitleDescription struct {
	Desc  string
	Delay int
	Files []SubtitleFile
}

type SubtitleFile struct {
	Ext      string
	Link     string
	FilmName *string `json:"-"`
}

func (this SubtitleFile) String() string {
	v, _ := json.MarshalIndent(this, "", "  ")
	return string(v)
}

func (this *SubtitleFile) Fetch(dirname string) (filename string, err error) {
	if this.FilmName != nil {
		filename = *this.FilmName + "." + this.Ext
	}
	req, err := http.NewRequest("GET", strings.Replace(this.Link, "https://", "http://", 1), nil)
	if err != nil {
		return
	}
	req.Header.Add("Origin", "http://shooter.cn/")
	req.Header.Add("User-Agent", userAgent)
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Referer", "http://shooter.cn/")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 > 3 {
		err = errors.New(resp.Status)
		return
	}

	splits := strings.Split(resp.Header.Get("Content-Disposition"), "filename=")
	if len(splits) > 1 {
		if filepath.Ext(splits[1])[1:] != this.Ext {
			err = errors.New("filename extension not matched.")
			return
		}
		filename = splits[1]
		if this.FilmName == nil {
			filmNameLen := len(splits[1]) - len(this.Ext)
			filmName := splits[1][:filmNameLen]
			this.FilmName = &filmName
		}
	} else if this.FilmName == nil {
		err = errors.New("filename not determined.")
		return
	}

	var saveFile = func(filename string) error {
		if len(dirname) > 0 {
			filename = filepath.Join(dirname, filename)
		}
		if _, err := os.Lstat(filename); os.IsNotExist(err) {
			if file, err := os.Create(filename); err == nil {
				if _, err = io.Copy(file, resp.Body); err != nil {
					return err
				}
				return nil
			} else {
				return err
			}
		}
		return os.ErrExist
	}

	err = saveFile(filename)
	i := 1
	for err == os.ErrExist {
		filename = fmt.Sprintf("%s-%d.%s", *this.FilmName, i, this.Ext)
		err = saveFile(filename)
		i++
	}
	return
}

// [
//   {
//     "Desc": "",
//     "Delay": 0,
//     "Files": [
//       {
//         "Ext": "ass",
//         "Link": "https://www.shooter.cn/api/subapi.php?fetch=MTQ1Mjk1MTI2Nnw5MFNuZ19EdGQ1Um5aQmtWa3JWSkFnbmVKTFNud1pmWHROMFhoVG1rc1loM3hRV1g5Mmg3bTJfbi1GX1hvNkdnZVBsVjQzbEFUMjZZNWNmeGVLUlozMUU3S1ZLSXduU1ZjdG4tN3J6dVlUWGkxT0ZnUVdleU1JcmdCams4R294SGZhbVoxaF9lTFkwckxLUT18tVUfkRy6VKi6-p3We7gbNo29gpZ2rJsdMLW6XLqRjhQ=&nonce=v%CE%1F%3A%40%C5%ED%F0%CF%AA%AC%3B%AF%15p%9D"
//       }
//     ]
//   },
//   {
//     "Desc": "",
//     "Delay": 0,
//     "Files": [
//       {
//         "Ext": "srt",
//         "Link": "https://www.shooter.cn/api/subapi.php?fetch=MTQ1Mjk1MTI2NnxnWkdHc0JDb3FpdnJXdkZ4STJsS2tPeVdfNFJVRkg0MkNQd1dKUjA5a2RyOV9oZFEteGp6RmVlZGpZT3A5c195U2VCNmw1LVczOUJaQ3BBY1NNbjk1Q3QtUGhoWlN0MmUtSkcwU1RTS3lpUDFsRDRlcW5kdWR0Tzc0T1Z1dlg5TTJqVi1UX0FyamdLVm1wcz18yAg8N3tsHJ9HJU8U7UElVBA7EBKnl20miy35fo7S1iY=&nonce=%C9%B3%00%0AzD%86%3C%EB%06D%F0%C10%03Y"
//       }
//     ]
//   }
// ]
// or
// 255
func Query(filehash, filename string) ([]SubtitleFile, error) {
	v := url.Values{}
	v.Set("filehash", filehash)
	v.Set("pathinfo", filename)
	v.Set("format", "json")
	req, err := http.NewRequest("POST", "http://shooter.cn/api/subapi.php", strings.NewReader(v.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Origin", "http://shooter.cn/")
	req.Header.Add("User-Agent", userAgent)
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Referer", "http://shooter.cn/")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 > 3 {
		return nil, errors.New(resp.Status)
	}
	rd := bufio.NewReader(resp.Body)
	if b, err := rd.Peek(1); err != nil {
		return nil, err
	} else if b[0] == 255 {
		return nil, errors.New("subtitles not found.")
	}
	var desc []subtitleDescription
	decoder := json.NewDecoder(rd)
	if err = decoder.Decode(&desc); err != nil {
		return nil, err
	}
	subfiles := make([]SubtitleFile, 0, len(desc))
	filmNameLen := len(filename) - len(filepath.Ext(filename))
	filmName := filename[:filmNameLen]
	for i, _ := range desc {
		for j, _ := range desc[i].Files {
			desc[i].Files[j].FilmName = &filmName
			subfiles = append(subfiles, desc[i].Files[j])
		}
	}
	return subfiles, nil
}
