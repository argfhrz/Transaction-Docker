package main

import (
	"log"
	"net/http"
	"virtual-account/apis"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to Virtual Account"))
	}).Methods("GET")

	loginApi := apis.LoginApi{}
	router.HandleFunc("/login", loginApi.Login).Methods("POST")
	router.HandleFunc("/logout", loginApi.Logout).Methods("POST")

	virtualAccountApi := apis.VirtualAccountApi{}
	router.HandleFunc("/register", virtualAccountApi.PostVirtualAccount).Methods("POST")
	router.HandleFunc("/virtual-accounts", virtualAccountApi.GetVirtualAccount).Methods("GET")
	router.HandleFunc("/virtual-account/id", virtualAccountApi.GetVirtualAccountByID).Methods("GET")
	router.HandleFunc("/virtual-account/pin", virtualAccountApi.UpdateVirtualAccountPin).Methods("PUT")
	router.HandleFunc("/virtual-account/name", virtualAccountApi.UpdateVirtualAccountName).Methods("PUT")
	router.HandleFunc("/delete-account", virtualAccountApi.RemoveVirtualAccount).Methods("DELETE")

	topUpApi := apis.TopUpApi{}
	router.HandleFunc("/top-up", topUpApi.PostTopUp).Methods("POST")
	router.HandleFunc("/top-up", topUpApi.RemoveTopUp).Methods("DELETE")

	transferApi := apis.TransferApi{}
	router.HandleFunc("/transfer", transferApi.PostTransfer).Methods("POST")
	router.HandleFunc("/transfer", transferApi.RemoveTransfer).Methods("DELETE")

	payApi := apis.PayApi{}
	router.HandleFunc("/pay", payApi.PostPay).Methods("POST")
	router.HandleFunc("/pay", payApi.RemovePay).Methods("DELETE")

	vaTransactionApi := apis.VaTransactionApi{}
	router.HandleFunc("/virtual-transactions", vaTransactionApi.GetVaTransaction).Methods("GET")
	router.HandleFunc("/virtual-transaction/id", vaTransactionApi.GetVaTransactionByID).Methods("GET")
	router.HandleFunc("/virtual-transaction/history", vaTransactionApi.GetVaTransactionByNo).Methods("GET")
	router.HandleFunc("/virtual-transaction/", vaTransactionApi.RemoveVirtualTransaction).Methods("DELETE")

	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})

	log.Println("Server listening to 0.0.0.0:8200")
	err := http.ListenAndServe("0.0.0.0:8200", handlers.CORS(headers, methods, origins)(router))
	if err != nil {
		log.Fatal(err)
	}
}
