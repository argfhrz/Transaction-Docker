package connection

import (
	"context"
	"fmt"
	"virtual-account/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func OpenMongoDb(ctx context.Context) (*mongo.Client, error) {
	config := config.MONGO_CONFIGS[config.DEV2]

	connString := fmt.Sprintf("mongodb://%s:%s@%s:%s", config.User, config.Pwd, config.Host, config.Port)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connString))
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	return client, nil
}

func CloseMongoDb(ctx context.Context, client *mongo.Client) error {

	err := client.Disconnect(ctx)
	if err != nil {
		return err
	}
	return nil
}
