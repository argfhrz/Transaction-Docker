package apis

import (
	"bank-account/config"
	"bank-account/connection"
	"bank-account/data"
	"encoding/json"
	"log"
	"net/http"
)

type BankAccountApi struct {
	BaseApi

	data.BankAccount
}

func (bankAccountApi BankAccountApi) GetBankAccount(w http.ResponseWriter, r *http.Request) {

	db := connection.OpenConnection(config.DEV)
	defer db.Close()

	bankAccountModel := data.CreateBankAccount(db)
	bankAccounts, err := bankAccountModel.GetListBankAccount()
	if err != nil {
		bankAccountApi.Error(w, err)
		return
	} else {
		bankAccountApi.Json(w, bankAccounts, http.StatusOK)
		return
	}

}

func (bankAccountApi BankAccountApi) GetBankAccountByID(w http.ResponseWriter, r *http.Request) {

	db := connection.OpenConnection(config.DEV)
	defer db.Close()

	bankAccountNo := bankAccountApi.QueryParam(r, "bankAccountNo")

	bankAccountModel := data.CreateBankAccount(db)
	bankAccount, err := bankAccountModel.FindBankAccountByNo(bankAccountNo)
	if err != nil {
		bankAccountApi.Error(w, err)
		return
	} else {
		bankAccountApi.Json(w, bankAccount, http.StatusOK)
		return
	}

}

func (bankAccountApi BankAccountApi) PostBankAccount(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&bankAccountApi)
	if err != nil {
		bankAccountApi.Error(w, err)
		return
	}

	db := connection.OpenConnection(config.DEV)
	defer db.Close()

	bankAccountModel := data.CreateBankAccount(db)

	bankAccountNo, err := bankAccountModel.AddBankAccount(bankAccountApi.BankAccountOwner, bankAccountApi.Saldo)
	if err != nil {
		bankAccountApi.Error(w, err)
		return
	} else {
		bankAccount, err := bankAccountModel.FindBankAccountByNo(bankAccountNo)
		if err != nil {
			bankAccountApi.Error(w, err)
			return
		} else {
			bankAccountApi.Json(w, bankAccount, http.StatusOK)
			log.Println(bankAccountNo)
			return
		}
	}

}

func (bankAccountApi BankAccountApi) UpdateBankAccountIdentity(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&bankAccountApi)
	if err != nil {
		bankAccountApi.Error(w, err)
		return
	}

	db := connection.OpenConnection(config.DEV)
	defer db.Close()

	bankAccountModel := data.CreateBankAccount(db)

	err = bankAccountModel.UpdateIdentity(bankAccountApi.BankAccountNo, bankAccountApi.BankAccountOwner)
	if err != nil {
		bankAccountApi.Error(w, err)
		return
	} else {
		vAccount, err := bankAccountModel.FindBankAccountByNo(bankAccountApi.BankAccountNo)
		if err != nil {
			bankAccountApi.Error(w, err)
			return
		} else {
			bankAccountApi.Json(w, vAccount, http.StatusOK)
			return
		}
	}
}

func (bankAccountApi BankAccountApi) RemoveBankAccount(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&bankAccountApi)
	if err != nil {
		bankAccountApi.Error(w, err)
		return
	}

	db := connection.OpenConnection(config.DEV)
	defer db.Close()

	bankAccountModel := data.CreateBankAccount(db)

	err = bankAccountModel.RemoveByNo(bankAccountApi.BankAccountNo)
	if err != nil {
		bankAccountApi.Error(w, err)
		return
	} else {
		bankAccountApi.Empty(w, http.StatusOK)
		return
	}

}
