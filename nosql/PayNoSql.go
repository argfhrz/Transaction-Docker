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
)

type PayNoSql struct {
	BaseNoSql

	ID                    string    `json:"id" bson:"id"`
	MerchantVaAccountNo   string    `json:"merchantVaAccountNo" bson:"merchantVaAccountNo"`
	MerchantVaAccountName string    `json:"merchantVaAccountName" bson:"merchantVaAccountName"`
	SrcVaAccountNo        string    `json:"srcVaAccountNo" bson:"srcVaAccountNo"`
	SrcVaAccountName      string    `json:"srcVaAccountName" bson:"srcVaAccountName"`
	PayAmount             float64   `json:"payAmount" bson:"payAmount"`
	CreatedAt             time.Time `json:"createdAt" bson:"createdAt"`
}

func CreatePayNoSql(client *mongo.Client) PayNoSql {
	pay := PayNoSql{}
	pay.Client = client
	return pay
}

func (pay PayNoSql) Collection() *mongo.Collection {
	return pay.Client.Database(config.DATABASE).Collection("pays")
}

func (pay PayNoSql) Truncate(ctx context.Context) error {
	return pay.Collection().Drop(ctx)
}

func (pay PayNoSql) ListPay(ctx context.Context) ([]PayNoSql, error) {

	pays := []PayNoSql{}
	filter := bson.D{}
	cursor, err := pay.Collection().Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &pays)
	if err != nil {
		return nil, err
	}

	return pays, nil

}

func (pay PayNoSql) AddPay(ctx context.Context,
	merchantVaAccountNo string, merchantVaAccountName string, srcVaAccountNo string, srcVaAccountName string, payAmount float64) (*PayNoSql, error) {

	pay.ID = uuid.New().String()
	pay.MerchantVaAccountNo = merchantVaAccountNo
	pay.MerchantVaAccountName = merchantVaAccountName
	pay.SrcVaAccountNo = srcVaAccountNo
	pay.SrcVaAccountName = srcVaAccountName
	pay.PayAmount = payAmount
	pay.CreatedAt = time.Now().UTC()

	_, err := pay.Collection().InsertOne(ctx, pay)
	if err != nil {
		return nil, err
	}
	return &pay, nil

}

func (pay *PayNoSql) FindOneByID(ctx context.Context, Id string) error {

	filter := bson.D{
		primitive.E{Key: "id", Value: Id},
	}

	result := pay.Collection().FindOne(ctx, filter)

	if result.Err() != nil {
		log.Println(result.Err())
		if result.Err().Error() == config.MONGO_NO_DOCUMENT {
			return errors.New("pay_not_found")
		}
		return result.Err()

	}

	err := result.Decode(&pay)
	if err != nil {
		return err
	}

	return nil

}

func (pay PayNoSql) Delete(ctx context.Context, id string) error {
	filter := bson.D{
		primitive.E{Key: "id", Value: id},
	}
	result, err := pay.Collection().DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	log.Println("Deleted successful", result)

	return nil
}
