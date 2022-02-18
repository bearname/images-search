package main

import (
	"archive/zip"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"photofinish/util"
	"strconv"
	"time"
)

func main() {
	start := time.Now()
	var filePath string
	flag.StringVar(&filePath, "p", "/1", "url to dropbox zip file")
	flag.Parse()
	//fmt.Println("fileUrl")
	//fmt.Println(filePath)
	//
	//downloader := dropbox.NewSDKDownloader("eJWCQYbm6HAAAAAAAAAAAQ_g-B7jQenyuEPZHqzMNwwY1fhD7ozQpNrWtulOZTgq")
	//err := downloader.Download(filePath)
	////err := DownloadFile(outputFile, fileUrl)
	//if err != nil {
	//    fmt.Println(err)
	//    return
	//}
	//fmt.Println("Downloaded: " + filePath)

	util.ClearFolder(filePath)

	zipFile := "./drop-originCopy.zip"
	body, err := ioutil.ReadFile(zipFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		fmt.Println(err)
		return
	}

	for i, file := range zipReader.File {
		fmt.Println("count " + strconv.Itoa(i))
		open, err := file.Open()
		if err != nil {
			fmt.Println(file.Name)
			fmt.Println(err)
		} else {
			all, err := ioutil.ReadAll(open)
			if err != nil {
				fmt.Println(err)
				continue
			}

			err = ioutil.WriteFile("out"+file.Name, all, 0644)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("Success write file " + file.Name)
			}
		}
	}

	fmt.Println(time.Since(start))
}
