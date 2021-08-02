package apis

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"virtual-account/config"
	"virtual-account/connection"
	"virtual-account/data"

	"github.com/gorilla/mux"
)

func TestVirtualAccountApi_Post(t *testing.T) {

	db := connection.OpenConnection(config.DEV)
	defer db.Close()

	virtualAccountModel := data.CreateVirtualAccount(db)
	err := virtualAccountModel.Truncate()
	if err != nil {
		t.Fatal(err)
	}

	router := mux.NewRouter()

	virtualAccountApi := VirtualAccountApi{}

	router.HandleFunc("/register", virtualAccountApi.PostVirtualAccount).Methods("POST")

	virtualAccount := data.VirtualAccount{}
	virtualAccount.PhoneNumber = "0812341141"
	virtualAccount.Email = "user1@gmail.com"
	virtualAccount.AccountName = "user1"
	virtualAccount.Pin = "123456"
	virtualAccount.Password = "123"

	body, err := json.Marshal(&virtualAccount)
	if err != nil {
		t.Fatal(err)
	}

	request, err := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
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

	newVirtualAccount := data.VirtualAccount{}
	err = json.Unmarshal(resp, &newVirtualAccount)
	if err != nil {
		t.Fatal(err)
	}
}

func TestVirtualAccountApi_Get(t *testing.T) {

	db := connection.OpenConnection(config.DEV)
	defer db.Close()

	router := mux.NewRouter()

	virtualAccountApi := VirtualAccountApi{}

	router.HandleFunc("/virtual-accounts", virtualAccountApi.GetVirtualAccount).Methods("GET")

	request, err := http.NewRequest("GET", "/virtual-accounts", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	resp := recorder.Body.Bytes()
	log.Println(string(resp))

}

func TestVirtualAccountApi_GetByID(t *testing.T) {

	db := connection.OpenConnection(config.DEV)
	defer db.Close()

	router := mux.NewRouter()

	virtualAccountApi := VirtualAccountApi{}

	router.HandleFunc("/virtual-account/id", virtualAccountApi.GetVirtualAccountByID).Methods("GET")

	request, err := http.NewRequest("GET", "/virtual-account/id?virtualAccountNo=0812234567-00003", nil)
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

func TestVirtualAccountApi_UpdatePin(t *testing.T) {

	db := connection.OpenConnection(config.DEV)
	defer db.Close()

	router := mux.NewRouter()

	virtualAccountApi := VirtualAccountApi{}

	router.HandleFunc("/virtual-account/pin", virtualAccountApi.UpdateVirtualAccountPin).Methods("POST")

	virtualAccount := data.VirtualAccount{}
	virtualAccount.VirtualAccountNo = "0812234567-00003"
	virtualAccount.Pin = "343434"

	body, err := json.Marshal(&virtualAccount)
	if err != nil {
		t.Fatal(err)
	}

	request, err := http.NewRequest("POST", "/virtual-account/pin", bytes.NewBuffer(body))
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

	newVirtualAccount := data.VirtualAccount{}
	err = json.Unmarshal(resp, &newVirtualAccount)
	if err != nil {
		t.Fatal(err)
	}

}

func TestVirtualAccountApi_UpdateName(t *testing.T) {

	db := connection.OpenConnection(config.DEV)
	defer db.Close()

	router := mux.NewRouter()

	virtualAccountApi := VirtualAccountApi{}

	router.HandleFunc("/virtual-account/name", virtualAccountApi.UpdateVirtualAccountPin).Methods("POST")

	virtualAccount := data.VirtualAccount{}
	virtualAccount.VirtualAccountNo = "0812234567-00003"
	virtualAccount.AccountName = "newName1"

	body, err := json.Marshal(&virtualAccount)
	if err != nil {
		t.Fatal(err)
	}

	request, err := http.NewRequest("POST", "/virtual-account/name", bytes.NewBuffer(body))
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

	newVirtualAccount := data.VirtualAccount{}
	err = json.Unmarshal(resp, &newVirtualAccount)
	if err != nil {
		t.Fatal(err)
	}

}

func TestVirtualAccountApi_Remove(t *testing.T) {

	db := connection.OpenConnection(config.DEV)
	defer db.Close()

	router := mux.NewRouter()

	virtualAccountApi := VirtualAccountApi{}

	router.HandleFunc("/delete-account", virtualAccountApi.RemoveVirtualAccount).Methods("DELETE")

	virtualAccount := data.VirtualAccount{}
	virtualAccount.VirtualAccountNo = "0812234567-00003"

	body, err := json.Marshal(&virtualAccount)
	if err != nil {
		t.Fatal(err)
	}

	request, err := http.NewRequest("DELETE", "/delete-account", bytes.NewBuffer(body))
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

}
