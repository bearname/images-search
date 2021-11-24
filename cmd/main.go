package main

import (
    "aws_rekognition_demo/internal/app/auth"
    s32 "aws_rekognition_demo/internal/app/aws/s3"
    "aws_rekognition_demo/internal/app/dropbox"
    "aws_rekognition_demo/internal/app/event"
    "aws_rekognition_demo/internal/app/picture"
    "aws_rekognition_demo/internal/common/infrarstructure/server"
    "aws_rekognition_demo/internal/domain"
    "aws_rekognition_demo/internal/infrastructure/postgres"
    "aws_rekognition_demo/internal/infrastructure/router"
    "aws_rekognition_demo/internal/infrastructure/transport"
    "context"
    "fmt"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/feature/s3/manager"
    "github.com/aws/aws-sdk-go-v2/service/s3"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/rekognition"
    "github.com/jackc/pgx"
    "github.com/pkg/errors"
    log "github.com/sirupsen/logrus"
    "net"
    "net/http"
    "os"
    "strconv"
    "time"
)

var svc *rekognition.Rekognition

//var bucket = flag.String("bucket", "img-rekongnition-test", "The name of the bucket")
//var photo = flag.String("photo", "pexels-oleg-magni-1427741.jpg", "The path to the photo file (JPEG, JPG, PNG)")

func main() {

    //flag.Parse()
    //runtime.GOMAXPROCS(4)
    log.SetFormatter(&log.JSONFormatter{})
    file, err := os.OpenFile("short.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
    if err == nil {
        log.SetOutput(file)
        defer func(file *os.File) {
            err := file.Close()
            if err != nil {
                log.Error(err)
            }
        }(file)
    }

    conf, err := ParseConfig()
    if err != nil {
        log.Fatal("Default settings" + err.Error())
    }

    //awsS3Bucket := "img-rekongnition-test"
    //awsS3region := "eu-central-1"
    //conf := &Config{
    //    ServeRestAddress: ":8000",
    //    DbAddress:        "localhost:5432",
    //    DbName:           "photofinish",
    //    DbUser:           "postgres",
    //    DbPassword:       "postgres",
    //    MaxConnections:   10,
    //    AcquireTimeout:   1,
    //    AwsS3Region:      awsS3region,
    //    AwsS3Bucket:      awsS3Bucket,
    //}

    connector, err := getConnector(conf)

    if err != nil {
        log.Fatal(err.Error())
    }
    pool, err := newConnectionPool(connector)

    if err != nil {
        log.Fatal(err.Error())
    }

    log.WithFields(log.Fields{"url": conf.ServeRestAddress}).Info("starting the httpServer")
    getenv := os.Getenv("PORT")
    if len(getenv) == 0 {
        getenv = "8000"
    }
    httpServer := server.HttpServer{}
    killSignalChan := httpServer.GetKillSignalChan()

    cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(conf.AwsS3Region))
    if err != nil {
        log.Fatal(err)
    }
    client := s3.NewFromConfig(cfg)
    uploader := manager.NewUploader(client)
    accessToken := os.Getenv("DROPBOX_ACCESS_TOKEN")
    if len(accessToken) == 0 {
        log.Fatal("DROPBOX_ACCESS_TOKEN not setted into env variable ")
        //accessToken = "Hoih4cfQJFsAAAAAAAAAAUXzsZDZ2k8o74P9PhFED1VwJAGYZT_qQIQIBa7zFlsq"
    }
    downloader := dropbox.NewSDKDownloader(accessToken)
    awsS3Uploader := s32.NewAwsS3Uploader(uploader, conf.AwsS3Bucket)
    handler := initHandlers(pool, awsS3Uploader, downloader)

    srv := httpServer.StartServer(getenv, handler)
    httpServer.WaitForKillSignal(killSignalChan)
    err = srv.Shutdown(context.TODO())
    if err != nil {
        log.Error(err)
        return
    }
}

func initHandlers(connPool *pgx.ConnPool, awsUploaderManager domain.Uploader, downloader domain.Downloader) http.Handler {
    pictureRepos := postgres.NewPictureRepository(connPool)
    compressor := picture.NewImageCompressor()
    pictureService := picture.NewPictureService(pictureRepos, svc, awsUploaderManager, compressor)
    pictureController := transport.NewPictureController(pictureService, downloader)
    eventService := event.NewEventService(postgres.NewEventRepository(connPool))
    eventController := transport.NewEventController(eventService)

    repository := postgres.NewUserRepository(connPool)
    service := auth.NewAuthService(repository)
    controller := transport.NewAuthController(service)
    handler := router.Router(pictureController, eventController, controller)
    return handler
}

type Config struct {
    ServeRestAddress string
    DbAddress        string
    DbName           string
    DbUser           string
    DbPassword       string
    MaxConnections   int
    AcquireTimeout   int
    AwsS3Bucket      string
    AwsS3Region      string
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

//postgres://cxzhykstalbgyi:cc469d81d726131f1ae98d8f8cbff091dbd7e9be7cce46a51674efa320f05f32@ec2-54-173-31-84.compute-1.amazonaws.com:5432/d7r53jijhba21i
func ParseConfig() (*Config, error) {
    var err error
    serveRestAddress, err := parseEnvString("SERVE_REST_ADDRESS", err)
    dbAddress, err := parseEnvString("DATABASE_ADDRESS", err)
    dbName, err := parseEnvString("DATABASE_NAME", err)
    dbUser, err := parseEnvString("DATABASE_USER", err)
    dbPassword, err := parseEnvString("DATABASE_PASSWORD", err)
    maxConnections, err := parseEnvInt("DATABASE_MAX_CONNECTION", err)
    acquireTimeout, err := parseEnvInt("DATABASE_ACQUIRE_TIMEOUT", err)
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
    }, nil
}

func getConnector(config *Config) (pgx.ConnPoolConfig, error) {
    databaseUri := "postgres://" + config.DbUser + ":" + config.DbPassword + "@" + config.DbAddress + "/" + config.DbName
    log.Info("databaseUri: " + databaseUri)
    pgxConnConfig, err := pgx.ParseURI(databaseUri)
    if err != nil {
        return pgx.ConnPoolConfig{}, errors.Wrap(err, "failed to parse database URI from environment variable")
    }
    pgxConnConfig.Dial = (&net.Dialer{Timeout: 10 * time.Second, KeepAlive: 5 * time.Minute}).Dial
    pgxConnConfig.RuntimeParams = map[string]string{
        "standard_conforming_strings": "on",
    }
    pgxConnConfig.PreferSimpleProtocol = true

    return pgx.ConnPoolConfig{
        ConnConfig:     pgxConnConfig,
        MaxConnections: config.MaxConnections,
        AcquireTimeout: time.Duration(config.AcquireTimeout) * time.Second,
    }, nil
}

func newConnectionPool(config pgx.ConnPoolConfig) (*pgx.ConnPool, error) {
    return pgx.NewConnPool(config)
}

func init() {
    //https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/sessions.html

    //flag.StringVar(bucket, "b", "img-rekongnition-test", "The name of the bucket")
    //flag.StringVar(photo, "p", "pexels-oleg-magni-1427741.jpg", "The path to the photo file (JPEG, JPG, PNG)")
    //flag.Parse()
    //
    //if *bucket == "" || *photo == "" {
    //	checks = false
    //	flag.PrintDefaults()
    //	fmt.Println("You must supply a bucket name (-b BUCKET) and photo file (-p PHOTO)")
    //	return
    //}
    //
    //fileExtension := filepath.Ext(*photo)
    //validExtension := map[string]bool{
    //    ".png":  true,
    //    ".jpg":  true,
    //    ".jpeg": true,
    //}
    //
    //if !validExtension[fileExtension] {
    //
    //    fmt.Println("Rekognition only supports jpeg, jpg or png")
    //    return
    //}

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
