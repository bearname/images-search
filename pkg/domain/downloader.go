package domain

import (
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
)

type Downloader interface {
	GetListFiles(dropboxPath string) ([]string, error)
	GetListFolder(path string) ([]string, error)
	DownloadFile(path string) (*files.FileMetadata, *[]byte, error)
}
