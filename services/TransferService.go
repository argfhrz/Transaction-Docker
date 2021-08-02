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

type TransferService struct {
	BaseService
}

func CreateTransferService(client *mongo.Client, db *sql.DB) TransferService {
	transferService := TransferService{}
	transferService.MongoClient = client
	transferService.DB = db
	return transferService
}

func (service TransferService) Transfer(ctx context.Context, srcVaAccountNo string, destVaAccountNo string, transferAmount float64, pin string) (*nosql.TransferNoSql, error) {

	//GET SRC VA ACCOUNT
	vaAccountModel := data.CreateVirtualAccount(service.DB)
	srcVaAccount, err := vaAccountModel.FindVirtualAccountByNo(srcVaAccountNo)
	if err != nil {
		log.Println(err)
		return nil, errors.New("virtual_account_not_found")
	}

	//GET DEST VA ACCOUNT
	destVaAccount, err := vaAccountModel.FindVirtualAccountByNo(destVaAccountNo)
	if err != nil {
		log.Println(err)
		return nil, errors.New("virtual_account_not_found")
	}

	//CHECK MONEY
	if srcVaAccountNo == destVaAccountNo {
		return nil, errors.New("invalid_tujuan_transfer")
	} else {
		if srcVaAccount.Saldo < transferAmount {
			return nil, errors.New("saldo_tidak_mencukupi")
		} else {

			srcSaldo := srcVaAccount.Saldo
			destSaldo := destVaAccount.Saldo
			pin := base64.StdEncoding.EncodeToString([]byte(pin))

			if srcVaAccount.VirtualAccountNo == srcVaAccountNo && destVaAccount.VirtualAccountNo == destVaAccountNo && pin == srcVaAccount.Pin {

				srcSaldo -= transferAmount
				err := vaAccountModel.UpdateSaldo(srcVaAccountNo, srcSaldo)
				if err != nil {
					log.Println(err)
					return nil, err
				}

				destSaldo += transferAmount
				err = vaAccountModel.UpdateSaldo(destVaAccountNo, destSaldo)
				if err != nil {
					log.Println(err)
					return nil, err
				}

				transferNoSql := nosql.CreateTransferNoSql(service.MongoClient)
				transfer, err := transferNoSql.AddTransfer(ctx, srcVaAccountNo, srcVaAccount.AccountName, destVaAccountNo, destVaAccount.AccountName, transferAmount)
				if err != nil {
					log.Println(err)
					return nil, err
				}

				description := "Transfer to :" + destVaAccountNo
				vaAccountTransactionNoSql := nosql.CreateVirtualAccountTransactionNoSql(service.MongoClient)
				vaAccountTransaction, err := vaAccountTransactionNoSql.AddVirtualAccountTransaction(ctx, srcVaAccountNo, srcVaAccount.AccountName, transferAmount, nosql.TRANSACTION_TYPE_TRANSFER, description, transfer.ID)
				if err != nil {
					log.Println(err)
					return nil, err
				}
				log.Println(vaAccountTransaction)

				return transfer, nil

			} else {
				return nil, errors.New("invalid_account_no_atau_pin")
			}

		}

	}

}
