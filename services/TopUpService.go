package services

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"log"
	"virtual-account/data"
	"virtual-account/nosql"

	"go.mongodb.org/mongo-driver/mongo"
)

type TopUpService struct {
	BaseService
}

func CreateTopUpService(mongoClient *mongo.Client, db *sql.DB) TopUpService {
	topUpService := TopUpService{}
	topUpService.MongoClient = mongoClient
	topUpService.DB = db
	return topUpService
}

func (service TopUpService) TopUp(ctx context.Context, bankAccountNo string, vaAccountNo string, pin string, transactionAmount float64) (*nosql.TopUpNoSql, *nosql.BankTransactionNoSql, *nosql.VirtualAccountTransactionNoSql, error) {

	//FIND BANK ACCOUNT
	bankAccountModel := data.CreateBankAccount(service.DB)
	findBankAccount, err := bankAccountModel.FindBankAccountByNo(bankAccountNo)
	if err != nil {
		log.Println(err)
		return nil, nil, nil, errors.New("bank_account_not_found")
	}

	//FIND VA ACCOUNT
	vaAccountModel := data.CreateVirtualAccount(service.DB)
	findVaAccount, err := vaAccountModel.FindVirtualAccountByNo(vaAccountNo)
	if err != nil {
		log.Println(err)
		return nil, nil, nil, err
	}

	bankSaldo := findBankAccount.Saldo
	vaSaldo := findVaAccount.Saldo
	pin = base64.StdEncoding.EncodeToString([]byte(pin))

	if findBankAccount.BankAccountNo == bankAccountNo && findVaAccount.VirtualAccountNo == vaAccountNo && findVaAccount.Pin == pin {

		bankSaldo -= transactionAmount
		vaSaldo += transactionAmount

		err = bankAccountModel.UpdateSaldo(bankAccountNo, bankSaldo)
		if err != nil {
			log.Println(err)
			return nil, nil, nil, err
		}

		err = vaAccountModel.UpdateSaldo(vaAccountNo, vaSaldo)
		if err != nil {
			log.Println(err)
			return nil, nil, nil, err
		}

		topUpNoSql := nosql.CreateTopUpNoSql(service.MongoClient)
		topUp, err := topUpNoSql.AddTopUp(ctx, bankAccountNo, findBankAccount.BankAccountOwner, vaAccountNo, findVaAccount.AccountName, transactionAmount)
		if err != nil {
			log.Println(err)
			return nil, nil, nil, err
		}

		references := vaAccountNo + " | " + findVaAccount.AccountName
		bankTransactionNoSql := nosql.CreateBankTransactionNoSql(service.MongoClient)
		bankTransaction, err := bankTransactionNoSql.AddBankTransaction(ctx, findBankAccount.BankAccountNo, findBankAccount.BankAccountOwner, transactionAmount, references)
		if err != nil {
			log.Println(err)
			return nil, nil, nil, err
		}
		log.Println(bankTransaction)

		description := findBankAccount.BankAccountNo + " | " + findBankAccount.BankAccountOwner
		virtualTransactionNoSql := nosql.CreateVirtualAccountTransactionNoSql(service.MongoClient)
		virtualTransaction, err := virtualTransactionNoSql.AddVirtualAccountTransaction(ctx, vaAccountNo, findVaAccount.AccountName, transactionAmount, nosql.TRANSACTION_TYPE_TOP_UP, description, topUp.ID)
		if err != nil {
			log.Println(err)
			return nil, nil, nil, err
		}
		log.Println(virtualTransaction)

		return topUp, bankTransaction, virtualTransaction, nil

	} else {
		return nil, nil, nil, errors.New("identitas_va_account_atau_bank_account_atau_pin")
	}

}
