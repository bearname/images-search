package domain

import (
	"bytes"
	"photofinish/pkg/domain/pictures"
)

type ImageCompressor interface {
	Compress(fileOrigin *[]byte, quality int, baseWidth int, format pictures.SupportedImgType) (*bytes.Buffer, bool)
}
