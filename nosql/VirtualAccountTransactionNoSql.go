package nosql

import (
	"context"
	"errors"
	"log"
	"virtual-account/config"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type VirtualAccountTransactionNoSql struct {
	BaseNoSql

	ID                string  `json:"id" bson:"id"`
	VaAccountNo       string  `json:"vaAccountNo" bson:"vaAccountNo"`
	VaAccountName     string  `json:"vaAccountName" bson:"vaAccountName"`
	TransactionAmount float64 `json:"transactionAmount" bson:"transactionAmount"`
	TransactionType   string  `json:"transactionType" bson:"transactionType"`
	Description       string  `json:"description" bson:"description"`
	TransactionID     string  `json:"TransactionID" bson:"TransactionID"`
}

const (
	TRANSACTION_TYPE_TOP_UP   = "TopUp"
	TRANSACTION_TYPE_PAY      = "Pay"
	TRANSACTION_TYPE_TRANSFER = "Transfer"
)

func CreateVirtualAccountTransactionNoSql(client *mongo.Client) VirtualAccountTransactionNoSql {
	virtualAccountTransaction := VirtualAccountTransactionNoSql{}
	virtualAccountTransaction.Client = client
	return virtualAccountTransaction
}

func (virtualAccountTransaction VirtualAccountTransactionNoSql) Collection() *mongo.Collection {
	return virtualAccountTransaction.Client.Database(config.DATABASE).Collection("virtual_account_transactions")
}

func (virtualAccountTransaction VirtualAccountTransactionNoSql) Truncate(ctx context.Context) error {
	return virtualAccountTransaction.Collection().Drop(ctx)
}

func (virtualAccountTransaction VirtualAccountTransactionNoSql) ListVirtualAccountTransaction(ctx context.Context) ([]VirtualAccountTransactionNoSql, error) {

	virtualAccountTransactions := []VirtualAccountTransactionNoSql{}
	filter := bson.D{}
	cursor, err := virtualAccountTransaction.Collection().Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &virtualAccountTransactions)
	if err != nil {
		return nil, err
	}

	return virtualAccountTransactions, nil

}

func (virtualAccountTransaction VirtualAccountTransactionNoSql) AddVirtualAccountTransaction(ctx context.Context,
	vaAccountNo string, vaAccountName string, transactionAmount float64, transactionType string, description string, bankTransactionID string) (*VirtualAccountTransactionNoSql, error) {

	virtualAccountTransaction.ID = uuid.New().String()
	virtualAccountTransaction.VaAccountNo = vaAccountNo
	virtualAccountTransaction.VaAccountName = vaAccountName
	virtualAccountTransaction.TransactionAmount = transactionAmount
	virtualAccountTransaction.TransactionType = transactionType
	virtualAccountTransaction.Description = description
	virtualAccountTransaction.TransactionID = bankTransactionID

	_, err := virtualAccountTransaction.Collection().InsertOne(ctx, virtualAccountTransaction)
	if err != nil {
		return nil, err
	}
	return &virtualAccountTransaction, nil

}

func (virtualAccountTransaction *VirtualAccountTransactionNoSql) FindOneByID(ctx context.Context, Id string) error {

	filter := bson.D{
		primitive.E{Key: "id", Value: Id},
	}

	result := virtualAccountTransaction.Collection().FindOne(ctx, filter)

	if result.Err() != nil {
		log.Println(result.Err())
		if result.Err().Error() == config.MONGO_NO_DOCUMENT {
			return errors.New("virtualAccountTransaction_not_found")
		}
		return result.Err()

	}

	err := result.Decode(&virtualAccountTransaction)
	if err != nil {
		return err
	}

	return nil

}

func (virtualAccountTransaction *VirtualAccountTransactionNoSql) FindOneByNo(ctx context.Context, vaAccountNo string) ([]VirtualAccountTransactionNoSql, error) {

	result, err := virtualAccountTransaction.Collection().Find(ctx, bson.D{primitive.E{Key: "vaAccountNo", Value: vaAccountNo}})
	if err != nil {
		return nil, err
	}
	var filtered []VirtualAccountTransactionNoSql
	if err = result.All(ctx, &filtered); err != nil {
		return nil, err
	}

	return filtered, nil

}

func (virtualAccountTransaction VirtualAccountTransactionNoSql) Delete(ctx context.Context, id string) error {
	filter := bson.D{
		primitive.E{Key: "id", Value: id},
	}
	result, err := virtualAccountTransaction.Collection().DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	log.Println("Deleted successful", result)

	return nil
}
