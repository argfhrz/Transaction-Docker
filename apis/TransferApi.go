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

type TransferApi struct {
	BaseApi

	nosql.TransferNoSql
	Pin string `json:"pin" bson:"pin"`
}

func (transferApi TransferApi) Connection(w http.ResponseWriter, r *http.Request) (*mongo.Client, context.Context, *sql.DB) {

	ctx := context.TODO()
	mongoClient, err := connection.OpenMongoDb(ctx)
	if err != nil {
		transferApi.Error(w, err)
		return nil, nil, nil
	}

	db := connection.OpenConnection(config.DEV)

	_, err = transferApi.ParseToken(r, ctx, mongoClient, db)
	if err != nil {
		transferApi.Error(w, err)
		return nil, nil, nil
	}

	return mongoClient, ctx, db

}

func (transferApi TransferApi) PostTransfer(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&transferApi)
	if err != nil {
		transferApi.Error(w, err)
		return
	}

	mongoClient, ctx, db := transferApi.Connection(w, r)
	defer connection.CloseMongoDb(ctx, mongoClient)
	defer db.Close()

	transferService := services.CreateTransferService(mongoClient, db)

	transferNo, err := transferService.Transfer(ctx, transferApi.ScrVaAccountNo, transferApi.DestVaAccountNo, transferApi.TransferAmount, transferApi.Pin)
	if err != nil {
		transferApi.Error(w, err)
		return
	} else {

		transferApi.Json(w, transferNo, http.StatusOK)
		log.Println(transferNo)

	}

}

func (transferApi TransferApi) RemoveTransfer(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&transferApi)
	if err != nil {
		transferApi.Error(w, err)
		return
	}

	mongoClient, ctx, db := transferApi.Connection(w, r)
	defer connection.CloseMongoDb(ctx, mongoClient)
	defer db.Close()

	transferNoSql := nosql.CreateTransferNoSql(mongoClient)

	err = transferNoSql.Delete(ctx, transferApi.ID)
	if err != nil {
		transferApi.Error(w, err)
		return
	} else {
		transferApi.Empty(w, http.StatusOK)
		return
	}

}
