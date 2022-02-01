package db

import (
	"github.com/jackc/pgx"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"net"
	"time"
)

func GetUrl(DbUser string, DbPassword string, DbAddress string, DbName string) string {
	return "postgres://" + DbUser + ":" + DbPassword + "@" + DbAddress + "/" + DbName + "?sslMode=false"
}

func GetDBConfig(databaseUri string, maxConnections int, acquireTimeout int) (pgx.ConnPoolConfig, error) {
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
		MaxConnections: maxConnections,
		AcquireTimeout: time.Duration(acquireTimeout) * time.Second,
	}, nil
}

func NewConnectionPool(config pgx.ConnPoolConfig) (*pgx.ConnPool, error) {
	return pgx.NewConnPool(config)
}
