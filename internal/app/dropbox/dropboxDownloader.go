package dropbox

import (
    "fmt"
    "github.com/pkg/errors"
    "io"
    "net/http"
    "os"
)

type Downloader struct {
}

func NewDropboxDownloader() *Downloader {
    return new(Downloader)
}

func (s *Downloader) Download(dropboxShareLinkToFile string) error {
    fileUrl := "https://dl.dropboxusercontent.com/" + dropboxShareLinkToFile[len("https:://www.dropbox.com"):]
    resp, err := http.Get(fileUrl)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    // Write the body to file
    //counter := &WriteCounter{}
    //_, err = io.Copy(out, io.TeeReader(resp.Body, counter))

    fmt.Println(resp.Header)
    get := resp.Header.Get("content-type")

    if len(get) == 0 || (len(get) != 0 && get != "application/zip") {
        return errors.New("not supported format " + get)
    }

    defer resp.Body.Close()

    f, err := os.OpenFile("./test.zip", os.O_WRONLY|os.O_CREATE, 0666)
    if err != nil {
        fmt.Println(err)
        return err
    }
    defer f.Close()
    written, err := io.Copy(f, resp.Body)
    fmt.Println(written)
    //bytesData, err := ioutil.ReadAll(content)
    //if err != nil {
    //    return nil, err
    //}
    return err

    //
    //all, err := ioutil.ReadAll(resp.Body)
    //if err != nil {
    //    return nil, err
    //}
    //return &all, err
}
