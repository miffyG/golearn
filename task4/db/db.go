package db

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DbConfig struct {
	DBHost     string `env:"DB_HOST" required:"true" envdefault:"localhost"`
	DBPort     string `env:"DB_PORT" required:"true" envdefault:"3306"`
	DBUser     string `env:"DB_USER" required:"true" envdefault:"root"`
	DBPassword string `env:"DB_PASSWORD" required:"true" envdefault:"password"`
	DBName     string `env:"DB_NAME" required:"true" envdefault:"test"`
}

var GormDb *gorm.DB

func InitGormDb(cfg *DbConfig) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("gorm连接数据库失败:", err)
	}
	GormDb = db
	fmt.Println("gorm连接数据库成功")
}

func CloseDBConnections() {
	if GormDb != nil {
		sqlDB, err := GormDb.DB()
		if err == nil {
			sqlDB.Close()
		}
	}
}
