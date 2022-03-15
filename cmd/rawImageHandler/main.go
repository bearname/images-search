package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rekognition"
	"github.com/col3name/images-search/pkg/infrastructure/external/aws/recognition"
	awsS3 "github.com/col3name/images-search/pkg/infrastructure/external/aws/s3"
	"github.com/col3name/images-search/pkg/infrastructure/external/dropbox"
	"github.com/col3name/images-search/pkg/infrastructure/external/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"

	"github.com/col3name/images-search/pkg/app/picture"
	"github.com/col3name/images-search/pkg/common/infrarstructure/db"
	"github.com/col3name/images-search/pkg/common/util"
	"github.com/col3name/images-search/pkg/domain/pictures"
	"github.com/col3name/images-search/pkg/infrastructure/postgres"
	"strconv"
	"sync"
	"time"
)

var svc *rekognition.Rekognition

func main() {
	util.LoadEnvFileIfNeeded()

	dbDSN := os.Getenv("DATABASE_DSN")
	if len(dbDSN) == 0 {
		log.Fatal("Failed get DATABASE_DSN environment variable")
	}

	maxConnectionsStr := os.Getenv("DATABASE_MAX_CONNECTION")
	if len(maxConnectionsStr) == 0 {
		log.Fatal(errors.New("failed get DATABASE_MAX_CONNECTION"))
	}

	maxConnections, err := strconv.Atoi(maxConnectionsStr)
	if err != nil {
		log.Fatal(errors.New("invalid DATABASE_MAX_CONNECTION"))
	}

	acquireTimeoutStr := os.Getenv("DATABASE_ACQUIRE_TIMEOUT")
	if len(acquireTimeoutStr) == 0 {
		log.Fatal(errors.New("failed get DATABASE_MAX_CONNECTION"))
	}

	acquireTimeout, err := strconv.Atoi(acquireTimeoutStr)
	if err != nil {
		log.Fatal(errors.New("invalid DATABASE_MAX_CONNECTION"))
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

	tgBotToken := os.Getenv("TG_BOT_TOKEN")
	if len(tgBotToken) == 0 {
		log.Fatal(errors.New("Failed get TG_BOT_TOKEN"))
	}
	tgSupportChatIdStr := os.Getenv("TG_SUPPORT_CHAT_ID")
	if len(tgBotToken) == 0 {
		log.Fatal(errors.New("Failed get TG_BOT_TOKEN"))
	}
	tgSupportChatId, err := strconv.Atoi(tgSupportChatIdStr)
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
	if err != nil {
		log.Fatal(err.Error())
	}

	downloader := dropbox.NewSDKDownloader(accessToken)
	awsS3Uploader := awsS3.NewAwsS3Uploader(uploader, awsBucket)

	pictureRepo := postgres.NewPictureRepository(pool)
	compressor := picture.NewImageCompressor()
	textDetector := recognition.NewAmazonTextRecognition(svc)

	if err != nil {
		log.Fatal(err.Error())
	}

	bot, err := tgbotapi.NewBotAPI(tgBotToken)
	if err != nil {
		log.Fatal(err.Error())
	}
	tgNotifier := telegram.NewNotifier(bot, tgSupportChatId)

	pictureCoordinator := picture.NewCoordinatorServiceImpl(2,
		pictureRepo,
		downloader,
		awsS3Uploader,
		textDetector,
		compressor,
		tgNotifier,
		0)

	wg := sync.WaitGroup{}
	wg.Add(1)
	ch := make(chan pictures.Picture)
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":2112", nil)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for t := 0; t < 10; t++ {
			wg.Add(1)
			go handleImageAsync(ch, &wg, pictureCoordinator)
		}
	}()

	sql := getFailedImageSQL()

	cron(5*time.Minute, func() error {
		defer wg.Done()
		pictureList, err := getUnhandledPictures(pool, sql)
		if err != nil {
			log.Error(err)
			return err
		}

		for _, image := range pictureList {
			ch <- image
		}
		close(ch)

		return nil
	})
	if err != nil {
		log.Println(err)
	}

	wg.Wait()
}

func getUnhandledPictures(pool *pgx.ConnPool, sql string) ([]pictures.Picture, error) {
	var pictureList []pictures.Picture
	err := db.WithTransaction(pool, func(tx *pgx.Tx) error {
		var data []interface{}
		data = append(data, pictures.Processing, pictures.Failed)
		rows, err := tx.Query(sql, data...)
		if err != nil {
			return err
		}
		if rows.Err() == pgx.ErrNoRows {
			return err
		}

		var img pictures.Picture
		for rows.Next() {
			err = rows.Scan(&img.Id,
				&img.EventId,
				&img.DropboxPath,
				&img.IsOriginalSaved,
				&img.IsPreviewSaved,
				&img.IsTextRecognized,
				&img.ProcessingStatus,
				&img.Attempts,
				&img.ExecuteAfter,
				&img.UpdatedAt,
				&img.TaskId)
			if err != nil {
				return err
			}

			pictureList = append(pictureList, img)
		}

		return nil
	})
	return pictureList, err
}

func getFailedImageSQL() string {
	return `UPDATE pictures
SET processing_status = $1, update_at = NOW(), execute_after = NOW()
WHERE id IN (
    SELECT id
    FROM pictures
    WHERE processing_status = $2 AND (execute_after < NOW() OR execute_after IS NULL) AND attempts < 10
    LIMIT 100
        FOR UPDATE SKIP LOCKED
) RETURNING id, eventid, dropbox_path, is_original_saved, is_preview_saved, is_text_recognized, processing_status, attempts, execute_after, update_at, task_id;`
}

func cron(duration time.Duration, fn func() error) {
	for {
		time.Sleep(duration)
		err := fn()
		if err != nil {
			log.Error(err)
		}
	}
}

var rw sync.Mutex
var i int

func handleImageAsync(ch chan pictures.Picture, wg *sync.WaitGroup, pictureCoordinator pictures.CoordinatorService) {
	for img := range ch {
		err := pictureCoordinator.PerformAddImage(&img)
		if err != nil {
			log.Println("error", err)
		} else {
			log.Println("success", img)
		}

		rw.Lock()
		i++
		rw.Unlock()
		log.Println(i)
	}

	wg.Done()
}

func init() {
	//Access keys are read from ~/.aws/credentials
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1"),
	})

	if err != nil {
		log.Println("Error while creating session,", err)
		return
	}

	svc = rekognition.New(sess)
	_ = svc
}
