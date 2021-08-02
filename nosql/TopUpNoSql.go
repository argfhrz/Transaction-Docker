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

type TopUpNoSql struct {
	BaseNoSql `bson:"-"`

	ID            string    `json:"id" bson:"id"`
	BankCode      string    `json:"bankCode" bson:"bankCode"`
	BankName      string    `json:"bankName" bson:"bankName"`
	VaAccountNo   string    `json:"vaAccountNo" bson:"vaAccountNo"`
	VaAccountName string    `json:"vaAccountName" bson:"vaAccountName"`
	TopUpAmount   float64   `json:"topUpAmount" bson:"topUpAmount"`
	CreatedAt     time.Time `json:"createdAt" bson:"createdAt"`
}

func CreateTopUpNoSql(client *mongo.Client) TopUpNoSql {
	topUp := TopUpNoSql{}
	topUp.Client = client
	return topUp
}

func (topUp TopUpNoSql) Collection() *mongo.Collection {
	return topUp.Client.Database(config.DATABASE).Collection("top_ups")
}

func (topUp TopUpNoSql) Truncate(ctx context.Context) error {
	return topUp.Collection().Drop(ctx)
}

func (topUp TopUpNoSql) ListTopUp(ctx context.Context) ([]TopUpNoSql, error) {

	topUps := []TopUpNoSql{}
	filter := bson.D{}
	cursor, err := topUp.Collection().Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &topUps)
	if err != nil {
		return nil, err
	}

	return topUps, nil

}

func (topUp TopUpNoSql) AddTopUp(ctx context.Context,
	bankCode string, bankName string, vaAccountNo string, vaAccountName string, topUpAmount float64) (*TopUpNoSql, error) {

	topUp.ID = uuid.New().String()
	topUp.BankCode = bankCode
	topUp.BankName = bankName
	topUp.VaAccountNo = vaAccountNo
	topUp.VaAccountName = vaAccountName
	topUp.TopUpAmount = topUpAmount
	topUp.CreatedAt = time.Now().UTC()

	_, err := topUp.Collection().InsertOne(ctx, topUp)
	if err != nil {
		return nil, err
	}
	return &topUp, nil

}

func (topUp *TopUpNoSql) FindOneByID(ctx context.Context, id string) error {

	filter := bson.D{
		primitive.E{Key: "id", Value: id},
	}

	result := topUp.Collection().FindOne(ctx, filter)

	if result.Err() != nil {

		log.Println(result.Err())
		if result.Err().Error() == config.MONGO_NO_DOCUMENT {
			return errors.New("topUp_not_found")
		}
		return result.Err()
	}

	err := result.Decode(&topUp)
	if err != nil {
		return err
	}

	return nil

}

func (topUp *TopUpNoSql) Update(ctx context.Context,
	id string, bankCode string, bankName string, vaAccountNo string, vaAccountName string, topUpAmount float64) error {

	filter := bson.D{
		primitive.E{Key: "id", Value: id},
	}

	set := bson.D{
		primitive.E{
			Key: "$set",
			Value: bson.D{
				primitive.E{Key: "bankCode", Value: bankCode},
				primitive.E{Key: "bankName", Value: bankName},
				primitive.E{Key: "vaAccountNo", Value: vaAccountNo},
				primitive.E{Key: "vaAccountName", Value: vaAccountName},
				primitive.E{Key: "topUpAmount", Value: topUpAmount},
			},
		},
	}

	optionsAfter := options.After
	updateOptions := &options.FindOneAndUpdateOptions{
		ReturnDocument: &optionsAfter,
	}

	result := topUp.Collection().FindOneAndUpdate(ctx, filter, set, updateOptions)
	if result.Err() != nil {
		return result.Err()
	}

	return result.Decode(&topUp)

}

func (topUp TopUpNoSql) Delete(ctx context.Context, id string) error {

	filter := bson.D{
		primitive.E{Key: "id", Value: id},
	}

	result, err := topUp.Collection().DeleteOne(ctx, filter, nil)
	if err != nil {
		return err
	}
	log.Println("Deleted successful", result)

	return nil

}
