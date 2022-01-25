package util

import (
	"encoding/base64"
	"flag"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"regexp"
)

func ImageBase64(buf []byte) ([]byte, error) {
	imgBase64Str := base64.StdEncoding.EncodeToString(buf)
	decodedImage, err := base64.StdEncoding.DecodeString(imgBase64Str)
	return decodedImage, err
}

func ExtractNumberFromString(input string) []string {
	re := regexp.MustCompile(`[-]?\d[\d,]*[.]?[\d{2}]*`)

	parts := re.FindAllString(input, -1)
	var result []string
	result = append(result, parts...)

	return result
}

func LoadEnvFileIfNeeded() {
	var isNeedLoadEnvFile string
	flag.StringVar(&isNeedLoadEnvFile, "d", "true", "is need load .env file")
	flag.Parse()
	if isNeedLoadEnvFile == "true" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
}

func IsUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}
