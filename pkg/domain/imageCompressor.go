package domain

import "bytes"

type ImageCompressor interface {
	Compress(fileOrigin *[]byte, quality int, baseWidth int, format string) (*bytes.Buffer, bool)
}
