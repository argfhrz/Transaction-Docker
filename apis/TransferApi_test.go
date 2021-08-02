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

func TestTransferApi_Post(t *testing.T) {
	ctx := context.TODO()

	conn, err := connection.OpenMongoDb(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer connection.CloseMongoDb(ctx, conn)

	transferModel := nosql.CreateTransferNoSql(conn)
	err = transferModel.Truncate(ctx)
	if err != nil {
		t.Fatal(err)
	}

	router := mux.NewRouter()

	transferApi := TransferApi{}

	router.HandleFunc("/transfer", transferApi.PostTransfer).Methods("POST")

	transferApi.ScrVaAccountNo = "0812234567-00003"
	transferApi.DestVaAccountNo = "0812313112-00005"
	transferApi.TransferAmount = float64(2000000)
	transferApi.Pin = "98765"

	body, err := json.Marshal(&transferApi)
	if err != nil {
		t.Fatal(err)
	}

	request, err := http.NewRequest("POST", "/transfer", bytes.NewBuffer(body))
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

	newTransfer := nosql.TransferNoSql{}
	err = json.Unmarshal(resp, &newTransfer)
	if err != nil {
		t.Fatal(err)
	}
}
