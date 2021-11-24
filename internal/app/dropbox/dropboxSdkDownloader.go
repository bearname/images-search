package dropbox

import (
	"fmt"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
	"io"
	"os"
)

type SDKDownloader struct {
	dbx files.Client
}

func NewSDKDownloader(accessToken string) *SDKDownloader {
	conf := dropbox.Config{
		Token:    accessToken,
		LogLevel: dropbox.LogInfo, // if needed, set the desired logging level. Default is off
	}

	s := new(SDKDownloader)
	s.dbx = files.New(conf)
	return s
}

func (s *SDKDownloader) Download(path string) error {
	_, content, err := s.dbx.DownloadZip(files.NewDownloadZipArg(path))
	if err != nil {
		fmt.Println(err)
		return err
	}

	f, err := os.OpenFile("./test.zip", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return err

	}
	defer f.Close()
	written, err := io.Copy(f, content)
	fmt.Println(written)
	//bytesData, err := ioutil.ReadAll(content)
	if err != nil {
		return err
	}
	return nil

	//return &bytesData, nil
}
