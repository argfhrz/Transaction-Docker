package apis

import (
	"bank-account/connection"
	"bank-account/nosql"
	"context"
	"encoding/json"
	"net/http"
)

type BankTransactionApi struct {
	BaseApi

	nosql.BankTransactionNoSql
}

func (bankTransactionApi BankTransactionApi) GetBankTransaction(w http.ResponseWriter, r *http.Request) {

	ctx := context.TODO()
	mongoClient, err := connection.OpenMongoDb(ctx)
	if err != nil {
		bankTransactionApi.Error(w, err)

	}
	defer connection.CloseMongoDb(ctx, mongoClient)

	bankTransactionNoSql := nosql.CreateBankTransactionNoSql(mongoClient)
	bankTransactions, err := bankTransactionNoSql.ListBankTransaction(ctx)
	if err != nil {
		bankTransactionApi.Error(w, err)
		return
	} else {
		bankTransactionApi.Json(w, bankTransactions, http.StatusOK)
		return
	}

}

func (bankTransactionApi BankTransactionApi) GetBankTransactionByID(w http.ResponseWriter, r *http.Request) {

	ctx := context.TODO()
	mongoClient, err := connection.OpenMongoDb(ctx)
	if err != nil {
		bankTransactionApi.Error(w, err)

	}
	defer connection.CloseMongoDb(ctx, mongoClient)

	bankTransactionNo := bankTransactionApi.QueryParam(r, "bankTransactionNo")

	bankTransactionNoSql := nosql.CreateBankTransactionNoSql(mongoClient)
	err = bankTransactionNoSql.FindOneByID(ctx, bankTransactionNo)
	if err != nil {
		bankTransactionApi.Error(w, err)
		return
	} else {
		bankTransactionApi.Json(w, bankTransactionNoSql, http.StatusOK)
		return
	}

}

func (bankTransactionApi BankTransactionApi) GetBankTransactionByNo(w http.ResponseWriter, r *http.Request) {

	ctx := context.TODO()
	mongoClient, err := connection.OpenMongoDb(ctx)
	if err != nil {
		bankTransactionApi.Error(w, err)

	}
	defer connection.CloseMongoDb(ctx, mongoClient)

	bankAccountNo := bankTransactionApi.QueryParam(r, "bankAccountNo")

	bankTransactionNoSql := nosql.CreateBankTransactionNoSql(mongoClient)
	bankTransaction, err := bankTransactionNoSql.FindOneByNo(ctx, bankAccountNo)
	if err != nil {
		bankTransactionApi.Error(w, err)
		return
	} else {
		bankTransactionApi.Json(w, bankTransaction, http.StatusOK)
		return
	}

}

func (bankTransactionApi BankTransactionApi) RemoveBankTransaction(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&bankTransactionApi)
	if err != nil {
		bankTransactionApi.Error(w, err)
		return
	}

	ctx := context.TODO()
	mongoClient, err := connection.OpenMongoDb(ctx)
	if err != nil {
		bankTransactionApi.Error(w, err)

	}
	defer connection.CloseMongoDb(ctx, mongoClient)

	bankTransactionNoSql := nosql.CreateBankTransactionNoSql(mongoClient)

	err = bankTransactionNoSql.Delete(ctx, bankTransactionApi.ID)
	if err != nil {
		bankTransactionApi.Error(w, err)
		return
	} else {
		bankTransactionApi.Empty(w, http.StatusOK)
		return
	}

}
