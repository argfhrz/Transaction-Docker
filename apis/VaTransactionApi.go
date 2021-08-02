package apis

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"virtual-account/config"
	"virtual-account/connection"
	"virtual-account/nosql"

	"go.mongodb.org/mongo-driver/mongo"
)

type VaTransactionApi struct {
	BaseApi

	nosql.VirtualAccountTransactionNoSql
}

func (vaTransactionApi VaTransactionApi) Connection(w http.ResponseWriter, r *http.Request) (*mongo.Client, context.Context, *sql.DB) {

	ctx := context.TODO()
	mongoClient, err := connection.OpenMongoDb(ctx)
	if err != nil {
		vaTransactionApi.Error(w, err)
		return nil, nil, nil
	}

	db := connection.OpenConnection(config.DEV)

	_, err = vaTransactionApi.ParseToken(r, ctx, mongoClient, db)
	if err != nil {
		vaTransactionApi.Error(w, err)
		return nil, nil, nil
	}

	return mongoClient, ctx, db

}

func (vaTransactionApi VaTransactionApi) GetVaTransaction(w http.ResponseWriter, r *http.Request) {

	mongoClient, ctx, db := vaTransactionApi.Connection(w, r)
	defer connection.CloseMongoDb(ctx, mongoClient)
	defer db.Close()

	vaTransactionNoSql := nosql.CreateVirtualAccountTransactionNoSql(mongoClient)
	vaTransactions, err := vaTransactionNoSql.ListVirtualAccountTransaction(ctx)
	if err != nil {
		vaTransactionApi.Error(w, err)
		return
	} else {
		vaTransactionApi.Json(w, vaTransactions, http.StatusOK)
		return
	}

}

func (vaTransactionApi VaTransactionApi) GetVaTransactionByID(w http.ResponseWriter, r *http.Request) {

	mongoClient, ctx, db := vaTransactionApi.Connection(w, r)
	defer connection.CloseMongoDb(ctx, mongoClient)
	defer db.Close()

	vaTransactionNo := vaTransactionApi.QueryParam(r, "vaTransactionNo")

	vaTransactionNoSql := nosql.CreateVirtualAccountTransactionNoSql(mongoClient)
	err := vaTransactionNoSql.FindOneByID(ctx, vaTransactionNo)
	if err != nil {
		vaTransactionApi.Error(w, err)
		return
	} else {
		vaTransactionApi.Json(w, vaTransactionNoSql, http.StatusOK)
		return
	}

}

func (vaTransactionApi VaTransactionApi) GetVaTransactionByNo(w http.ResponseWriter, r *http.Request) {

	mongoClient, ctx, db := vaTransactionApi.Connection(w, r)
	defer connection.CloseMongoDb(ctx, mongoClient)
	defer db.Close()

	virtualAccountNo := vaTransactionApi.QueryParam(r, "virtualAccountNo")

	vaTransactionNoSql := nosql.CreateVirtualAccountTransactionNoSql(mongoClient)
	vaTransaction, err := vaTransactionNoSql.FindOneByNo(ctx, virtualAccountNo)
	if err != nil {
		vaTransactionApi.Error(w, err)
		return
	} else {
		vaTransactionApi.Json(w, vaTransaction, http.StatusOK)
		return
	}

}

func (vaTransactionApi VaTransactionApi) RemoveVirtualTransaction(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&vaTransactionApi)
	if err != nil {
		vaTransactionApi.Error(w, err)
		return
	}

	mongoClient, ctx, db := vaTransactionApi.Connection(w, r)
	defer connection.CloseMongoDb(ctx, mongoClient)
	defer db.Close()

	vaTransactionNoSql := nosql.CreateVirtualAccountTransactionNoSql(mongoClient)

	err = vaTransactionNoSql.Delete(ctx, vaTransactionApi.ID)
	if err != nil {
		vaTransactionApi.Error(w, err)
		return
	} else {
		vaTransactionApi.Empty(w, http.StatusOK)
		return
	}

}
