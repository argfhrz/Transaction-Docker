package nosql

import "go.mongodb.org/mongo-driver/mongo"

type BaseNoSql struct {
	Client *mongo.Client `json:"-" bson:"-"`
}
