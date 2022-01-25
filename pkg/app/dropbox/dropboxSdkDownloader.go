package dropbox

import (
	"fmt"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
)

type SDKDownloader struct {
	dbx files.Client
}

func NewSDKDownloader(accessToken string) *SDKDownloader {
	//client := &http.Client{
	//    Transport: &http.Transport{
	//        TLSClientConfig: &tls.Config{
	//            InsecureSkipVerify: true,
	//        },
	//    },
	//}
	conf := dropbox.Config{
		Token:    accessToken,
		LogLevel: dropbox.LogInfo, // if needed, set the desired logging level. Default is off
	}

	s := new(SDKDownloader)
	s.dbx = files.New(conf)
	return s
}

//
//func (s *SDKDownloader) DownloadZip(path string) error {
//	_, content, err := s.dbx.DownloadZip(files.NewDownloadZipArg(path))
//	if err != nil {
//		fmt.Println(err)
//		return err
//	}
//
//	f, err := os.OpenFile("./test.zip", os.O_WRONLY|os.O_CREATE, 0666)
//	if err != nil {
//		fmt.Println(err)
//		return err
//
//	}
//	defer f.Close()
//	written, err := io.Copy(f, content)
//	fmt.Println(written)
//	//bytesData, err := ioutil.ReadAll(content)
//	if err != nil {
//		return err
//	}
//	return nil
//
//	//return &bytesData, nil
//}

func (s *SDKDownloader) GetListFolder(path string, recursive bool, isNeedFile bool) ([]string, error) {
	folder, err := s.dbx.ListFolder(&files.ListFolderArg{
		Path:      path,
		Recursive: recursive,
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var fileList []string

	fileList = s.fillResult(fileList, folder, isNeedFile)

	hasMore := folder.HasMore
	cursor := folder.Cursor
	for hasMore {
		folder, err = s.dbx.ListFolderContinue(&files.ListFolderContinueArg{
			Cursor: cursor,
		})
		if err != nil {
			log.Println(err)
			return fileList, err
		}
		hasMore = folder.HasMore
		cursor = folder.Cursor

		fileList = s.fillResult(fileList, folder, isNeedFile)
	}

	return fileList, nil
}
func (s *SDKDownloader) fillResult(fileList []string, folder *files.ListFolderResult, isFile bool) []string {
	for _, entry := range folder.Entries {
		if isFile {
			switch entry.(type) {
			case *files.FileMetadata:
				fileEntry, _ := entry.(*files.FileMetadata)
				filePath := fileEntry.PathLower
				fmt.Println("Document found :" + filePath)
				fileList = append(fileList, filePath)
			}
		} else {
			switch entry.(type) {
			case *files.FolderMetadata:
				fileEntry, _ := entry.(*files.FolderMetadata)
				filePath := fileEntry.PathLower
				fmt.Println("File found :" + filePath)
				fileList = append(fileList, filePath)
			}
		}
	}

	return fileList
}
func (s *SDKDownloader) DownloadFile(path string) (*files.FileMetadata, *[]byte, error) {
	fileMetadata, content, err := s.dbx.Download(&files.DownloadArg{Path: path})
	if err != nil {
		return nil, nil, err
	}
	defer func(content io.ReadCloser) {
		err := content.Close()
		if err != nil {
			return
		}
	}(content)
	data, err := ioutil.ReadAll(content)
	if err != nil {
		return nil, nil, err
	}

	return fileMetadata, &data, nil
}
