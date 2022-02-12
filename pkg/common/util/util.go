package util

import (
	"encoding/base64"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"net"
	"os"
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
	flag.StringVar(&isNeedLoadEnvFile, "d", "false", "is need load .env file")
	flag.Parse()
	if isNeedLoadEnvFile == "true" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
}
func LogToFileIfNeeded() {
	log.SetFormatter(&log.JSONFormatter{})
	var isNeedFileLog string
	flag.StringVar(&isNeedFileLog, "log", "true", "is need load .env file")
	flag.Parse()
	if isNeedFileLog == "false" {
		file, err := os.OpenFile("short.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 6444)
		if err == nil {
			log.SetOutput(file)
			defer func(file *os.File) {
				err = file.Close()
				if err != nil {
					log.Error(err)
				}
			}(file)
		}
	}
}

func ParseEnvString(key string, err error) (string, error) {
	if err != nil {
		return "", err
	}
	str, ok := os.LookupEnv(key)
	if !ok {
		return "", fmt.Errorf("undefined environment variable %v", key)
	}
	return str, nil
}

func IsUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}

func GetRemoteIp(remoteAddr string) (string, error) {
	ip, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		return "", err
	}
	return ip, err
}
