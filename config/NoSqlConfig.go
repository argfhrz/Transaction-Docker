package config

const (
	DEV2 = "DEV2"

	CURRENT_PHASE2 = DEV2
)

type MongoConfig struct {
	Host           string
	Port           string
	User           string
	Pwd            string
	Authentication bool
}

const (
	DATABASE = "zfinaltest"
)

var MONGO_CONFIGS map[string]MongoConfig = map[string]MongoConfig{
	DEV2: {
		Authentication: true,
		User:           "admin",
		Pwd:            "password",
		Host:           "10.116.0.2",
		Port:           "27017",
	},
}
