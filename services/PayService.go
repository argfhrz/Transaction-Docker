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

type PayService struct {
	BaseService
}

func CreatePayService(client *mongo.Client, db *sql.DB) PayService {
	payService := PayService{}
	payService.MongoClient = client
	payService.DB = db
	return payService
}

func (service PayService) Pay(ctx context.Context, mrchAccountNo string, srcVaAccountNo string, payAmount float64, pin string) (*nosql.PayNoSql, error) {

	//GET MERCHANT VA ACCOUNT
	vaAccountModel := data.CreateVirtualAccount(service.DB)
	merchantAccount, err := vaAccountModel.FindVirtualAccountByNo(mrchAccountNo)
	if err != nil {
		log.Println(err)
		return nil, errors.New("virtual_account_not_found")
	}

	//GET SRC VA ACCOUNT
	srcVaAccount, err := vaAccountModel.FindVirtualAccountByNo(srcVaAccountNo)
	if err != nil {
		log.Println(err)
		return nil, errors.New("virtual_account_not_found")
	}

	//CHECK CONDITION
	if srcVaAccountNo == mrchAccountNo {
		return nil, errors.New("invalid_tujuan_transfer")
	} else {
		if srcVaAccount.Saldo < payAmount {
			return nil, errors.New("saldo_tidak_mencukupi")
		} else {
			mrchSaldo := merchantAccount.Saldo
			srcSaldo := srcVaAccount.Saldo
			pin := base64.StdEncoding.EncodeToString([]byte(pin))

			if srcVaAccount.VirtualAccountNo == srcVaAccountNo && merchantAccount.VirtualAccountNo == mrchAccountNo && pin == srcVaAccount.Pin {

				srcSaldo -= payAmount
				err := vaAccountModel.UpdateSaldo(srcVaAccountNo, srcSaldo)
				if err != nil {
					log.Println(err)
					return nil, err
				}

				mrchSaldo += payAmount
				err = vaAccountModel.UpdateSaldo(mrchAccountNo, mrchSaldo)
				if err != nil {
					log.Println(err)
					return nil, err
				}

				payNoSql := nosql.CreatePayNoSql(service.MongoClient)
				pay, err := payNoSql.AddPay(ctx, mrchAccountNo, merchantAccount.AccountName, srcVaAccountNo, srcVaAccount.AccountName, payAmount)
				if err != nil {
					log.Println(err)
					return nil, err
				}

				description := mrchAccountNo + " | " + merchantAccount.AccountName
				vaAccountTransactionNoSql := nosql.CreateVirtualAccountTransactionNoSql(service.MongoClient)
				vaAccountTransaction, err := vaAccountTransactionNoSql.AddVirtualAccountTransaction(ctx, srcVaAccountNo, srcVaAccount.AccountName, payAmount, nosql.TRANSACTION_TYPE_PAY, description, pay.ID)
				if err != nil {
					log.Println(err)
					return nil, err
				}
				log.Println(vaAccountTransaction)

				return pay, nil

			} else {
				return nil, errors.New("invalid_account_no_atau_pin")
			}

		}

	}
}
