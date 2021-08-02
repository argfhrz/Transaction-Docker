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

func TestTopUpService(t *testing.T) {
	ctx := context.TODO()
	mongoClient, err := connection.OpenMongoDb(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer connection.CloseMongoDb(ctx, mongoClient)

	db := connection.OpenConnection(config.DEV)
	defer db.Close()

	topUpNoSql := nosql.CreateTopUpNoSql(mongoClient)
	err = topUpNoSql.Truncate(ctx)
	if err != nil {
		t.Fatal(err)
	}

	bankAccountNo := "12321312311231"
	vaAccountNo := "0812312312-00001"
	pin := "123456"
	transactionAmount := float64(2000000)

	topUpService := CreateTopUpService(mongoClient, db)
	topUp, bankTransaction, virtualAccounTransaction, err := topUpService.TopUp(ctx, bankAccountNo, vaAccountNo, pin, transactionAmount)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(helpers.ToJson(topUp))
	log.Println(helpers.ToJson(bankTransaction))
	log.Println(helpers.ToJson(virtualAccounTransaction))

}
