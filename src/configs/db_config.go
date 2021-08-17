package configs

import (
	"fmt"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	"github.com/golang-migrate/migrate/source/file"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	gormmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"io/ioutil"
)

type DbConfig struct {
	Host     string
	Port     string
	DbName   string
	User     string
	Pass     string
	PassFile string
}

//func InitDB(host, port, dbName, user, pass, passFile string) (*gorm.DB, error) {
//InitDB Opens new connection with Mysql
func InitDB(config *DbConfig) (*gorm.DB, error) {
	if config.PassFile != "" {
		b, err := ioutil.ReadFile(config.PassFile)
		if err != nil {
			return nil, err
		}
		config.Pass = string(b)
	}
	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=true", config.User, config.Pass, config.Host, config.Port, config.DbName)
	log.Info(dataSource)
	db, err := gorm.Open(gormmysql.Open(dataSource), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	sqlDB, _ := db.DB()
	driver, err := mysql.WithInstance(sqlDB, &mysql.Config{})
	if err != nil {
		return nil, err
	}
	fsrc, err := (&file.File{}).Open("file://migrations")
	if err != nil {
		return nil, err
	}
	m, err := migrate.NewWithInstance(
		"file",
		fsrc,
		config.DbName,
		driver,
	)
	if err != nil {
		return nil, errors.Errorf("An error occurred while migration instantiating.. %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, errors.Errorf("An error occurred while syncing the database.. %v", err)
	}

	return db, err
}
