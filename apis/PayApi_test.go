package apis

import (
	"bytes"
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

func TestPayApi_Post(t *testing.T) {

	ctx := context.TODO()

	conn, err := connection.OpenMongoDb(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer connection.CloseMongoDb(ctx, conn)

	payModel := nosql.CreatePayNoSql(conn)
	err = payModel.Truncate(ctx)
	if err != nil {
		t.Fatal(err)
	}

	router := mux.NewRouter()

	payApi := PayApi{}

	router.HandleFunc("/top-up", payApi.PostPay).Methods("POST")

	payApi.MerchantVaAccountNo = "0812234567-00003"
	payApi.SrcVaAccountNo = "0812313112-00005"
	payApi.PayAmount = float64(2000000)
	payApi.Pin = "123123"

	body, err := json.Marshal(&payApi)
	if err != nil {
		t.Fatal(err)
	}

	request, err := http.NewRequest("POST", "/top-up", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	token := "JDJhJDEwJFo3eFpzOTZOR1VaTXFNS0p1WG96Sy5ZVk5tV3FHZ1BvMXpTWWthbHQ3RDF6eWxMcTR3QlVX"

	request.Header.Add("Authorization", token)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	resp := recorder.Body.Bytes()
	log.Println(string(resp))

	if recorder.Code != http.StatusOK {
		t.Fatal(string(resp))
	}

	newPay := nosql.PayNoSql{}
	err = json.Unmarshal(resp, &newPay)
	if err != nil {
		t.Fatal(err)
	}
}
