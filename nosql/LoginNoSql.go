package nosql

import (
	"context"
	"errors"
	"log"
	"time"
	"virtual-account/config"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LoginNoSql struct {
	BaseNoSql `bson:"-"`

	ID          string    `json:"id" bson:"id"`
	Token       string    `json:"token" bson:"token"`
	AccountNo   string    `json:"accountNo" bson:"accountNo"`
	CreatedAt   time.Time `json:"createdAt" bson:"createdAt"`
	IsExpired   bool      `json:"isExpired" bson:"isExpired"`
	ExpiredTime time.Time `json:"expiredTime" bson:"expiredTime"`
}

func CreateLoginNoSql(client *mongo.Client) LoginNoSql {
	login := LoginNoSql{}
	login.Client = client
	return login
}

func (login LoginNoSql) Collection() *mongo.Collection {
	return login.Client.Database(config.DATABASE).Collection("logins")
}

func (login LoginNoSql) Truncate(ctx context.Context) error {
	return login.Collection().Drop(ctx)
}

func (login LoginNoSql) AddLogin(ctx context.Context, token string, accountNo string) error {

	login.ID = uuid.New().String()
	login.Token = token
	login.AccountNo = accountNo
	login.CreatedAt = time.Now().UTC()
	login.IsExpired = false
	login.ExpiredTime = time.Now().UTC().Add(time.Hour * 1)

	_, err := login.Collection().InsertOne(ctx, login)
	if err != nil {
		return err
	}
	return nil

}

func (login *LoginNoSql) FindOneByToken(ctx context.Context, token string) error {

	filter := bson.D{
		primitive.E{Key: "token", Value: token},
	}

	result := login.Collection().FindOne(ctx, filter)
	if result.Err() != nil {

		log.Println(result.Err())
		if result.Err().Error() == config.MONGO_NO_DOCUMENT {
			return errors.New("login_not_found")
		}
		return result.Err()
	}

	err := result.Decode(&login)
	if err != nil {
		return err
	}

	return nil
}

func (login LoginNoSql) Delete(ctx context.Context, id string) error {

	filter := bson.D{
		primitive.E{Key: "id", Value: id},
	}

	result, err := login.Collection().DeleteOne(ctx, filter, nil)
	if err != nil {
		return err
	}
	log.Println(result)
	return nil

}

func (login LoginNoSql) ListLogins(ctx context.Context) ([]LoginNoSql, error) {

	logins := []LoginNoSql{}
	filter := bson.D{}
	cursor, err := login.Collection().Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &logins)
	if err != nil {
		return nil, err
	}

	return logins, nil

}

func (login *LoginNoSql) UpdateExpired(ctx context.Context, accountNo string) error {
	filter := bson.D{
		primitive.E{Key: "accountNo", Value: accountNo},
	}

	set := bson.D{
		primitive.E{
			Key: "$set",
			Value: bson.D{
				primitive.E{Key: "token", Value: ""},
				primitive.E{Key: "isExpired", Value: true},
				primitive.E{Key: "updatedAt", Value: time.Now().UTC()},
			},
		},
	}

	optionsAfter := options.After
	updateOptions := &options.FindOneAndUpdateOptions{
		ReturnDocument: &optionsAfter,
	}

	result := login.Collection().FindOneAndUpdate(ctx, filter, set, updateOptions)
	if result.Err() != nil {
		return result.Err()
	}

	return result.Decode(&login)

}
