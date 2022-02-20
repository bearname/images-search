package main

import (
	"errors"
	"github.com/jackc/pgx"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"photofinish/pkg/common/infrarstructure/db"
	"photofinish/pkg/common/util"
	"photofinish/pkg/domain/pictures"
	"strconv"
	"time"
)

func main() {
	util.LoadEnvFileIfNeeded()

	dsn := os.Getenv("DATABASE_DSN")
	if len(dsn) == 0 {
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

	connector, err := db.GetDBConfig(dsn, maxConnections, acquireTimeout)

	if err != nil {
		log.Fatal(err.Error())
	}
	pool, err := db.NewConnectionPool(connector)
	if err != nil {
		log.Fatal(err.Error())
	}

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":2112", nil)
	}()
	for {
		log.Println("getSqlReleaseBlockPicture()")
		err = db.WithTransaction(pool, func(tx *pgx.Tx) error {
			var data []interface{}
			data = append(data, pictures.Failed, pictures.Processing)
			_, err = tx.Exec(getSqlReleaseBlockPicture(), data...)
			if err != nil {
				log.Println(err)
			}
			return err
		})
		if err != nil {
			log.Error(err)
		}
		time.Sleep(10 * time.Minute)
	}
}

func getSqlReleaseBlockPicture() string {
	return `UPDATE pictures
            SET processing_status = $1, execute_after = NOW()
            WHERE (execute_after < NOW() -  '10 minutes'::interval OR execute_after IS NULL) AND processing_status = $2;`
}
