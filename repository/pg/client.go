package pg

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/vgheri/gennaker/engine"
)

type pgRepository struct {
	db *sql.DB
}

func (pg pgRepository) getDB() (*sql.DB, error) {
	if pg.db == nil {
		return nil, fmt.Errorf("Connection is not initialized")
	}
	return pg.db, nil
}

// NewClient returns a new postgres client
func NewClient(host, port, username, password, dbname string, maxconn int) (engine.DeploymentRepository, error) {
	var err error
	var dsn string
	var conn *sql.DB
	if password == "" {
		dsn = fmt.Sprintf("host=%s port=%s user=%s dbname=%s connect_timeout=10 statement_timeout=3000", host, port, username, dbname)
	} else {
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s connect_timeout=10 statement_timeout=3000", host, port, username, password, dbname)
	}

	connected := false
	for retries := 1; retries <= 5; retries++ {
		conn, err = sql.Open("postgres", dsn)
		if err == nil {
			err = conn.Ping()
			if err == nil {
				conn.SetMaxOpenConns(maxconn)
				conn.SetMaxIdleConns(maxconn)
				connected = true
				break
			}
		}
		fmt.Println("Error trying to connect to db. Retrying...")
		time.Sleep(time.Duration(retries) * time.Second)
	}
	if !connected {
		return nil, err
	}

	return &pgRepository{
		db: conn,
	}, nil
}
