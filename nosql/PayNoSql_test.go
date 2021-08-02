package nosql

import (
	"context"
	"log"
	"testing"
	"virtual-account/connection"
	"virtual-account/helpers"
)

func TestPayNoSql(t *testing.T) {
	ctx := context.TODO()

	conn, err := connection.OpenMongoDb(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer connection.CloseMongoDb(ctx, conn)

	payNoSql := CreatePayNoSql(conn)
	err = payNoSql.Truncate(ctx)
	if err != nil {
		t.Fatal(err)
	}

	merchantVaAccountNo := "0812234567-00003"
	merchantVaAccountName := "merchant1"
	vaAccountNo := "0812312312-00001"
	vaAccountName := "account1"
	payAmount := float64(100000)

	pay, err := payNoSql.AddPay(ctx, merchantVaAccountNo, merchantVaAccountName, vaAccountNo, vaAccountName, payAmount)
	if err != nil {
		t.Fatal(err)
	}

	log.Println(pay)

	pays, err := payNoSql.ListPay(ctx)
	if err != nil {
		t.Fatal(err)
	}

	log.Println(helpers.ToJson(pays))

	err = payNoSql.FindOneByID(ctx, pay.ID)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(helpers.ToJson(payNoSql))

}
