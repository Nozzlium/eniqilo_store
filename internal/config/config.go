package config

type Config struct {
	DB         DBConfig
	JWTSecret  string `json:"JWT_SECRET"`
	BCryptSalt uint8  `json:"BCRYPT_SALT"`
}

type DBConfig struct {
	DBName     string `json:"DB_NAME"`
	DBPort     string `json:"DB_PORT"`
	DBHost     string `json:"DB_HOST"`
	DBUsername string `json:"DB_USERNAME"`
	DBPassword string `json:"DB_PASSWORD"`
	DBParams   string `json:"DB_PARAMS"`
}
