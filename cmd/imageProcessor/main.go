package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/streadway/amqp"
	"net/http"
	_ "net/http/pprof"
	"photofinish/pkg/common/util"
	"runtime"
	"time"

	"github.com/aws/aws-sdk-go/service/rekognition"
	"photofinish/pkg/app/aws/recognition"

	"github.com/pkg/errors"
	"log"
	"os"
	//"photofinish/pkg/app/aws/recognition/rekognition"
	s32 "photofinish/pkg/app/aws/s3"
	"photofinish/pkg/app/dropbox"
	"photofinish/pkg/app/picture"
	rabbitmq "photofinish/pkg/common/infrarstructure/amqp"
	"photofinish/pkg/common/infrarstructure/db"
	"photofinish/pkg/domain/pictures"
	"photofinish/pkg/infrastructure/postgres"
	"strconv"
	"sync"
)

var svc *rekognition.Rekognition

func main() {
	util.LoadEnvFileIfNeeded()

	go func() {
		err := http.ListenAndServe("0.0.0.0:8081", nil)
		if err != nil {
			log.Fatal(err)
		}
	}()
	amqpServerURL := os.Getenv("AMQP_SERVER_URL")
	if len(amqpServerURL) == 0 {
		log.Fatal(errors.New("failed get AMQP_SERVER_URL"))
	}

	awsS3Region := os.Getenv("AWS_S3_REGION")
	if len(awsS3Region) == 0 {
		log.Fatal(errors.New("Failed get AWS_S3_REGION"))
	}
	awsBucket := os.Getenv("AWS_S3_BUCKET")
	if len(awsBucket) == 0 {
		log.Fatal(errors.New("Failed get AWS_S3_BUCKET"))
	}

	accessToken := os.Getenv("DROPBOX_ACCESS_TOKEN")
	if len(accessToken) == 0 {
		log.Fatal("DROPBOX_ACCESS_TOKEN not set into env variable ")
	}

	dbDSN := os.Getenv("DATABASE_DSN")
	if len(dbDSN) == 0 {
		log.Fatal(errors.New("Failed get DATABASE_DSN"))
	}

	maxConnectionsStr := os.Getenv("DATABASE_MAX_CONNECTION")
	if len(maxConnectionsStr) == 0 {
		log.Fatal(errors.New("Failed get DATABASE_MAX_CONNECTION"))
	}

	maxConnections, err := strconv.Atoi(maxConnectionsStr)
	if err != nil {
		log.Fatal(errors.New("invalid DATABASE_MAX_CONNECTION"))
	}

	acquireTimeoutStr := os.Getenv("DATABASE_ACQUIRE_TIMEOUT")
	if len(acquireTimeoutStr) == 0 {
		log.Fatal(errors.New("Failed get DATABASE_MAX_CONNECTION"))
	}

	acquireTimeout, err := strconv.Atoi(acquireTimeoutStr)
	if err != nil {
		log.Fatal(errors.New("invalid DATABASE_MAX_CONNECTION"))
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(awsS3Region))
	if err != nil {
		log.Fatal(err)
	}
	client := s3.NewFromConfig(cfg)
	uploader := manager.NewUploader(client)
	connector, err := db.GetDBConfig(dbDSN, maxConnections, acquireTimeout)

	if err != nil {
		log.Fatal(err.Error())
	}
	pool, err := db.NewConnectionPool(connector)

	fmt.Println(pool)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(uploader)
	downloader := dropbox.NewSDKDownloader(accessToken)
	awsS3Uploader := s32.NewAwsS3Uploader(uploader, awsBucket)

	pictureRepo := postgres.NewPictureRepository(pool)
	compressor := picture.NewImageCompressor()
	textDetector := recognition.NewAmazonTextRecognition(svc)
	err = pictureRepo.IsExists("1")
	pictureCoordinator := picture.NewCoordinatorServiceImpl(2,
		pictureRepo,
		downloader,
		awsS3Uploader,
		textDetector,
		compressor,
		0)

	dial, err := rabbitmq.Dial(amqpServerURL, rabbitmq.TargetQueue)
	if err != nil {
		log.Fatal(err)
	}
	messages, err := rabbitmq.Consume(dial, rabbitmq.TargetQueue)
	if err != nil {
		log.Fatal(err)
	}
	forever := make(chan bool)
	wg := sync.WaitGroup{}
	wg.Add(1)
	ch := make(chan pictures.DropboxImage)
	go func() {
		defer wg.Done()
		for message := range messages {
			// For example, show received message in a console.
			handleMessage(message, err, ch)
			close(ch)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for t := 0; t < runtime.NumCPU(); t++ {
			wg.Add(1)
			go handleImageAsync(ch, &wg, pictureCoordinator)
		}
	}()

	wg.Wait()

	<-forever
}

func handleMessage(message amqp.Delivery, err error, ch chan pictures.DropboxImage) {
	log.Printf(" > Received message: %s\n", message.Body)
	var initial pictures.DropboxImages

	err = json.Unmarshal(message.Body, &initial)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(len(initial.Images))
	time.Sleep(2 * time.Second)
	//id := initial.EventId

	for _, image := range initial.Images {
		ch <- image
	}
}

func handleImageAsync(ch chan pictures.DropboxImage, wg *sync.WaitGroup, pictureCoordinator pictures.CoordinatorService) {
	for img := range ch {
		p := pictures.Picture{
			Id:              img.Id,
			DropboxPath:     img.Path,
			EventId:         img.EventId,
			IsOriginalSaved: false,
		}
		fmt.Println(p)
		err := pictureCoordinator.PerformAddImage(&p)
		if err != nil {
			log.Println(err)
		} else {
			fmt.Println("success")
			fmt.Println(p)
		}

	}
	wg.Done()
}

func init() {
	//Access keys are read from ~/.aws/credentials
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1"),
	})

	if err != nil {
		fmt.Println("Error while creating session,", err)
		return
	}

	svc = rekognition.New(sess)
	_ = svc
}
