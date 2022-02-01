package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rekognition"
	"github.com/jackc/pgx"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"net/http"
	"os"
	"photofinish/pkg/app/auth"
	"photofinish/pkg/app/dropbox"
	"photofinish/pkg/app/event"
	"photofinish/pkg/app/paySystem"
	"photofinish/pkg/app/picture"
	"photofinish/pkg/app/tasks"
	"photofinish/pkg/app/user"
	rabbitmq "photofinish/pkg/common/infrarstructure/amqp"
	"photofinish/pkg/common/infrarstructure/db"
	"photofinish/pkg/common/infrarstructure/server"
	"photofinish/pkg/common/util"
	"photofinish/pkg/infrastructure/postgres"
	"photofinish/pkg/infrastructure/router"
	"photofinish/pkg/infrastructure/transport"
	"runtime"
	"strconv"
	"time"
)

var svc *rekognition.Rekognition

func main() {
	util.LoadEnvFileIfNeeded()

	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetFormatter(&log.JSONFormatter{})
	file, err := os.OpenFile("short.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err == nil {
		log.SetOutput(file)
		defer func(file *os.File) {
			err = file.Close()
			if err != nil {
				log.Error(err)
			}
		}(file)
	}

	c, err := ParseConfig()
	if err != nil {
		log.Fatal("Default settings" + err.Error())
	}

	conf := c
	fmt.Println("\n\n\n\n\n\n\n\nconf.DropboxAccessToken")
	fmt.Println(conf.DropboxAccessToken)
	fmt.Println("\n\n\n\n\n\n\n\nconf.DropboxAccessToken")
	url := db.GetUrl(conf.DbUser, conf.DbPassword, conf.DbAddress, conf.DbName)
	connector, err := db.GetDBConfig(url, conf.MaxConnections, conf.AcquireTimeout)

	if err != nil {
		log.Fatal(err.Error())
	}
	pool, err := db.NewConnectionPool(connector)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.WithFields(log.Fields{"url": conf.ServeRestAddress}).Info("starting the httpServer")
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8000"
	}
	httpServer := server.HttpServer{}
	killSignalChan := httpServer.GetKillSignalChan()

	amqpChannel, err := rabbitmq.Dial(conf.AmqpServerURL, rabbitmq.TargetQueue)
	if err != nil {
		log.Fatal("Failed connect to amqp server")
	}
	if amqpChannel == nil {
		log.Fatal("amqpChannel nil")
		return
	}
	defer func(amqpChannel *amqp.Channel) {
		err = amqpChannel.Close()
		if err != nil {
			return
		}
	}(amqpChannel)

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(c.AwsS3Region))
	if err != nil {
		log.Fatal(err)
	}
	client := s3.NewFromConfig(cfg)

	handler := initHandlers(pool, amqpChannel, conf.DropboxAccessToken, conf.StripeSecretKey, client, c.AwsS3Bucket)

	log.Println("Start on port '" + port + " 'at " + time.Now().String())
	srv := httpServer.StartServer(port, handler)
	httpServer.WaitForKillSignal(killSignalChan)
	err = srv.Shutdown(context.TODO())

	log.Println("Stop at " + time.Now().String())

	if err != nil {
		log.Error(err)
		return
	}
}

func initHandlers(connPool *pgx.ConnPool, amqpChannel *amqp.Channel, dropboxAccessToken string, stripeSecretKey string, s3Client *s3.Client, bucket string) http.Handler {
	downloader := dropbox.NewSDKDownloader(dropboxAccessToken)

	pictureRepo := postgres.NewPictureRepository(connPool)
	pictureService := picture.NewPictureService(pictureRepo, amqpChannel, downloader, s3Client, bucket) //, svc, uploader, compressor
	pictureController := transport.NewPictureController(pictureService)

	eventRepo := postgres.NewEventRepository(connPool)
	eventService := event.NewEventService(eventRepo)
	eventController := transport.NewEventController(eventService)

	userRepo := postgres.NewUserRepository(connPool)
	authService := auth.NewAuthService(userRepo)
	userService := user.NewUserService(userRepo)
	authController := transport.NewAuthController(authService)

	tasksRepo := postgres.NewTasksRepositoryImpl(connPool)
	tasksService := tasks.NewService(tasksRepo)
	tasksController := transport.NewTasksController(tasksService)

	orderRepo := postgres.NewOrderRepository(connPool)
	stripeService := paySystem.NewStripeService(stripeSecretKey)
	payService := paySystem.NewOrderService(orderRepo)
	yookassaService, err := paySystem.NewYookassaService(54401, "test_Fh8hUAVVBGUGbjmlzba6TB0iyUbos_lueTHE-axOwM0")
	if err != nil {
		log.Fatal(err)
		return nil
	}

	orderController := transport.NewOrderController(userService, payService, stripeService, yookassaService)

	controllers := router.Controllers{
		PictureController: pictureController,
		EventsController:  eventController,
		AuthController:    authController,
		TasksController:   tasksController,
		OrderController:   orderController,
	}
	return router.Router(controllers)
}

type Config struct {
	ServeRestAddress   string
	DbAddress          string
	DbName             string
	DbUser             string
	DbPassword         string
	MaxConnections     int
	AcquireTimeout     int
	AwsS3Region        string
	AwsS3Bucket        string
	AmqpServerURL      string
	DropboxAccessToken string
	StripeSecretKey    string
}

func parseEnvString(key string, err error) (string, error) {
	if err != nil {
		return "", err
	}
	str, ok := os.LookupEnv(key)
	if !ok {
		return "", fmt.Errorf("undefined environment variable %v", key)
	}
	return str, nil
}

func parseEnvInt(key string, err error) (int, error) {
	s, err := parseEnvString(key, err)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(s)
}

func ParseConfig() (*Config, error) {
	var err error
	serveRestAddress, err := parseEnvString("SERVE_REST_ADDRESS", err)
	dbAddress, err := parseEnvString("DATABASE_ADDRESS", err)
	dbName, err := parseEnvString("DATABASE_NAME", err)
	dbUser, err := parseEnvString("DATABASE_USER", err)
	dbPassword, err := parseEnvString("DATABASE_PASSWORD", err)
	maxConnections, err := parseEnvInt("DATABASE_MAX_CONNECTION", err)
	acquireTimeout, err := parseEnvInt("DATABASE_ACQUIRE_TIMEOUT", err)
	amqpServerUrl, err := parseEnvString("AMQP_SERVER_URL", err)
	dropboxAccessToken, err := parseEnvString("DROPBOX_ACCESS_TOKEN", err)
	stripeSecretKey, err := parseEnvString("STRIPE_SECRET_KEY", err)
	awsS3Region, err := parseEnvString("AWS_S3_REGION", err)
	awsS3Bucket, err := parseEnvString("AWS_S3_BUCKET", err)

	if err != nil {
		log.Info("error " + err.Error())
		return nil, err
	}

	return &Config{
		serveRestAddress,
		dbAddress,
		dbName,
		dbUser,
		dbPassword,
		maxConnections,
		acquireTimeout,
		awsS3Region,
		awsS3Bucket,
		amqpServerUrl,
		dropboxAccessToken,
		stripeSecretKey,
	}, nil
}

func init() {
	//https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/sessions.html

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
