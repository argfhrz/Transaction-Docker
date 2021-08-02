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

func TestPayService(t *testing.T) {
	ctx := context.TODO()
	mongoClient, err := connection.OpenMongoDb(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer connection.CloseMongoDb(ctx, mongoClient)

	db := connection.OpenConnection(config.DEV)
	defer db.Close()

	payNoSql := nosql.CreatePayNoSql(mongoClient)
	err = payNoSql.Truncate(ctx)
	if err != nil {
		t.Fatal(err)
	}

	srcVaAccountNo := "0812312312-00001"
	mrchVaAccountNo := "0812234567-00003"
	pin := "123456"
	payAmount := float64(1000000)

	payService := CreatePayService(mongoClient, db)
	pay, err := payService.Pay(ctx, mrchVaAccountNo, srcVaAccountNo, payAmount, pin)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(helpers.ToJson(pay))
}
