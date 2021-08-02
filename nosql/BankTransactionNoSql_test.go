package nosql

import (
	"context"
	"log"
	"testing"
	"virtual-account/connection"
	"virtual-account/helpers"
)

func TestBankTransactionNoSql(t *testing.T) {
	ctx := context.TODO()

	conn, err := connection.OpenMongoDb(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer connection.CloseMongoDb(ctx, conn)

	bankTransactionNoSql := CreateBankTransactionNoSql(conn)
	err = bankTransactionNoSql.Truncate(ctx)
	if err != nil {
		t.Fatal(err)
	}

	bankAccountNo := "12321312311231"
	bankAccountOwner := "owner1"
	transactionAmount := float64(200000)
	reference := ""

	bankTransaction, err := bankTransactionNoSql.AddBankTransaction(ctx, bankAccountNo, bankAccountOwner, transactionAmount, reference)
	if err != nil {
		t.Fatal(err)
	}

	log.Println(bankTransaction)

	bankTransactions, err := bankTransactionNoSql.ListBankTransaction(ctx)
	if err != nil {
		t.Fatal(err)
	}

	log.Println(helpers.ToJson(bankTransactions))

	err = bankTransactionNoSql.FindOneByID(ctx, bankTransaction.ID)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(helpers.ToJson(bankTransactionNoSql))

}
