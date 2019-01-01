package db

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"gocleancode/config"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

type Db interface {
	Close() error
	Transact(txFunc func(*sql.Tx) error) (err error)
}

type DB struct {
	*sql.DB
	DataSourceName string
}

var instances []Db

func init() {
	closeAllDb := make(chan os.Signal)
	signal.Notify(closeAllDb,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		<-closeAllDb
		for _, db := range instances {
			db.Close()
		}
	}()
}
func NewMysqlDb(config config.Configuration) (DB, error) {
	ds := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", config.DbUser, config.DbPass, config.DbHost, strconv.Itoa(config.DbPort), config.DbName)

	db, err := sql.Open("mysql", ds)
	_db := DB{db, ds}

	instances = append(instances, _db)

	if err != nil {
		return _db, err
	}

	err = db.Ping() // Check if can connect to db
	if err != nil {
		log.Error(err)
		return _db, errors.New("Cannot connect to " + ds)
	} else {
		log.Info(fmt.Sprintf("Successfully connected to %s:%v/%s", config.DbHost, config.DbPort, config.DbName))
	}
	return _db, err
}

func (db DB) Close() error {
	log.Info("Closing db " + db.DataSourceName)
	return db.DB.Close()
}

func (db DB) Transact(txFunc func(*sql.Tx) error) (err error) {
	tx, err := db.Begin()
	if err != nil {
		log.Error(err)
		return err
	}
	defer func() {
		if p := recover(); p != nil { // Under normal circumstances a panic should not occur.
			tx.Rollback() //  Should panic really happen, make sure we rollback
			panic(p)      // Then re-throw panic after Rollback
		} else if err != nil {
			tx.Rollback() // err is non-nil; don't change it
		} else {
			err = tx.Commit() // err is nil; if Commit returns error update err
		}
	}()
	err = txFunc(tx)
	return err
}
