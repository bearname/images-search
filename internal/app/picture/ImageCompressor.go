package picture

import (
	"bytes"
	"fmt"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"image/png"
	"strings"
)

type ImageCompressor struct {
}

func NewImageCompressor() *ImageCompressor {
	return new(ImageCompressor)
}

func (c *ImageCompressor) Compress(fileOrigin []byte, quality int, baseWidth int, format string) (bytes.Buffer, bool) {
	var fileOut bytes.Buffer

	format = strings.ToLower(format)
	origin, typeImage, config, ok := c.decodeImage(format, fileOrigin)
	if !ok {
		return fileOut, false
	}
	width := uint(baseWidth)
	height := uint(baseWidth * config.Height / config.Width)
	canvas := resize.Thumbnail(width, height, origin, resize.Lanczos3)

	if typeImage == 0 {
		err := png.Encode(&fileOut, canvas)
		if err != nil {
			fmt.Println("Failed to compress image")
			return fileOut, false

		}
	} else {
		err := jpeg.Encode(&fileOut, canvas, &jpeg.Options{Quality: quality})
		if err != nil {
			fmt.Println("Failed to compress image")
			return fileOut, false

		}
	}

	return fileOut, true
}

func (c *ImageCompressor) decodeImage(format string, fileData []byte) (image.Image, int64, image.Config, bool) {
	var origin image.Image
	var config image.Config
	var typeImage int64
	var err error
	buffer := bytes.NewBuffer(fileData)
	if format == "jpg" || format == "jpeg" {
		typeImage = 1
		origin, err = jpeg.Decode(buffer)
		if err != nil {
			fmt.Println("jpeg.Decode(fileData)")
			return nil, 0, image.Config{}, false
		}
		tmp := bytes.NewBuffer(fileData)
		config, err = jpeg.DecodeConfig(tmp)
		if err != nil {
			fmt.Println("jpeg.DecodeConfig(temp)")
			return nil, 0, image.Config{}, false
		}
	} else if format == "png" {
		typeImage = 0
		origin, err = png.Decode(buffer)
		if err != nil {
			fmt.Println("png.Decode(fileData)")
			return nil, 0, image.Config{}, false

		}
		tmp := bytes.NewBuffer(fileData)
		config, err = jpeg.DecodeConfig(tmp)

		config, err = png.DecodeConfig(tmp)
		if err != nil {
			fmt.Println("png.DecodeConfig(temp)")
			return nil, 0, image.Config{}, false
		}
	}

	return origin, typeImage, config, true
}
