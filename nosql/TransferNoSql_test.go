package nosql

import (
	"context"
	"log"
	"testing"
	"virtual-account/connection"
	"virtual-account/helpers"
)

func TestTransferNoSql(t *testing.T) {
	ctx := context.TODO()

	conn, err := connection.OpenMongoDb(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer connection.CloseMongoDb(ctx, conn)

	transferNoSql := CreateTransferNoSql(conn)
	err = transferNoSql.Truncate(ctx)
	if err != nil {
		t.Fatal(err)
	}

	srcVaAccountNo := "0812312312-00001"
	srcVaAccountName := "account1"
	destVaAccountNo := "0812312312-00002"
	destVaAccountName := "account2"
	transferAmount := float64(200000)

	transfer, err := transferNoSql.AddTransfer(ctx, srcVaAccountNo, srcVaAccountName, destVaAccountNo, destVaAccountName, transferAmount)
	if err != nil {
		t.Fatal(err)
	}

	log.Println(transfer)

	transfers, err := transferNoSql.ListTransfer(ctx)
	if err != nil {
		t.Fatal(err)
	}

	log.Println(helpers.ToJson(transfers))

	err = transferNoSql.FindOneByID(ctx, transfer.ID)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(helpers.ToJson(transferNoSql))

}
