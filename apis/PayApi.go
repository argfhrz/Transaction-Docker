package apis

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"virtual-account/config"
	"virtual-account/connection"
	"virtual-account/nosql"
	"virtual-account/services"

	"go.mongodb.org/mongo-driver/mongo"
)

type PayApi struct {
	BaseApi

	nosql.PayNoSql
	Pin string `json:"pin" bson:"pin"`
}

func (payApi PayApi) Connection(w http.ResponseWriter, r *http.Request) (*mongo.Client, context.Context, *sql.DB) {

	ctx := context.TODO()
	mongoClient, err := connection.OpenMongoDb(ctx)
	if err != nil {
		payApi.Error(w, err)
		return nil, nil, nil
	}

	db := connection.OpenConnection(config.DEV)

	_, err = payApi.ParseToken(r, ctx, mongoClient, db)
	if err != nil {
		payApi.Error(w, err)
		return nil, nil, nil
	}

	return mongoClient, ctx, db

}

func (payApi PayApi) PostPay(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&payApi)
	if err != nil {
		payApi.Error(w, err)
		return
	}

	mongoClient, ctx, db := payApi.Connection(w, r)
	defer connection.CloseMongoDb(ctx, mongoClient)
	defer db.Close()

	payService := services.CreatePayService(mongoClient, db)

	payNo, err := payService.Pay(ctx, payApi.MerchantVaAccountNo, payApi.SrcVaAccountNo, payApi.PayAmount, payApi.Pin)
	if err != nil {
		payApi.Error(w, err)
		return
	} else {

		payApi.Json(w, payNo, http.StatusOK)
		log.Println(payNo)

	}

}

func (payApi PayApi) RemovePay(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&payApi)
	if err != nil {
		payApi.Error(w, err)
		return
	}

	mongoClient, ctx, db := payApi.Connection(w, r)
	defer connection.CloseMongoDb(ctx, mongoClient)
	defer db.Close()

	payNoSql := nosql.CreatePayNoSql(mongoClient)

	err = payNoSql.Delete(ctx, payApi.ID)
	if err != nil {
		payApi.Error(w, err)
		return
	} else {
		payApi.Empty(w, http.StatusOK)
		return
	}

}
