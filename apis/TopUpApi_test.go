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

func TestTopUpApi_Post(t *testing.T) {

	ctx := context.TODO()

	conn, err := connection.OpenMongoDb(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer connection.CloseMongoDb(ctx, conn)

	topUpModel := nosql.CreateTopUpNoSql(conn)
	err = topUpModel.Truncate(ctx)
	if err != nil {
		t.Fatal(err)
	}

	router := mux.NewRouter()

	topUpApi := TopUpApi{}

	router.HandleFunc("/top-up", topUpApi.PostTopUp).Methods("POST")

	topUpApi.BankCode = "a053bc97-dedd-48ee-a27f-e1d2aaa04539"
	topUpApi.VaAccountNo = "0812234567-00003"
	topUpApi.TopUpAmount = float64(2000000)
	topUpApi.Pin = "98765"

	body, err := json.Marshal(&topUpApi)
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

	newTopUp := nosql.TopUpNoSql{}
	err = json.Unmarshal(resp, &newTopUp)
	if err != nil {
		t.Fatal(err)
	}
}
