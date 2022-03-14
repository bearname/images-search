package picture

import (
	"bytes"
	"github.com/col3name/images-search/pkg/domain/pictures"
	"github.com/nfnt/resize"
	log "github.com/sirupsen/logrus"
	"image"
	"image/jpeg"
	"image/png"
)

type ImageCompressor struct {
}

func NewImageCompressor() *ImageCompressor {
	return new(ImageCompressor)
}

func (c *ImageCompressor) Compress(fileOrigin *[]byte, quality int, baseWidth int, format pictures.SupportedImgType) (*bytes.Buffer, bool) {
	var fileOut bytes.Buffer

	origin, typeImage, config, ok := c.decodeImage(format, fileOrigin)

	if !ok {
		return &fileOut, false
	}
	width := uint(baseWidth)
	height := uint(baseWidth * config.Height / config.Width)
	canvas := resize.Thumbnail(width, height, *origin, resize.Lanczos3)

	if typeImage == 0 {
		err := png.Encode(&fileOut, canvas)
		if err != nil {
			log.Println("Failed to compress image")
			return &fileOut, false

		}
	} else {
		err := jpeg.Encode(&fileOut, canvas, &jpeg.Options{Quality: quality})
		if err != nil {
			log.Println("Failed to compress image")
			return &fileOut, false
		}
	}

	return &fileOut, true
}

func (c *ImageCompressor) decodeImage(format pictures.SupportedImgType, fileData *[]byte) (*image.Image, int64, *image.Config, bool) {
	var origin image.Image
	var config image.Config
	var typeImage int64
	var err error
	buffer := bytes.NewBuffer(*fileData)
	if format == pictures.JPEG || format == pictures.JPG {
		typeImage = 1
		origin, err = jpeg.Decode(buffer)
		if err != nil {
			log.Println("jpeg.Decode(fileData)")
			return nil, 0, nil, false
		}
		tmp := bytes.NewBuffer(*fileData)
		config, err = jpeg.DecodeConfig(tmp)
		if err != nil {
			log.Println("jpeg.DecodeConfig(temp)")
			return nil, 0, nil, false
		}
	} else if format == pictures.PNG {
		typeImage = 0
		origin, err = png.Decode(buffer)
		if err != nil {
			log.Println("png.Decode(fileData)")
			return nil, 0, nil, false

		}
		tmp := bytes.NewBuffer(*fileData)
		config, err = png.DecodeConfig(tmp)
		if err != nil {
			log.Println("png.DecodeConfig(temp)")
			return nil, 0, nil, false
		}
	}

	return &origin, typeImage, &config, true
}
