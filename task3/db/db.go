package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var GormDb *gorm.DB
var SqlxDb *sqlx.DB

func InitGormDb(cfg *DbConfig) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("gorm连接数据库失败:", err)
	}
	GormDb = db
	fmt.Println("gorm连接数据库成功")
}

func InitSqlxDb(cfg *DbConfig) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name)
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		fmt.Println("sqlx连接数据库失败:", err)
		return
	}
	SqlxDb = db
	fmt.Println("sqlx连接数据库成功")
}

func CloseDBConnections() {
	if GormDb != nil {
		sqlDB, err := GormDb.DB()
		if err == nil {
			sqlDB.Close()
		}
	}
	if SqlxDb != nil {
		SqlxDb.Close()
	}
}
