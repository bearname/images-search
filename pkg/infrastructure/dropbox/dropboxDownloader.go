package dropbox

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
)

type Downloader struct {
}

func (s *Downloader) Download(dropboxShareLinkToFile string) error {
	fileUrl := "https://dl.dropboxusercontent.com/" + dropboxShareLinkToFile[len("https:://www.dropbox.com"):]
	resp, err := http.Get(fileUrl)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	get := resp.Header.Get("content-type")

	if len(get) == 0 || (len(get) != 0 && get != "application/zip") {
		return errors.New("not supported format " + get)
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	f, err := os.OpenFile("./test.zip", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Println(err)
		return err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)
	_, err = io.Copy(f, resp.Body)
	return err
}
