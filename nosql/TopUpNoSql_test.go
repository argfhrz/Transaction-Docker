package nosql

import (
	"context"
	"log"
	"testing"
	"virtual-account/connection"
	"virtual-account/helpers"
)

func TestTopUpNoSql(t *testing.T) {
	ctx := context.TODO()

	conn, err := connection.OpenMongoDb(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer connection.CloseMongoDb(ctx, conn)

	topUpNoSql := CreateTopUpNoSql(conn)
	err = topUpNoSql.Truncate(ctx)
	if err != nil {
		t.Fatal(err)
	}

	bankCode := "12321312311231"
	bankName := "owner1"
	vaAccountNo := "0812312312-00001"
	vaAccountName := "account1"
	topUpAmount := float64(200000)

	topUp, err := topUpNoSql.AddTopUp(ctx, bankCode, bankName, vaAccountNo, vaAccountName, topUpAmount)
	if err != nil {
		t.Fatal(err)
	}

	log.Println(topUp)

	topUps, err := topUpNoSql.ListTopUp(ctx)
	if err != nil {
		t.Fatal(err)
	}

	log.Println(helpers.ToJson(topUps))

	err = topUpNoSql.FindOneByID(ctx, topUp.ID)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(helpers.ToJson(topUpNoSql))

}
