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

type BankTransactionNoSql struct {
	BaseNoSql `bson:"-"`

	ID                string  `json:"id" bson:"id"`
	BankAccountNo     string  `json:"bankAccountNo" bson:"bankAccountNo"`
	BankAccountOwner  string  `json:"bankAccountOwner" bson:"bankAccountOwner"`
	TransactionAmount float64 `json:"transactionAmount" bson:"transactionAmount"`
	Reference         string  `json:"reference" bson:"reference"`
}

func CreateBankTransactionNoSql(client *mongo.Client) BankTransactionNoSql {
	transaction := BankTransactionNoSql{}
	transaction.Client = client
	return transaction
}

func (transaction BankTransactionNoSql) Collection() *mongo.Collection {
	return transaction.Client.Database(config.DATABASE).Collection("bank_transactions")
}

func (transaction BankTransactionNoSql) Truncate(ctx context.Context) error {
	return transaction.Collection().Drop(ctx)
}

func (transaction BankTransactionNoSql) ListBankTransaction(ctx context.Context) ([]BankTransactionNoSql, error) {

	transactions := []BankTransactionNoSql{}
	filter := bson.D{}
	cursor, err := transaction.Collection().Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &transactions)
	if err != nil {
		return nil, err
	}

	return transactions, nil

}

func (transaction BankTransactionNoSql) AddBankTransaction(ctx context.Context, bankAccountNo string, bankAccountOwner string, transactionAmount float64, reference string) (*BankTransactionNoSql, error) {

	transaction.ID = uuid.New().String()
	transaction.BankAccountNo = bankAccountNo
	transaction.BankAccountOwner = bankAccountOwner
	transaction.TransactionAmount = transactionAmount
	transaction.Reference = reference

	_, err := transaction.Collection().InsertOne(ctx, transaction)
	if err != nil {
		return nil, err
	}
	return &transaction, nil

}

func (transaction *BankTransactionNoSql) FindOneByID(ctx context.Context, Id string) error {

	filter := bson.D{
		primitive.E{Key: "id", Value: Id},
	}

	result := transaction.Collection().FindOne(ctx, filter)

	if result.Err() != nil {
		log.Println(result.Err())
		if result.Err().Error() == config.MONGO_NO_DOCUMENT {
			return errors.New("bank_transaction_not_found")
		}
		return result.Err()

	}

	err := result.Decode(&transaction)
	if err != nil {
		return err
	}

	return nil

}

func (transaction BankTransactionNoSql) Delete(ctx context.Context, id string) ([]BankTransactionNoSql, error) {
	filter := bson.D{
		primitive.E{Key: "id", Value: id},
	}
	result, err := transaction.Collection().DeleteOne(ctx, filter)
	if err != nil {
		return nil, err
	}
	log.Println("Deleted successful", result)
	newResult, err := transaction.ListBankTransaction(ctx)
	if err != nil {
		return nil, err
	}

	return newResult, nil
}
