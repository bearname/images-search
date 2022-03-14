package domain

import (
	"bytes"
	"github.com/col3name/images-search/pkg/domain/pictures"
)

type ImageCompressor interface {
	Compress(fileOrigin *[]byte, quality int, baseWidth int, format pictures.SupportedImgType) (*bytes.Buffer, bool)
}
