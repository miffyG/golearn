package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/miffyG/golearn/task4/db"
)

type Secret struct {
	JwtSecret string `env:"JWT_SECRET" unset:"true"`
}

func LoadConfig() {

	if err := godotenv.Load(".env"); err != nil {
		fmt.Printf("no .env loaded: %v", err)
	}
	fmt.Println("加载环境变量成功")
}

func GetDbConfig() *db.DbConfig {
	// 将环境变量映射到结构体（自动类型转换 + 默认值）
	var cfg db.DbConfig
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("解析配置失败: %v\n", err)
		return nil
	}

	fmt.Printf("数据库host：%s, 端口：%s, 用户：%s, 密码：%s, 数据库名：%s\n",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	return &cfg
}

func GetSecretConfig() *Secret {
	var secret Secret
	if err := env.Parse(&secret); err != nil {
		fmt.Printf("解析JWT密钥失败: %v\n", err)
		return nil
	}
	return &secret
}
