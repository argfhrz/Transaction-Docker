package apis

import (
	"context"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"virtual-account/config"
	"virtual-account/connection"
	"virtual-account/data"
	"virtual-account/helpers"
	"virtual-account/nosql"
	"virtual-account/services"

	"go.mongodb.org/mongo-driver/mongo"
)

const (
	SERVER_BANK = "http://172.25.0.3:8300"
)

type TopUpApi struct {
	BaseApi

	nosql.TopUpNoSql
	Pin string `json:"pin" bson:"pin"`
}

func (topUpApi TopUpApi) Connection(w http.ResponseWriter, r *http.Request) (*mongo.Client, context.Context, *sql.DB) {

	ctx := context.TODO()
	mongoClient, err := connection.OpenMongoDb(ctx)
	if err != nil {
		topUpApi.Error(w, err)
		return nil, nil, nil
	}

	db := connection.OpenConnection(config.DEV)

	_, err = topUpApi.ParseToken(r, ctx, mongoClient, db)
	if err != nil {
		topUpApi.Error(w, err)
		return nil, nil, nil
	}

	return mongoClient, ctx, db

}

func (topUpApi TopUpApi) PostTopUp(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&topUpApi)
	if err != nil {
		topUpApi.Error(w, err)
		return
	}

	mongoClient, ctx, db := topUpApi.Connection(w, r)
	defer connection.CloseMongoDb(ctx, mongoClient)
	defer db.Close()

	response, err := http.Get(SERVER_BANK + "/bank-account/id?bankAccountNo=" + url.QueryEscape(topUpApi.BankCode))
	if err != nil {
		topUpApi.Error(w, err)
		return
	}

	defer response.Body.Close()

	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		topUpApi.Error(w, err)
		return
	}

	bankAccounts := data.BankAccount{}
	err = json.Unmarshal(respBody, &bankAccounts)
	if err != nil {
		topUpApi.Error(w, err)
		return
	}

	if bankAccounts.Saldo < topUpApi.TopUpAmount {
		topUpApi.Error(w, err)
		return
	} else {

		topUpService := services.CreateTopUpService(mongoClient, db)
		topUp, bankTransaction, virtualAccountTransaction, err := topUpService.TopUp(ctx, topUpApi.BankCode, topUpApi.VaAccountNo, topUpApi.Pin, topUpApi.TopUpAmount)
		if err != nil {
			topUpApi.Error(w, err)
			return
		} else {
			topUpApi.Json(w, topUp, http.StatusOK)
			log.Println(helpers.ToJson(bankTransaction))
			log.Println(helpers.ToJson(virtualAccountTransaction))
			return

		}
	}

}

func (topUpApi TopUpApi) RemoveTopUp(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&topUpApi)
	if err != nil {
		topUpApi.Error(w, err)
		return
	}

	mongoClient, ctx, db := topUpApi.Connection(w, r)
	defer connection.CloseMongoDb(ctx, mongoClient)
	defer db.Close()

	topUpNoSql := nosql.CreateTopUpNoSql(mongoClient)

	err = topUpNoSql.Delete(ctx, topUpApi.ID)
	if err != nil {
		topUpApi.Error(w, err)
		return
	} else {
		topUpApi.Empty(w, http.StatusOK)
		return
	}

}
