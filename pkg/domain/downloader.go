package domain

import (
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
)

type Downloader interface {
	GetListFolder(path string, recursive bool, isNeedFile bool) ([]string, error)
	DownloadFile(path string) (*files.FileMetadata, *[]byte, error)
	// DownloadZip(url string) error
}
