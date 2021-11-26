package picture

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go/service/rekognition"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
	"photofinish/internal/common/uuid"
	"photofinish/internal/domain"
	"photofinish/internal/domain/picture"
	"regexp"
	"strconv"
)

type ServiceImpl struct {
	pictureRepo    picture.Repository
	awsRekognition *rekognition.Rekognition
	//awsBucket      string
	//uploader       *manager.Uploader
	uploader   domain.Uploader
	compressor domain.ImageCompressor
}

func NewPictureService(pictureRepo picture.Repository, awsRekognition *rekognition.Rekognition, uploader domain.Uploader, compressor domain.ImageCompressor) *ServiceImpl {
	s := new(ServiceImpl)
	s.pictureRepo = pictureRepo
	s.awsRekognition = awsRekognition
	s.uploader = uploader
	s.compressor = compressor
	//s.awsBucket = awsBucket
	//s.uploader = uploader
	return s
}

//func (s *ServiceImpl) upload(filename string, file io.Reader, acl types.ObjectCannedACL) (*manager.UploadOutput, error) {
//    index := strings.Index(filename, ".")
//    fileName := filename[:index] + uuid.Generate().String() + filename[index+1:]
//
//    upload, err := s.uploader.Upload(context.TODO(), &s3.PutObjectInput{
//        Bucket: aws.String(s.awsBucket),
//        Key:    aws.String(fileName),
//        Body:   file,
//        ACL:    acl,
//    })
//    return upload, err
//}

func (s *ServiceImpl) Create(imageTextDetectionDto *picture.TextDetectionOnImageDto) error {
	return (s.pictureRepo).Store(imageTextDetectionDto)
}

const MOCK_DATA = `{
    "TextDetections": [
        {
            "DetectedText": "groma",
            "Type": "LINE",
            "Id": 0,
            "Confidence": 71.6114730834961,
            "Geometry": {
                "BoundingBox": {
                    "Width": 0.04115382954478264,
                    "Height": 0.02400638721883297,
                    "Left": 0.3590179979801178,
                    "Top": 0.33090221881866455
                },
                "Polygon": [
                    {
                        "X": 0.35934239625930786,
                        "Y": 0.33090221881866455
                    },
                    {
                        "X": 0.40017181634902954,
                        "Y": 0.33221757411956787
                    },
                    {
                        "X": 0.3998474180698395,
                        "Y": 0.3549085855484009
                    },
                    {
                        "X": 0.3590179979801178,
                        "Y": 0.35359323024749756
                    }
                ]
            }
        },
        {
            "DetectedText": "Nacional",
            "Type": "LINE",
            "Id": 1,
            "Confidence": 99.84120178222656,
            "Geometry": {
                "BoundingBox": {
                    "Width": 0.043958697468042374,
                    "Height": 0.019329529255628586,
                    "Left": 0.354614794254303,
                    "Top": 0.415447860956192
                },
                "Polygon": [
                    {
                        "X": 0.354614794254303,
                        "Y": 0.416672945022583
                    },
                    {
                        "X": 0.3983484208583832,
                        "Y": 0.415447860956192
                    },
                    {
                        "X": 0.39857348799705505,
                        "Y": 0.4335522949695587
                    },
                    {
                        "X": 0.35483986139297485,
                        "Y": 0.4347773790359497
                    }
                ]
            }
        },
        {
            "DetectedText": "5000 e 10000",
            "Type": "LINE",
            "Id": 2,
            "Confidence": 45.583614349365234,
            "Geometry": {
                "BoundingBox": {
                    "Width": 0.04651162773370743,
                    "Height": 0.014639639295637608,
                    "Left": 0.35408851504325867,
                    "Top": 0.43355855345726013
                },
                "Polygon": [
                    {
                        "X": 0.35408851504325867,
                        "Y": 0.43355855345726013
                    },
                    {
                        "X": 0.4006001353263855,
                        "Y": 0.43355855345726013
                    },
                    {
                        "X": 0.4006001353263855,
                        "Y": 0.44819819927215576
                    },
                    {
                        "X": 0.35408851504325867,
                        "Y": 0.44819819927215576
                    }
                ]
            }
        },
        {
            "DetectedText": "e",
            "Type": "LINE",
            "Id": 3,
            "Confidence": 94.34638214111328,
            "Geometry": {
                "BoundingBox": {
                    "Width": 0.007501875516027212,
                    "Height": 0.009009009227156639,
                    "Left": 0.3720930218696594,
                    "Top": 0.4369369447231293
                },
                "Polygon": [
                    {
                        "X": 0.3720930218696594,
                        "Y": 0.4369369447231293
                    },
                    {
                        "X": 0.3795948922634125,
                        "Y": 0.4369369447231293
                    },
                    {
                        "X": 0.3795948922634125,
                        "Y": 0.44594594836235046
                    },
                    {
                        "X": 0.3720930218696594,
                        "Y": 0.44594594836235046
                    }
                ]
            }
        },
        {
            "DetectedText": "10000",
            "Type": "LINE",
            "Id": 4,
            "Confidence": 97.83013153076172,
            "Geometry": {
                "BoundingBox": {
                    "Width": 0.02100525051355362,
                    "Height": 0.010135134682059288,
                    "Left": 0.3795948922634125,
                    "Top": 0.4380630552768707
                },
                "Polygon": [
                    {
                        "X": 0.3795948922634125,
                        "Y": 0.4380630552768707
                    },
                    {
                        "X": 0.4006001353263855,
                        "Y": 0.4380630552768707
                    },
                    {
                        "X": 0.4006001353263855,
                        "Y": 0.44819819927215576
                    },
                    {
                        "X": 0.3795948922634125,
                        "Y": 0.44819819927215576
                    }
                ]
            }
        },
        {
            "DetectedText": "6",
            "Type": "LINE",
            "Id": 5,
            "Confidence": 99.43851470947266,
            "Geometry": {
                "BoundingBox": {
                    "Width": 0.03150787577033043,
                    "Height": 0.05968468636274338,
                    "Left": 0.6849212050437927,
                    "Top": 0.4121621549129486
                },
                "Polygon": [
                    {
                        "X": 0.6849212050437927,
                        "Y": 0.4121621549129486
                    },
                    {
                        "X": 0.7164291143417358,
                        "Y": 0.4121621549129486
                    },
                    {
                        "X": 0.7164291143417358,
                        "Y": 0.4718468487262726
                    },
                    {
                        "X": 0.6849212050437927,
                        "Y": 0.4718468487262726
                    }
                ]
            }
        },
        {
            "DetectedText": "252",
            "Type": "LINE",
            "Id": 6,
            "Confidence": 100,
            "Geometry": {
                "BoundingBox": {
                    "Width": 0.10352588444948196,
                    "Height": 0.07432432472705841,
                    "Left": 0.32033008337020874,
                    "Top": 0.4515765905380249
                },
                "Polygon": [
                    {
                        "X": 0.3210802674293518,
                        "Y": 0.4515765905380249
                    },
                    {
                        "X": 0.4238559603691101,
                        "Y": 0.4515765905380249
                    },
                    {
                        "X": 0.42310577630996704,
                        "Y": 0.5259009003639221
                    },
                    {
                        "X": 0.32033008337020874,
                        "Y": 0.5247747898101807
                    }
                ]
            }
        },
        {
            "DetectedText": "LISBOA 2020",
            "Type": "LINE",
            "Id": 7,
            "Confidence": 99.46501922607422,
            "Geometry": {
                "BoundingBox": {
                    "Width": 0.0649627298116684,
                    "Height": 0.023643307387828827,
                    "Left": 0.3401569724082947,
                    "Top": 0.5431105494499207
                },
                "Polygon": [
                    {
                        "X": 0.3406303822994232,
                        "Y": 0.5431105494499207
                    },
                    {
                        "X": 0.40511971712112427,
                        "Y": 0.5465085506439209
                    },
                    {
                        "X": 0.4046463370323181,
                        "Y": 0.5667538642883301
                    },
                    {
                        "X": 0.3401569724082947,
                        "Y": 0.5633558630943298
                    }
                ]
            }
        },
        {
            "DetectedText": "groma",
            "Type": "WORD",
            "Id": 8,
            "ParentId": 0,
            "Confidence": 71.6114730834961,
            "Geometry": {
                "BoundingBox": {
                    "Width": 0.04052577540278435,
                    "Height": 0.022522522136569023,
                    "Left": 0.3593398332595825,
                    "Top": 0.3310810923576355
                },
                "Polygon": [
                    {
                        "X": 0.3593398332595825,
                        "Y": 0.3310810923576355
                    },
                    {
                        "X": 0.39984995126724243,
                        "Y": 0.33220720291137695
                    },
                    {
                        "X": 0.39984995126724243,
                        "Y": 0.3547297418117523
                    },
                    {
                        "X": 0.3593398332595825,
                        "Y": 0.3536036014556885
                    }
                ]
            }
        },
        {
            "DetectedText": "Nacional",
            "Type": "WORD",
            "Id": 9,
            "ParentId": 1,
            "Confidence": 99.84120178222656,
            "Geometry": {
                "BoundingBox": {
                    "Width": 0.04352544993162155,
                    "Height": 0.018018018454313278,
                    "Left": 0.35483869910240173,
                    "Top": 0.4166666567325592
                },
                "Polygon": [
                    {
                        "X": 0.35483869910240173,
                        "Y": 0.4166666567325592
                    },
                    {
                        "X": 0.3983495831489563,
                        "Y": 0.41554054617881775
                    },
                    {
                        "X": 0.3983495831489563,
                        "Y": 0.43355855345726013
                    },
                    {
                        "X": 0.35483869910240173,
                        "Y": 0.434684693813324
                    }
                ]
            }
        },
        {
            "DetectedText": "6",
            "Type": "WORD",
            "Id": 13,
            "ParentId": 5,
            "Confidence": 99.43851470947266,
            "Geometry": {
                "BoundingBox": {
                    "Width": 0.03150787577033043,
                    "Height": 0.05968468636274338,
                    "Left": 0.6849212050437927,
                    "Top": 0.4121621549129486
                },
                "Polygon": [
                    {
                        "X": 0.6849212050437927,
                        "Y": 0.4121621549129486
                    },
                    {
                        "X": 0.7164291143417358,
                        "Y": 0.4121621549129486
                    },
                    {
                        "X": 0.7164291143417358,
                        "Y": 0.4718468487262726
                    },
                    {
                        "X": 0.6849212050437927,
                        "Y": 0.4718468487262726
                    }
                ]
            }
        },
        {
            "DetectedText": "5000 e 10000",
            "Type": "WORD",
            "Id": 10,
            "ParentId": 2,
            "Confidence": 45.583614349365234,
            "Geometry": {
                "BoundingBox": {
                    "Width": 0.04651162773370743,
                    "Height": 0.014639639295637608,
                    "Left": 0.35408851504325867,
                    "Top": 0.43355855345726013
                },
                "Polygon": [
                    {
                        "X": 0.35408851504325867,
                        "Y": 0.43355855345726013
                    },
                    {
                        "X": 0.4006001353263855,
                        "Y": 0.43355855345726013
                    },
                    {
                        "X": 0.4006001353263855,
                        "Y": 0.44819819927215576
                    },
                    {
                        "X": 0.35408851504325867,
                        "Y": 0.44819819927215576
                    }
                ]
            }
        },
        {
            "DetectedText": "252",
            "Type": "WORD",
            "Id": 14,
            "ParentId": 6,
            "Confidence": 100,
            "Geometry": {
                "BoundingBox": {
                    "Width": 0.1027756929397583,
                    "Height": 0.07432810962200165,
                    "Left": 0.3210802674293518,
                    "Top": 0.4515765905380249
                },
                "Polygon": [
                    {
                        "X": 0.3210802674293518,
                        "Y": 0.4515765905380249
                    },
                    {
                        "X": 0.4238559603691101,
                        "Y": 0.4515765905380249
                    },
                    {
                        "X": 0.42310577630996704,
                        "Y": 0.5259009003639221
                    },
                    {
                        "X": 0.32033008337020874,
                        "Y": 0.5247747898101807
                    }
                ]
            }
        },
        {
            "DetectedText": "e",
            "Type": "WORD",
            "Id": 11,
            "ParentId": 3,
            "Confidence": 94.34638214111328,
            "Geometry": {
                "BoundingBox": {
                    "Width": 0.007501875516027212,
                    "Height": 0.009009009227156639,
                    "Left": 0.3720930218696594,
                    "Top": 0.4369369447231293
                },
                "Polygon": [
                    {
                        "X": 0.3720930218696594,
                        "Y": 0.4369369447231293
                    },
                    {
                        "X": 0.3795948922634125,
                        "Y": 0.4369369447231293
                    },
                    {
                        "X": 0.3795948922634125,
                        "Y": 0.44594594836235046
                    },
                    {
                        "X": 0.3720930218696594,
                        "Y": 0.44594594836235046
                    }
                ]
            }
        },
        {
            "DetectedText": "10000",
            "Type": "WORD",
            "Id": 12,
            "ParentId": 4,
            "Confidence": 97.83013153076172,
            "Geometry": {
                "BoundingBox": {
                    "Width": 0.02100525051355362,
                    "Height": 0.010135134682059288,
                    "Left": 0.3795948922634125,
                    "Top": 0.4380630552768707
                },
                "Polygon": [
                    {
                        "X": 0.3795948922634125,
                        "Y": 0.4380630552768707
                    },
                    {
                        "X": 0.4006001353263855,
                        "Y": 0.4380630552768707
                    },
                    {
                        "X": 0.4006001353263855,
                        "Y": 0.44819819927215576
                    },
                    {
                        "X": 0.3795948922634125,
                        "Y": 0.44819819927215576
                    }
                ]
            }
        },
        {
            "DetectedText": "LISBOA",
            "Type": "WORD",
            "Id": 15,
            "ParentId": 7,
            "Confidence": 99.4224624633789,
            "Geometry": {
                "BoundingBox": {
                    "Width": 0.03675919026136398,
                    "Height": 0.018033629283308983,
                    "Left": 0.3405851423740387,
                    "Top": 0.545045018196106
                },
                "Polygon": [
                    {
                        "X": 0.3405851423740387,
                        "Y": 0.545045018196106
                    },
                    {
                        "X": 0.3773443400859833,
                        "Y": 0.545045018196106
                    },
                    {
                        "X": 0.3765941560268402,
                        "Y": 0.5630630850791931
                    },
                    {
                        "X": 0.3405851423740387,
                        "Y": 0.5619369149208069
                    }
                ]
            }
        },
        {
            "DetectedText": "2020",
            "Type": "WORD",
            "Id": 16,
            "ParentId": 7,
            "Confidence": 99.50758361816406,
            "Geometry": {
                "BoundingBox": {
                    "Width": 0.027756938710808754,
                    "Height": 0.01915883645415306,
                    "Left": 0.3773443400859833,
                    "Top": 0.5472972989082336
                },
                "Polygon": [
                    {
                        "X": 0.3773443400859833,
                        "Y": 0.5472972989082336
                    },
                    {
                        "X": 0.4051012694835663,
                        "Y": 0.5472972989082336
                    },
                    {
                        "X": 0.4043510854244232,
                        "Y": 0.5664414167404175
                    },
                    {
                        "X": 0.3773443400859833,
                        "Y": 0.565315306186676
                    }
                ]
            }
        }
    ],
    "TextModelVersion": "3.0"
}`

func (s *ServiceImpl) DetectImageFromArchive(root string, minConfidence int, eventId int64) error {
	var pictures []*picture.TextDetectionOnImageDto
	err := filepath.Walk(root,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			fmt.Println(path, info.Size())
			if !info.IsDir() {
				textDetection, err := s.textTextOnImageFile(eventId, path, minConfidence)
				if err != nil {
					log.Println(err.Error())
					return err
				}

				pictures = append(pictures, textDetection)
				if len(pictures) > 400 {
					err = s.storeAllImages(pictures)
					if err != nil {
						return err
					}
				}
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}

	//for _, file := range files {
	//    fmt.Printf("=%s\n", file.Name)
	//    if !file.FileInfo().IsDir() {
	//
	//        textDetection, err := s.textTextOnImageFile(eventId, file, minConfidence)
	//        if err != nil {
	//            log.Println(err.Error())
	//            return err
	//        }
	//
	//        pictures = append(pictures, textDetection)
	//        if len(pictures) > 400 {
	//            err = s.storeAllImages(pictures)
	//            if err != nil {
	//                return err
	//            }
	//        }
	//    }
	//}

	if len(pictures) > 0 {
		err := s.storeAllImages(pictures)
		if err != nil {
			return err
		}
	}
	pictures = nil

	return nil
}

func (s *ServiceImpl) storeAllImages(pictures []*picture.TextDetectionOnImageDto) error {
	err := (s.pictureRepo).StoreAll(pictures)
	if err != nil {
		log.Println(err.Error())
		return err
	} else {
		pictures = []*picture.TextDetectionOnImageDto{}
	}
	return nil
}

func (s *ServiceImpl) Search(dto picture.SearchPictureDto) (picture.SearchPictureResultDto, error) {
	return (s.pictureRepo).Search(dto)
}

func (s *ServiceImpl) Delete(pictureId string) error {
	err := s.pictureRepo.FindById(pictureId)
	if err != nil {
		log.Error(err)
		return err
	}
	return (s.pictureRepo).Delete(pictureId)
}

func (s *ServiceImpl) textTextOnImageFile(eventId int64, path string, minConfidence int) (*picture.TextDetectionOnImageDto, error) {
	//fileBytes, err := readAll(path)
	//if err != nil {
	//    return nil, err
	//}
	fileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	filename := uuid.Generate().String()
	originalImageUploadOutput, err := s.uploader.Upload(strconv.FormatInt(eventId, 10)+"/"+filename+".jpg", bytes.NewReader(fileBytes), types.ObjectCannedACLBucketOwnerRead)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	compressBuffer, ok := s.compressor.Compress(fileBytes, 100, 300, "jpg")

	var thumbnailImage *domain.UploadOutput
	if ok {
		//all, err := ioutil.ReadAll(&compressBuffer)
		//if err != nil {
		//    log.Error(err)
		//    return nil, err
		//}

		//filename = "web/images/" + filename + "-thumb.jpg"
		//fmt.Print("filename")
		//fmt.Println(filename)
		//err = writeFile(filename, all)
		//if err != nil {
		//    log.Error(err)
		//}
		thumbnailImage, err = s.uploader.Upload(strconv.FormatInt(eventId, 10)+"/"+filename+"-thumb.jpg", &compressBuffer, types.ObjectCannedACLPublicRead)
		if err != nil {
			log.Error(err)
			return nil, err
		}
	}

	var decodedImage []byte
	decodedImage, err = imageBase64(fileBytes)

	if err != nil {
		log.Println("decodestring")
		log.Println(err.Error())
		return nil, err
	}
	var textDetection []picture.TextDetection

	textDetection, err = s.detectText(decodedImage, minConfidence)
	if err != nil {
		log.Println("detect")
		log.Println(err.Error())
		return nil, err
	}
	decodedImage = nil

	log.Println(textDetection)

	return picture.NewImageTextDetectionDto(eventId, originalImageUploadOutput.Location, thumbnailImage.Location, textDetection), nil
	//return picture.NewImageTextDetectionDto(eventId, "originalImageUploadOutput.Location", " thumbnailImage.Location", textDetection), nil
}

func readAll(file *zip.File) ([]byte, error) {
	fc, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer closeFile(fc)

	content, err := ioutil.ReadAll(fc)
	if err != nil {
		return nil, err
	}

	return content, nil
}

type myCloser interface {
	Close() error
}

func closeFile(f myCloser) {
	err := f.Close()
	if err != nil {
		log.Error(err)
	}
}

func writeFile(filename string, bytes []byte) error {
	err := ioutil.WriteFile(filename, bytes, 0644)
	if err != nil {
		return err
	}

	return err
}

func imageBase64(buf []byte) ([]byte, error) {
	imgBase64Str := base64.StdEncoding.EncodeToString(buf)
	decodedImage, err := base64.StdEncoding.DecodeString(imgBase64Str)
	return decodedImage, err
}

func (s *ServiceImpl) detectText(decodedImage []byte, minConfidence int) ([]picture.TextDetection, error) {
	input := &rekognition.DetectTextInput{
		Image: &rekognition.Image{
			Bytes: decodedImage,
		},
		//Image: &rekognition.Image{
		//    S3Object: &rekognition.S3Object{
		//        Bucket: bucket,
		//        Name:   photo,
		//    },
		//},
	}

	//var result rekognition.DetectTextOutput
	//err := json.Unmarshal([]byte(MOCK_DATA), &result)
	result, err := s.awsRekognition.DetectText(input)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	log.Println(result.GoString())

	set := domain.MakeSet()
	for _, detection := range result.TextDetections {
		//log.Println("===")
		//log.Println("'int(*detection.Confidence*100)'")
		if int(*detection.Confidence) >= minConfidence {
			detectedText := *detection.DetectedText
			//log.Println(detectedText)
			numbers := extractNumberFromString(detectedText)
			for _, number := range numbers {
				set.Add(number, *detection.Confidence)
			}
		}
	}

	log.Println("set.GetKeys()")
	var arr []picture.TextDetection

	for detectedText, confidence := range set.GetAll() {
		fmt.Printf("%s %f\n", detectedText, confidence)
		detection := picture.NewTextDetection(detectedText, confidence)
		arr = append(arr, *detection)
	}

	log.Println(set.GetKeys())

	output, err := json.Marshal(result)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	log.Println(string(output))
	return arr, nil
}

func extractNumberFromString(input string) []string {
	re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)

	//fmt.Printf("Pattern: %v\n", re.String()) // Print Pattern
	//
	//fmt.Printf("String contains any match: %v\n", re.MatchString(input)) // True

	submatchall := re.FindAllString(input, -1)
	var result []string
	for _, element := range submatchall {
		result = append(result, element)
		//log.Println(element)
	}
	return result
}
