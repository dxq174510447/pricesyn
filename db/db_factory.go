package db

import (
	"context"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
	"time"
)

type DbFactory struct {
	initLock sync.Once
	db       *gorm.DB
	DbUser   string
	DbPwd    string
	DbHost   string
	DbPort   int
	DbName   string
	MaxOpen  int
	MaxIdle  int
}

func (d *DbFactory) init() error {
	var err error
	d.initLock.Do(func() {
		var connStr = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
			d.DbUser, d.DbPwd, d.DbHost, d.DbPort, d.DbName)
		mysqlDB, err1 := gorm.Open(mysql.Open(connStr))
		if err1 != nil {
			err = err1
		} else {
			db, _ := mysqlDB.DB()
			db.SetMaxIdleConns(d.MaxIdle)
			db.SetMaxOpenConns(d.MaxOpen)
			db.SetConnMaxLifetime(time.Duration(30) * time.Minute)
			d.db = mysqlDB
		}
	})
	return err
}

func (d *DbFactory) GetDb(ctx context.Context) (*gorm.DB, error) {
	err := d.init()
	if err != nil {
		return nil, err
	}
	return d.db, nil
}
