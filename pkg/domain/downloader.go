package domain

type Downloader interface {
	Download(url string) error
}
