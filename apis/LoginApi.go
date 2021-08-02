package apis

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
	"virtual-account/config"
	"virtual-account/connection"
	"virtual-account/nosql"
	"virtual-account/services"
)

type LoginApi struct {
	BaseApi

	nosql.LoginNoSql
}

func (loginApi LoginApi) Login(w http.ResponseWriter, r *http.Request) {

	authorization := r.Header.Get("Authorization")
	authorization = strings.Replace(authorization, "Basic ", "", -1)
	auth, err := base64.StdEncoding.DecodeString(authorization)
	if err != nil {
		loginApi.Error(w, err)
		return
	}
	auth1 := strings.Split(string(auth), ":")

	ctx := context.TODO()
	mongoClient, err := connection.OpenMongoDb(ctx)
	if err != nil {
		loginApi.Error(w, err)
		return
	}

	defer connection.CloseMongoDb(ctx, mongoClient)
	db := connection.OpenConnection(config.DEV)
	defer db.Close()

	loginService := services.CreateLoginService(mongoClient, db)
	encodedToken, err := loginService.Login(ctx, auth1[0], auth1[1])
	if err != nil {
		loginApi.Error(w, err)
		return
	} else {

		data := map[string]string{
			"token": encodedToken,
		}
		loginApi.Json(w, data, http.StatusOK)
		return
	}

}

func (loginApi LoginApi) Logout(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&loginApi)
	if err != nil {
		loginApi.Error(w, err)
		return
	}

	ctx := context.TODO()
	mongoClient, err := connection.OpenMongoDb(ctx)
	if err != nil {
		loginApi.Error(w, err)
		return
	}
	defer connection.CloseMongoDb(ctx, mongoClient)
	db := connection.OpenConnection(config.DEV)
	defer db.Close()

	loginService := services.CreateLoginService(mongoClient, db)
	err = loginService.Logout(ctx, loginApi.AccountNo)
	if err != nil {
		loginApi.Error(w, err)
		return
	}
}
