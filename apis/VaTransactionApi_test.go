package apis

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"virtual-account/config"
	"virtual-account/connection"

	"github.com/gorilla/mux"
)

func TestVirtualAccountTransactionApi_Get(t *testing.T) {

	db := connection.OpenConnection(config.DEV)
	defer db.Close()

	router := mux.NewRouter()

	vaTransactionApi := VaTransactionApi{}

	router.HandleFunc("/virtual-transactions", vaTransactionApi.GetVaTransaction).Methods("GET")

	request, err := http.NewRequest("GET", "/virtual-transactions", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	resp := recorder.Body.Bytes()
	log.Println(string(resp))

}

func TestVirtualAccountTransactionApi_GetByID(t *testing.T) {

	db := connection.OpenConnection(config.DEV)
	defer db.Close()

	router := mux.NewRouter()

	vaTransactionApi := VaTransactionApi{}

	router.HandleFunc("/virtual-transaction/id", vaTransactionApi.GetVaTransactionByID).Methods("GET")

	request, err := http.NewRequest("GET", "/virtual-transaction/id?vaTransactionNo=b1cbbe91-34fc-42df-8929-0e978bf2d580", nil)
	if err != nil {
		t.Fatal(err)
	}

	token := "JDJhJDEwJFo3eFpzOTZOR1VaTXFNS0p1WG96Sy5ZVk5tV3FHZ1BvMXpTWWthbHQ3RDF6eWxMcTR3QlVX"

	request.Header.Add("Authorization", token)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	resp := recorder.Body.Bytes()
	log.Println(string(resp))

}

func TestVirtualAccountTransactionApi_GetByNo(t *testing.T) {

	db := connection.OpenConnection(config.DEV)
	defer db.Close()

	router := mux.NewRouter()

	vaTransactionApi := VaTransactionApi{}

	router.HandleFunc("/virtual-transaction/history", vaTransactionApi.GetVaTransactionByNo).Methods("GET")

	request, err := http.NewRequest("GET", "/virtual-transaction/history?virtualAccountNo=0812234567-00003", nil)
	if err != nil {
		t.Fatal(err)
	}

	token := "JDJhJDEwJFo3eFpzOTZOR1VaTXFNS0p1WG96Sy5ZVk5tV3FHZ1BvMXpTWWthbHQ3RDF6eWxMcTR3QlVX"

	request.Header.Add("Authorization", token)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	resp := recorder.Body.Bytes()
	log.Println(string(resp))

}
