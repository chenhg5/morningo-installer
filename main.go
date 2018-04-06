package main

import (
	"net/http"
	"os"
	"io"
	//"log"
	"archive/zip"
	"path/filepath"
	"io/ioutil"
	"fmt"
)

// moringo 项目安装器
// 下载，安装

func main() {
	url := "https://api.github.com/repos/chenhg5/moringo/zipball/master"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/vnd.github.v3+json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()

	file, _ := os.Create("tmp.zip")
	io.Copy(file, res.Body)

	unzipDir("tmp.zip", "moringo")

	files, _ := ioutil.ReadDir("./moringo")

	os.Rename("./moringo/" + files[0].Name(), "./" + files[0].Name())
	os.Remove("moringo")
	os.Remove("tmp.zip")
	os.Rename("./" + files[0].Name(), "./moringo")

	fmt.Println("install ok!")
}

func unzipDir(zipFile, dir string) {

	r, err := zip.OpenReader(zipFile)
	if err != nil {
		//log.Fatalf("Open zip file failed: %s\n", err.Error())
	}
	defer r.Close()

	for _, f := range r.File {
		func() {
			path := dir + string(filepath.Separator) + f.Name
			os.MkdirAll(filepath.Dir(path), 0755)
			fDest, err := os.Create(path)
			if err != nil {
				//log.Printf("Create failed: %s\n", err.Error())
				return
			}
			defer fDest.Close()

			fSrc, err := f.Open()
			if err != nil {
				//log.Printf("Open failed: %s\n", err.Error())
				return
			}
			defer fSrc.Close()

			_, err = io.Copy(fDest, fSrc)
			if err != nil {
				//log.Printf("Copy failed: %s\n", err.Error())
				return
			}
		}()
	}
}