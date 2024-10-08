package database

import (
	"fmt"
	"rincon/config"
	"rincon/model"
	"rincon/utils"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitializeSQL() {
	db, err := SelectDB()
	if err != nil {
		if dbRetries < 5 {
			dbRetries++
			utils.SugarLogger.Errorln("Failed to connect database, retrying in 5s... ")
			time.Sleep(time.Second * 5)
			InitializeSQL()
		} else {
			utils.SugarLogger.Errorln("Failed to connect database after 5 attempts, defaulting to local storage")
			config.StorageMode = "local"
			InitializeLocal()
			return
		}
	} else {
		utils.SugarLogger.Infoln("Connected to database")
		err = db.AutoMigrate(&model.Service{}, &model.ServiceDependency{}, &model.Route{})
		if err != nil {
			utils.SugarLogger.Errorln("AutoMigration failed", err)
		}
		utils.SugarLogger.Infoln("AutoMigration complete")
		DB = db
	}
}

func SelectDB() (*gorm.DB, error) {
	if config.DatabaseDriver == "mysql" {
		return ConnectMysql()
	} else if config.DatabaseDriver == "postgres" {
		return ConnectPostgres()
	}
	return nil, fmt.Errorf("invalid database driver")
}

func ConnectPostgres() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC", config.DatabaseHost, config.DatabaseUser, config.DatabasePassword, config.DatabaseName, config.DatabasePort)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func ConnectMysql() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=UTC", config.DatabaseUser, config.DatabasePassword, config.DatabaseHost, config.DatabasePort, config.DatabaseName)
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}
