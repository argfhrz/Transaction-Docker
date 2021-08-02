package apis

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"virtual-account/connection"
	"virtual-account/nosql"

	"github.com/gorilla/mux"
)

func TestLoginApi(t *testing.T) {

	ctx := context.TODO()
	mongoClient, err := connection.OpenMongoDb(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer connection.CloseMongoDb(ctx, mongoClient)

	loginApi := LoginApi{}

	router := mux.NewRouter()

	router.HandleFunc("/login", loginApi.Login).Methods("POST")

	phoneNumber := "0812312312"
	password := "123"

	request, err := http.NewRequest("POST", "/login", nil)
	if err != nil {
		t.Fatal(err)
	}

	request.SetBasicAuth(phoneNumber, password)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	resp := recorder.Body.Bytes()
	log.Println(string(resp))

	if recorder.Code != http.StatusOK {
		t.Fatal(string(resp))
	}

	newLogin := nosql.LoginNoSql{}
	err = json.Unmarshal(resp, &newLogin)
	if err != nil {
		t.Fatal(err)
	}

}
