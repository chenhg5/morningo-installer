package main

import (
	"net/http"
	"os"
	"io"
	"archive/zip"
	"path/filepath"
	"io/ioutil"
	"fmt"
	"flag"
	"strings"
	"time"
)

// morningo 项目安装器
// 下载，安装

func main() {

	go func() {
		fmt.Printf("%s", "[")
		consoleStr := "█"
		for i := 0; i != 10; i = i + 1 {
			//log.Println(consoleStr)
			fmt.Printf("%s", consoleStr)
			time.Sleep(time.Second * 1)
		}
	}()

	url := "https://api.github.com/repos/chenhg5/morningo/zipball/master"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/vnd.github.v3+json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()

	file, _ := os.Create("tmp.zip")
	io.Copy(file, res.Body)

	unzipDir("tmp.zip", "tmp")

	files, _ := ioutil.ReadDir("./tmp")

	var name string
	flag.StringVar(&name, "project-name", "morningo", "project name")
	flag.Parse()

	os.Rename("./tmp/" + files[0].Name(), "./" + files[0].Name())
	os.Remove("tmp")
	os.Remove("tmp.zip")
	os.Rename("./" + files[0].Name(), "./" + name)

	os.Remove("./" + files[0].Name() + "/README.md")
	os.Remove("./" + files[0].Name() + "/Dockerfile")
	os.Remove("./" + files[0].Name() + "/LICENSE")
	os.Remove("./" + files[0].Name() + "/pid")
	os.Remove("./" + files[0].Name() + "/.dockerignore")
	os.Remove("./" + files[0].Name() + "/.travis.yml")

	renameProject( "./" + name, name)

	fmt.Printf("%s", "] 100% \ninstall ok!\n\n")
	fmt.Printf("%s", "|  \\/  | ___  _ __ _ __ (_)_ __   __ _ / ___| ___\n")
	fmt.Printf("%s", "| |\\/| |/ _ \\| '__| '_ \\| | '_ \\ / _` | |  _ / _ \\\n")
	fmt.Printf("%s", "| |  | | (_) | |  | | | | | | | | (_| | |_| | (_) |\n")
	fmt.Printf("%s", "|_|  |_|\\___/|_|  |_| |_|_|_| |_|\\__, |\\____|\\___/\n")
	fmt.Printf("%s", "                                 |___/")
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