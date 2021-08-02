package config

const (
	DEV = "DEV"

	CURRENT_PHASE = DEV
)

type DbConfig struct {
	Driver   string
	User     string
	Password string
	Host     string
	DbName   string
	SslMode  string
}

var DB_CONFIGS map[string]DbConfig = map[string]DbConfig{
	DEV: {
		Driver:   "postgres",
		User:     "postgres",
		Password: "123456",
		Host:     "10.116.0.2",
		DbName:   "arigo",
		SslMode:  "disable",
	},
}
