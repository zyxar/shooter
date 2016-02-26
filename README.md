# Shooter 射手字幕下载

[![Go Report Card](https://goreportcard.com/badge/github.com/zyxar/shooter)](https://goreportcard.com/report/github.com/zyxar/shooter)
[![GoDoc](https://godoc.org/github.com/zyxar/shooter?status.svg)](https://godoc.org/github.com/zyxar/shooter)

cli and lib for downloading subtitles from SHOOTER

## Install

- cli: `go get github.com/zyxar/shooter/cmd/shooter`
- lib: `go get github.com/zyxar/shooter`

## Usage of `cli:shooter`

```shell
$ shooter -dir {TARGET_DIRECTORY} {FILM_FILE}...
```

- Default {TARGET_DIRECTORY} is present working directory.
- `http_proxy` environment variable is respected.

## Alternative

- [npmjs.com/package/shooter](https://github.com/zyxar/shooter.js)

However, the Node.js implementation has some potential issue (contention) when saving subtitle files.
