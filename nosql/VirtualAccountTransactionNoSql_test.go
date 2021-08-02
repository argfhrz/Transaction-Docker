package nosql

import (
	"context"
	"log"
	"testing"
	"virtual-account/connection"
	"virtual-account/helpers"
)

func TestVirtualAccountTransactionNoSql(t *testing.T) {
	ctx := context.TODO()

	conn, err := connection.OpenMongoDb(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer connection.CloseMongoDb(ctx, conn)

	virtualAccountTransactionNoSql := CreateVirtualAccountTransactionNoSql(conn)
	err = virtualAccountTransactionNoSql.Truncate(ctx)
	if err != nil {
		t.Fatal(err)
	}

	vaAccountNo := "12321312311231"
	vaAccountName := "owner1"
	transactionAmount := float64(200000)
	transactionType := "TopUp"
	description := "testing"
	bankTransactionID := "123123131"

	virtualAccountTransaction, err := virtualAccountTransactionNoSql.AddVirtualAccountTransaction(ctx, vaAccountNo, vaAccountName, transactionAmount, transactionType, description, bankTransactionID)
	if err != nil {
		t.Fatal(err)
	}

	log.Println(virtualAccountTransaction)

	virtualAccountTransactions, err := virtualAccountTransactionNoSql.ListVirtualAccountTransaction(ctx)
	if err != nil {
		t.Fatal(err)
	}

	log.Println(helpers.ToJson(virtualAccountTransactions))

	err = virtualAccountTransactionNoSql.FindOneByID(ctx, virtualAccountTransaction.ID)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(helpers.ToJson(virtualAccountTransactionNoSql))

}
