package apis

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"virtual-account/config"
	"virtual-account/connection"
	"virtual-account/data"

	"go.mongodb.org/mongo-driver/mongo"
)

type VirtualAccountApi struct {
	BaseApi

	data.VirtualAccount
}

func (virtualAccountApi VirtualAccountApi) Connection(w http.ResponseWriter, r *http.Request) (*mongo.Client, context.Context, *sql.DB) {

	ctx := context.TODO()
	mongoClient, err := connection.OpenMongoDb(ctx)
	if err != nil {
		virtualAccountApi.Error(w, err)
		return nil, nil, nil
	}

	db := connection.OpenConnection(config.DEV)

	_, err = virtualAccountApi.ParseToken(r, ctx, mongoClient, db)
	if err != nil {
		virtualAccountApi.Error(w, err)
		return nil, nil, nil
	}

	return mongoClient, ctx, db

}

func (virtualAccountApi VirtualAccountApi) GetVirtualAccount(w http.ResponseWriter, r *http.Request) {

	mongoClient, ctx, db := virtualAccountApi.Connection(w, r)
	defer connection.CloseMongoDb(ctx, mongoClient)
	defer db.Close()

	virtualAccountModel := data.CreateVirtualAccount(db)
	virtualAccounts, err := virtualAccountModel.GetListVirtualAccount()
	if err != nil {
		virtualAccountApi.Error(w, err)
		return
	} else {
		virtualAccountApi.Json(w, virtualAccounts, http.StatusOK)
		return
	}

}

func (virtualAccountApi VirtualAccountApi) GetVirtualAccountByID(w http.ResponseWriter, r *http.Request) {

	mongoClient, ctx, db := virtualAccountApi.Connection(w, r)
	defer connection.CloseMongoDb(ctx, mongoClient)
	defer db.Close()

	virtualAccountNo := virtualAccountApi.QueryParam(r, "virtualAccountNo")

	virtualAccountModel := data.CreateVirtualAccount(db)
	virtualAccount, err := virtualAccountModel.FindVirtualAccountByNo(virtualAccountNo)
	if err != nil {
		virtualAccountApi.Error(w, err)
		return
	} else {
		virtualAccountApi.Json(w, virtualAccount, http.StatusOK)
		return
	}

}

func (virtualAccountApi VirtualAccountApi) PostVirtualAccount(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&virtualAccountApi)
	if err != nil {
		virtualAccountApi.Error(w, err)
		return
	}

	db := connection.OpenConnection(config.DEV)
	defer db.Close()

	virtualAccountModel := data.CreateVirtualAccount(db)

	virtualAccountNo, err := virtualAccountModel.Add(virtualAccountApi.PhoneNumber, virtualAccountApi.Email, virtualAccountApi.AccountName, virtualAccountApi.Pin, virtualAccountApi.Password)
	if err != nil {
		virtualAccountApi.Error(w, err)
		return
	} else {
		vAccount, err := virtualAccountModel.FindVirtualAccountByNo(virtualAccountNo)
		if err != nil {
			virtualAccountApi.Error(w, err)
			return
		} else {
			virtualAccountApi.Json(w, vAccount, http.StatusOK)
			log.Println(virtualAccountNo)
			return
		}
	}

}

func (virtualAccountApi VirtualAccountApi) UpdateVirtualAccountName(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&virtualAccountApi)
	if err != nil {
		virtualAccountApi.Error(w, err)
		return
	}

	mongoClient, ctx, db := virtualAccountApi.Connection(w, r)
	defer connection.CloseMongoDb(ctx, mongoClient)
	defer db.Close()

	virtualAccountModel := data.CreateVirtualAccount(db)

	err = virtualAccountModel.Update(virtualAccountApi.VirtualAccountNo, virtualAccountApi.AccountName)
	if err != nil {
		virtualAccountApi.Error(w, err)
		return
	} else {
		vAccount, err := virtualAccountModel.FindVirtualAccountByNo(virtualAccountApi.VirtualAccountNo)
		if err != nil {
			virtualAccountApi.Error(w, err)
			return
		} else {
			virtualAccountApi.Json(w, vAccount, http.StatusOK)
			return
		}
	}
}

func (virtualAccountApi VirtualAccountApi) UpdateVirtualAccountPin(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&virtualAccountApi)
	if err != nil {
		virtualAccountApi.Error(w, err)
		return
	}

	mongoClient, ctx, db := virtualAccountApi.Connection(w, r)
	defer connection.CloseMongoDb(ctx, mongoClient)
	defer db.Close()

	virtualAccountModel := data.CreateVirtualAccount(db)

	err = virtualAccountModel.UpdatePin(virtualAccountApi.VirtualAccountNo, virtualAccountApi.Pin)
	if err != nil {
		virtualAccountApi.Error(w, err)
		return
	} else {
		vAccount, err := virtualAccountModel.FindVirtualAccountByNo(virtualAccountApi.VirtualAccountNo)
		if err != nil {
			virtualAccountApi.Error(w, err)
			return
		} else {
			virtualAccountApi.Json(w, vAccount, http.StatusOK)
			return
		}
	}
}

func (virtualAccountApi VirtualAccountApi) RemoveVirtualAccount(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&virtualAccountApi)
	if err != nil {
		virtualAccountApi.Error(w, err)
		return
	}

	mongoClient, ctx, db := virtualAccountApi.Connection(w, r)
	defer connection.CloseMongoDb(ctx, mongoClient)
	defer db.Close()

	virtualAccountModel := data.CreateVirtualAccount(db)

	err = virtualAccountModel.RemoveByNo(virtualAccountApi.VirtualAccountNo)
	if err != nil {
		virtualAccountApi.Error(w, err)
		return
	} else {
		virtualAccountApi.Empty(w, http.StatusOK)
		return
	}

}
