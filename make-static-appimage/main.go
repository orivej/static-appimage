package main

import (
	"archive/zip"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/orivej/e"
)

const runtime = "static-appimage-runtime"

func main() {
	log.SetFlags(0)

	runtimePath, err := exec.LookPath(runtime)
	if err != nil {
		log.Fatalf("failed to find '%s' src PATH: %s", runtime, err)
	}

	if len(os.Args) != 3 {
		log.Fatalf("usage: %s APPDIR DESTINATION", os.Args[0])
	}
	appdir := os.Args[1]
	dstPath := os.Args[2]

	runtimeFile, err := os.Open(runtimePath)
	e.Exit(err)
	dst, err := os.OpenFile(dstPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	e.Exit(err)
	_, err = io.Copy(dst, runtimeFile)
	e.Exit(err)
	e.CloseOrExit(runtimeFile)

	zipOffset, err := dst.Seek(0, io.SeekCurrent)
	e.Exit(err)

	zipWriter := zip.NewWriter(dst)
	zipWriter.SetOffset(zipOffset)

	err = filepath.Walk(appdir, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fi.IsDir() {
			return nil
		}

		relpath, err := filepath.Rel(appdir, path)
		e.Exit(err)

		zipHeader := &zip.FileHeader{
			Name:   relpath,
			Method: zip.Deflate,
		}
		zipHeader.SetMode(fi.Mode())

		w, err := zipWriter.CreateHeader(zipHeader)
		e.Exit(err)

		if fi.Mode()&os.ModeSymlink != 0 {
			target, err := os.Readlink(path)
			e.Exit(err)
			_, err = w.Write([]byte(target))
			e.Exit(err)
		} else {
			src, err := os.Open(path)
			e.Exit(err)
			_, err = io.Copy(w, src)
			e.Exit(err)
			e.CloseOrExit(src)
		}

		return nil
	})
	e.Exit(err)

	e.CloseOrExit(zipWriter)
}