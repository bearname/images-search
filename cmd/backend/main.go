package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rekognition"
	"github.com/col3name/images-search/pkg/app/auth"
	"github.com/col3name/images-search/pkg/app/event"
	"github.com/col3name/images-search/pkg/app/picture"
	"github.com/col3name/images-search/pkg/app/tasks"
	"github.com/col3name/images-search/pkg/app/user"
	demon2 "github.com/col3name/images-search/pkg/common/demon"
	rabbitmq "github.com/col3name/images-search/pkg/common/infrarstructure/amqp"
	"github.com/col3name/images-search/pkg/common/infrarstructure/db"
	"github.com/col3name/images-search/pkg/common/infrarstructure/server"
	"github.com/col3name/images-search/pkg/common/util"
	"github.com/col3name/images-search/pkg/infrastructure/dropbox"
	paySystem2 "github.com/col3name/images-search/pkg/infrastructure/paySystem"
	"github.com/col3name/images-search/pkg/infrastructure/postgres"
	"github.com/col3name/images-search/pkg/infrastructure/router"
	"github.com/col3name/images-search/pkg/infrastructure/transport"
	"github.com/jackc/pgx"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"
)

var svc *rekognition.Rekognition

func main() {
	util.LoadEnvFileIfNeeded()

	runtime.GOMAXPROCS(runtime.NumCPU())

	util.LogToFileIfNeeded()

	c, err := ParseConfig()
	if err != nil {
		log.Fatal("Default settings" + err.Error())
	}

	conf := c
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

	handler := initHandlers(pool, amqpChannel, conf.DropboxAccessToken, conf.StripeSecretKey, "addImageTopic")
	orderRepo := postgres.NewOutboxRepo(pool)
	amqpService := rabbitmq.NewAmqpService(amqpChannel)
	go demon2.HandleDemon(orderRepo, amqpService)
	log.Println("Start on port '" + port + " 'at " + time.Now().String())
	srv := httpServer.StartServer(port, handler)

	killSignalChan := httpServer.GetKillSignalChan()
	httpServer.WaitForKillSignal(killSignalChan)
	err = srv.Shutdown(context.TODO())
	log.Println("Stop at " + time.Now().String())

	if err != nil {
		log.Error(err)
		return
	}
}

func initHandlers(connPool *pgx.ConnPool, amqpChannel *amqp.Channel, dropboxAccessToken string, stripeSecretKey string, brokerAddImageTopic string) http.Handler {
	downloader := dropbox.NewSDKDownloader(dropboxAccessToken)

	pictureRepo := postgres.NewPictureRepository(connPool)
	tasksRepo := postgres.NewTasksRepo(connPool)
	tasksService := tasks.NewService(tasksRepo)

	amqpService := rabbitmq.NewAmqpService(amqpChannel)
	pictureService := picture.NewPictureService(pictureRepo, amqpService, downloader, brokerAddImageTopic, tasksService)
	pictureController := transport.NewPictureController(pictureService)

	eventRepo := postgres.NewEventRepository(connPool)
	eventService := event.NewEventService(eventRepo)
	eventController := transport.NewEventController(eventService)

	userRepo := postgres.NewUserRepository(connPool)
	authService := auth.NewAuthService(userRepo)
	userService := user.NewUserService(userRepo)
	authController := transport.NewAuthController(authService)

	tasksController := transport.NewTasksController(tasksService)

	orderRepo := postgres.NewOrderRepository(connPool)
	stripeService := paySystem2.NewStripeService(stripeSecretKey)
	payService := paySystem2.NewOrderService(orderRepo)
	yookassaService, err := paySystem2.NewYookassaService(54401, "test_Fh8hUAVVBGUGbjmlzba6TB0iyUbos_lueTHE-axOwM0")
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
	return router.Router(&controllers)
}

type Config struct {
	ServeRestAddress   string
	DbAddress          string
	DbName             string
	DbUser             string
	DbPassword         string
	MaxConnections     int
	AcquireTimeout     int
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
		log.Println("Error while creating session,", err)
		return
	}

	svc = rekognition.New(sess)
	_ = svc
}
