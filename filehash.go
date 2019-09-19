package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func checkErrf(e error) {
	if e != nil {
		log.Fatalln(e)
	}
}

func getArgs() (arg string, err error) {
	if len(os.Args) == 1 {
		err = errors.New("cannot find command argument")
	} else {
		arg, err = filepath.Abs(os.Args[1])
	}
	return
}

func sums(f *os.File) {
	md5 := md5.New()
	sha1 := sha1.New()
	sha256 := sha256.New()
	w := io.MultiWriter(md5, sha1, sha256)
	if _, err := io.Copy(w, f); err != nil {
		log.Println(err)
		return
	}
	fmt.Printf("File: %s\n", f.Name())
	fmt.Printf("MD5: %X\n", md5.Sum(nil))
	fmt.Printf("SHA-1: %X\n", sha1.Sum(nil))
	fmt.Printf("SHA-256: %X\n\n", sha256.Sum(nil))
}

func main() {
	path, err := getArgs()
	checkErrf(err)

	f, err := os.Open(path)
	checkErrf(err)
	defer f.Close()

	if s, _ := f.Stat(); !s.IsDir() {
		sums(f)
		return

	}
	err = os.Chdir(path)
	checkErrf(err)

	names, err := f.Readdir(0)
	checkErrf(err)

	for _, s := range names {
		if s.IsDir() {
			continue
		}

		fpath, err := filepath.Abs(s.Name())
		if err != nil {
			log.Println(err)
			continue
		}

		file, err := os.Open(fpath)
		if err != nil {
			log.Println(err)
			continue
		}
		sums(file)
		file.Close()
	}
}
