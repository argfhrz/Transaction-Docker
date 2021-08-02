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

func TestLoginService(t *testing.T) {

	ctx := context.TODO()
	mongoClient, err := connection.OpenMongoDb(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer connection.CloseMongoDb(ctx, mongoClient)

	db := connection.OpenConnection(config.DEV)
	defer db.Close()

	loginModel := nosql.CreateLoginNoSql(mongoClient)
	err = loginModel.Truncate(ctx)
	if err != nil {
		t.Fatal(err)
	}

	phoneNumber := "0812312312"
	password := "123"

	loginService := CreateLoginService(mongoClient, db)
	token, err := loginService.Login(ctx, phoneNumber, password)
	if err != nil {
		t.Fatal(err)
	}

	log.Println("Token=", token)

	login, err := loginService.ParseToken(ctx, token)
	if err != nil {
		t.Fatal(err)
	}

	log.Println("login=", helpers.ToJson(login))

}
