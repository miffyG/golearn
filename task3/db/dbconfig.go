package db

type DbConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

func LoadConfig() (DbConfig, error) {
	var cfg DbConfig

	cfg.Host = "127.0.0.1"
	cfg.Port = "3306"
	cfg.Name = "test"
	cfg.User = "root"
	cfg.Password = "123456"

	return cfg, nil
}
