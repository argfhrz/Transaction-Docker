package services

import (
	"context"
	"log"
	"testing"
	"virtual-account/config"
	"virtual-account/connection"
	"virtual-account/helpers"
	"virtual-account/nosql"
)

func TestTransferService(t *testing.T) {
	ctx := context.TODO()
	mongoClient, err := connection.OpenMongoDb(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer connection.CloseMongoDb(ctx, mongoClient)

	db := connection.OpenConnection(config.DEV)
	defer db.Close()

	transferNoSql := nosql.CreateTransferNoSql(mongoClient)
	err = transferNoSql.Truncate(ctx)
	if err != nil {
		t.Fatal(err)
	}

	srcVaAccountNo := "0812312312-00001"
	destVaAccountNo := "0812312313-00002"
	pin := "123456"
	transferAmount := float64(2000000)

	transferService := CreateTransferService(mongoClient, db)
	transfer, err := transferService.Transfer(ctx, srcVaAccountNo, destVaAccountNo, transferAmount, pin)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(helpers.ToJson(transfer))
}
