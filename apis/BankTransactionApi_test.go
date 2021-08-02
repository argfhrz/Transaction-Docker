package apis

import (
	"bank-account/config"
	"bank-account/connection"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestBankAccountTransactionApi_Get(t *testing.T) {

	db := connection.OpenConnection(config.DEV)
	defer db.Close()

	router := mux.NewRouter()

	bankTransactionApi := BankTransactionApi{}

	router.HandleFunc("/bank-transactions", bankTransactionApi.GetBankTransaction).Methods("GET")

	request, err := http.NewRequest("GET", "/bank-transactions", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	resp := recorder.Body.Bytes()
	log.Println(string(resp))

}

func TestBankAccountTransactionApi_GetByID(t *testing.T) {

	db := connection.OpenConnection(config.DEV)
	defer db.Close()

	router := mux.NewRouter()

	bankTransactionApi := BankTransactionApi{}

	router.HandleFunc("/bank-transaction/id", bankTransactionApi.GetBankTransactionByID).Methods("GET")

	request, err := http.NewRequest("GET", "/bank-transaction/id?bankTransactionNo=ab2a1ef1-6c95-4452-b21f-20332abbbef5", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	resp := recorder.Body.Bytes()
	log.Println(string(resp))

}

func TestBankAccountTransactionApi_GetByNo(t *testing.T) {

	db := connection.OpenConnection(config.DEV)
	defer db.Close()

	router := mux.NewRouter()

	bankTransactionApi := BankTransactionApi{}

	router.HandleFunc("/bank-transactions/history", bankTransactionApi.GetBankTransactionByNo).Methods("GET")

	request, err := http.NewRequest("GET", "/bank-transactions/history?bankAccountNo=a053bc97-dedd-48ee-a27f-e1d2aaa04539", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	resp := recorder.Body.Bytes()
	log.Println(string(resp))

}
