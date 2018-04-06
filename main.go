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
	"flag"
	"strings"
)

// morningo 项目安装器
// 下载，安装

func main() {
	url := "https://api.github.com/repos/chenhg5/morningo/zipball/master"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/vnd.github.v3+json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()

	file, _ := os.Create("tmp.zip")
	io.Copy(file, res.Body)

	unzipDir("tmp.zip", "morningo")

	files, _ := ioutil.ReadDir("./morningo")

	var name string
	flag.StringVar(&name, "project-name", "morningo", "project name")
	flag.Parse()

	os.Rename("./morningo/" + files[0].Name(), "./" + files[0].Name())
	os.Remove("morningo")
	os.Remove("tmp.zip")
	os.Rename("./" + files[0].Name(), "./" + name)

	renameProject( "./" + name, name)

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

func renameProject(fileDir string, projectName string)  {
	//fmt.Println("path: " +  fileDir)
	files, _ := ioutil.ReadDir(fileDir)
	for _,file := range files{
		if file.IsDir(){
			renameProject(fileDir + "/" + file.Name(), projectName)
		} else {
			path := fileDir + "/" + file.Name()
			//fmt.Println("replace path: " +  path)
			buf, _ := ioutil.ReadFile(path)
			content := string(buf)

			//替换
			newContent := strings.Replace(content, "morningo/", projectName + "/", -1)

			//重新写入
			ioutil.WriteFile(path, []byte(newContent), 0)
		}
	}
}