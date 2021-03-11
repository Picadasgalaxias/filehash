package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
)

func checkErrfatal(e error) {
	if e != nil {
		log.Fatalln(e)
	}
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
	fmt.Printf("File: %s\nMD5: %x\nSHA-1: %x\nSHA-256: %x\n\n",
		f.Name(), md5.Sum(nil), sha1.Sum(nil), sha256.Sum(nil))
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("not enough arguments, accept files or directories")
		return
	}
	for _, path := range os.Args[1:] {
		func() {
			file, err := os.Open(path)
			checkErrfatal(err)
			defer file.Close()

			if s, _ := file.Stat(); !s.IsDir() {
				sums(file)
				return

			}

			err = os.Chdir(path)
			checkErrfatal(err)

			names, err := file.Readdir(0)
			checkErrfatal(err)

			for _, s := range names {
				if s.IsDir() {
					continue
				}

				file, err := os.Open(s.Name())
				if err != nil {
					log.Println(err)
					continue
				}
				sums(file)
				file.Close()
			}
		}()
	}
}
