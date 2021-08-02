package services

import (
	"database/sql"

	"go.mongodb.org/mongo-driver/mongo"
)

type BaseService struct {
	DB          *sql.DB       `json:",omitempty"`
	MongoClient *mongo.Client `json:",omitempty"`
}
