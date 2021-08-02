package apis

import (
	"bank-account/config"
	"bank-account/connection"
	"bank-account/data"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestBankAccountApi_Post(t *testing.T) {
	db := connection.OpenConnection(config.DEV)
	defer db.Close()

	bankAccountModel := data.CreateBankAccount(db)
	err := bankAccountModel.Truncate()
	if err != nil {
		t.Fatal(err)
	}

	router := mux.NewRouter()

	bankAccountApi := BankAccountApi{}

	router.HandleFunc("/bank-account", bankAccountApi.PostBankAccount).Methods("POST")

	bankAccount := data.BankAccount{}
	bankAccount.BankAccountOwner = "Owner1"
	bankAccount.Saldo = float64(1000000000)

	body, err := json.Marshal(&bankAccount)
	if err != nil {
		t.Fatal(err)
	}

	request, err := http.NewRequest("POST", "/bank-account", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	resp := recorder.Body.Bytes()
	log.Println(string(resp))

	if recorder.Code != http.StatusOK {
		t.Fatal(string(resp))
	}

	newBankAccount := data.BankAccount{}
	err = json.Unmarshal(resp, &newBankAccount)
	if err != nil {
		t.Fatal(err)
	}

}

func TestBankAccountApi_Get(t *testing.T) {

	db := connection.OpenConnection(config.DEV)
	defer db.Close()

	router := mux.NewRouter()

	bankAccountApi := BankAccountApi{}

	router.HandleFunc("/bank-accounts", bankAccountApi.GetBankAccount).Methods("GET")

	request, err := http.NewRequest("GET", "/bank-accounts", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	resp := recorder.Body.Bytes()
	log.Println(string(resp))

}

func TestBankAccountApi_GetByID(t *testing.T) {

	db := connection.OpenConnection(config.DEV)
	defer db.Close()

	router := mux.NewRouter()

	bankAccountApi := BankAccountApi{}

	router.HandleFunc("/bank-account/id", bankAccountApi.GetBankAccountByID).Methods("GET")

	request, err := http.NewRequest("GET", "/bank-account/id?bankAccountNo=a053bc97-dedd-48ee-a27f-e1d2aaa04539", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	resp := recorder.Body.Bytes()
	log.Println(string(resp))

}

func TestVirtualAccountApi_UpdateIdentity(t *testing.T) {

	db := connection.OpenConnection(config.DEV)
	defer db.Close()

	router := mux.NewRouter()

	bankAccountApi := BankAccountApi{}

	router.HandleFunc("/bank-account", bankAccountApi.UpdateBankAccountIdentity).Methods("POST")

	bankAccount := data.BankAccount{}
	bankAccount.BankAccountNo = "a053bc97-dedd-48ee-a27f-e1d2aaa04539"
	bankAccount.BankAccountOwner = "owner2"

	body, err := json.Marshal(&bankAccount)
	if err != nil {
		t.Fatal(err)
	}

	request, err := http.NewRequest("POST", "/bank-account", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	resp := recorder.Body.Bytes()
	log.Println(string(resp))

	if recorder.Code != http.StatusOK {
		t.Fatal(string(resp))
	}

	newBankAccount := data.BankAccount{}
	err = json.Unmarshal(resp, &newBankAccount)
	if err != nil {
		t.Fatal(err)
	}

}
