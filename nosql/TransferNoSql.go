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

type TransferNoSql struct {
	BaseNoSql `bson:"-"`

	ID                string    `json:"id" bson:"id"`
	ScrVaAccountNo    string    `json:"srcVaAccountNo" bson:"srcVaAccountNo"`
	SrcVaAccountName  string    `json:"srcVaAccountName" bson:"srcVaAccountName"`
	DestVaAccountNo   string    `json:"destVaAccountNo" bson:"destVaAccountNo"`
	DestVaAccountName string    `json:"destVaAccountName" bson:"destVaAccountName"`
	TransferAmount    float64   `json:"transferAmount" bson:"transferAmount"`
	CreatedAt         time.Time `json:"createdAt" bson:"createdAt"`
}

func CreateTransferNoSql(client *mongo.Client) TransferNoSql {
	transfer := TransferNoSql{}
	transfer.Client = client
	return transfer
}

func (transfer TransferNoSql) Collection() *mongo.Collection {
	return transfer.Client.Database(config.DATABASE).Collection("transfers")
}

func (transfer TransferNoSql) Truncate(ctx context.Context) error {
	return transfer.Collection().Drop(ctx)
}

func (transfer TransferNoSql) ListTransfer(ctx context.Context) ([]TransferNoSql, error) {

	transfers := []TransferNoSql{}
	filter := bson.D{}
	cursor, err := transfer.Collection().Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &transfers)
	if err != nil {
		return nil, err
	}

	return transfers, nil

}

func (transfer TransferNoSql) AddTransfer(ctx context.Context, srcVaAccountNo string, srcVaAccountName string, destVaAccountNo string, destVaAccountName string, transferAmount float64) (*TransferNoSql, error) {

	transfer.ID = uuid.New().String()
	transfer.ScrVaAccountNo = srcVaAccountNo
	transfer.SrcVaAccountName = srcVaAccountName
	transfer.DestVaAccountNo = destVaAccountNo
	transfer.DestVaAccountName = destVaAccountName
	transfer.TransferAmount = transferAmount

	_, err := transfer.Collection().InsertOne(ctx, transfer)
	if err != nil {
		return nil, err
	}
	return &transfer, nil
}

func (transfer *TransferNoSql) FindOneByID(ctx context.Context, Id string) error {

	filter := bson.D{
		primitive.E{Key: "id", Value: Id},
	}

	result := transfer.Collection().FindOne(ctx, filter)

	if result.Err() != nil {
		log.Println(result.Err())
		if result.Err().Error() == config.MONGO_NO_DOCUMENT {
			return errors.New("transfer_not_found")
		}
		return result.Err()

	}

	err := result.Decode(&transfer)
	if err != nil {
		return err
	}

	return nil

}

func (transfer TransferNoSql) Delete(ctx context.Context, id string) error {
	filter := bson.D{
		primitive.E{Key: "id", Value: id},
	}
	result, err := transfer.Collection().DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	log.Println("Deleted successful", result)

	return nil

}
