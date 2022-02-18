package main

import (
    "flag"
    "fmt"
    "github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
    "github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
    "github.com/pkg/errors"
    log "github.com/sirupsen/logrus"
    "io/ioutil"
    "net/http"
    _ "net/http/pprof"
    "photofinish/pkg/app/picture"
    "photofinish/pkg/domain/pictures"
    "photofinish/util"
    "runtime"
    "sync"
    "time"
)

var compressor picture.ImageCompressor

func main() {
    go func() {
        err := http.ListenAndServe("0.0.0.0:8081", nil)
        if err != nil {
            log.Fatal(err)
        }
    }()
    compressor = *picture.NewImageCompressor()
    start := time.Now()
    var path string
    flag.StringVar(&path, "p", "/1", "dropbox folder path")
    flag.Parse()
    accessToken := "eJWCQYbm6HAAAAAAAAAAAQ_g-B7jQenyuEPZHqzMNwwY1fhD7ozQpNrWtulOZTgq"
    handleFile(accessToken, path)
    log.Println(time.Since(start))
}

func handleFile(accessToken string, path string) {
    conf := dropbox.Config{
        Token:    accessToken,
        LogLevel: dropbox.LogInfo,
    }
    dbx := files.New(conf)
    var fileList []string
    hasMore := true
    for hasMore {
        folder, err := dbx.ListFolder(&files.ListFolderArg{
            Path:      path,
            Recursive: true,
        })
        if err != nil {
            log.Println(err)
            return
        }

        for _, entry := range folder.Entries {
            switch entry.(type) {
            case *files.FileMetadata:
                fileEntry, _ := entry.(*files.FileMetadata)
                filePath := fileEntry.PathLower
                log.Println("Document found :" + filePath)
                fileList = append(fileList, filePath)
            }
        }
        hasMore = folder.HasMore
    }

    rootFolder := path[1:]
    util.ClearFolder(path)

    ch := make(chan string)
    wg := sync.WaitGroup{}
    //for t := 0; t < runtime.NumCPU(); t++ {
    //    wg.Add(1)
    //    go downloadFileAsync(ch, &wg, dbx, rootFolder)
    //}
    consumeFile(ch, &wg, dbx, rootFolder)

    senderFile(fileList, ch, &wg)
}

func senderFile(fileList []string, ch chan string, wg *sync.WaitGroup) {
    for _, line := range fileList {
        ch <- line
    }
    close(ch)

    wg.Wait()
}

func consumeFile(ch chan string, wg *sync.WaitGroup, dbx files.Client, rootFolder string) {
    for t := 0; t < runtime.NumCPU(); t++ {
        wg.Add(1)
        go downloadFileAsync(ch, wg, dbx, rootFolder)
    }
}

func downloadFileAsync(ch chan string, wg *sync.WaitGroup, dbx files.Client, rootFolder string) {
    for filePath := range ch {
        err := downloadFile(dbx, rootFolder, filePath)
        if err != nil {
            log.Println(err)
        } else {
            log.Println("File '" + filePath + "' download success")
        }
    }
    wg.Done()
}

func downloadFile(dbx files.Client, rootFolder string, filePath string) error {
    download, content, err := dbx.Download(&files.DownloadArg{Path: filePath})
    if err != nil {
        return err
    }
    data, err := ioutil.ReadAll(content)
    if err != nil {
        return err
    }
    compressed, ok := compressor.Compress(&data, 80, 300, pictures.JPEG)
    if !ok {
        err = errors.New("Failed scale image")
        data = nil
        //c.handleError(image, err)
        return err
    }

    all, err := ioutil.ReadAll(compressed)
    if err != nil {
        return err
    }
    ioutil.WriteFile(rootFolder+"/compr-"+download.Name, all, 0644)
    return ioutil.WriteFile(rootFolder+"/"+download.Name, data, 0644)
}
