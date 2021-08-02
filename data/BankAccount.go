package data

import (
	"database/sql"
	"errors"
	"log"
	"time"
	"virtual-account/helpers"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type BankAccount struct {
	BaseData

	BankAccountNo    string  `json:"bankAccountNo"`
	BankAccountOwner string  `json:"bankAccountOwner"`
	Saldo            float64 `json:"saldo"`

	CreatedAt time.Time `json:"createdAt"`
}

func CreateBankAccount(db *sql.DB) BankAccount {
	bankAccount := BankAccount{}
	bankAccount.DB = db
	return bankAccount
}

func CreateBankAccountWithTransaction(transaction *sql.Tx) BankAccount {
	bankAccount := BankAccount{}
	bankAccount.Transaction = transaction
	bankAccount.UseTransaction = true
	return bankAccount
}

func (bankAccount BankAccount) AddBankAccount(bankAccountOwner string, saldo float64) (string, error) {
	sql := `insert into bank_accounts (bank_account_no, bank_account_owner, saldo, created_at) values 
	($1, $2, $3, $4)`

	bankAccountNo := uuid.New().String()
	result, err := bankAccount.Exec(sql, bankAccountNo, bankAccountOwner, saldo, time.Now().UTC())
	if err != nil {
		return "", err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return "", err
	}

	log.Println("BankAccount.Add.AffectedRows=", affectedRows)
	return bankAccountNo, nil
}

func (bankAccount BankAccount) fetchRow(cursor *sql.Rows) (BankAccount, error) {
	b := BankAccount{}
	err := cursor.Scan(&b.BankAccountNo, &b.BankAccountOwner, &b.Saldo, &b.CreatedAt)
	if err != nil {
		return BankAccount{}, err
	}
	return b, nil
}

func (bankAccount BankAccount) selectQuery() string {
	sql := `select * from bank_accounts`
	return sql
}

func (bankAccount BankAccount) GetListBankAccount() ([]BankAccount, error) {
	//                0             1
	sql := bankAccount.selectQuery()
	cursor, err := bankAccount.Query(sql)
	if err != nil {
		return nil, err
	}

	bankAccounts := []BankAccount{}

	for cursor.Next() {
		c, err := bankAccount.fetchRow(cursor)
		if err != nil {
			return nil, err
		}
		bankAccounts = append(bankAccounts, c)
	}

	return bankAccounts, nil
}

func (bankAccount BankAccount) FindBankAccountByNo(bankAccountNo string) (*BankAccount, error) {

	if helpers.Empty(bankAccountNo) {
		return nil, errors.New("invalid_bank_account")
	}

	sql := bankAccount.selectQuery()
	sql += ` where bank_account_no=$1`

	cursor, err := bankAccount.Query(sql, bankAccountNo)
	if err != nil {
		return nil, err
	}

	if cursor.Next() {
		c, err := bankAccount.fetchRow(cursor)
		if err != nil {
			return nil, err
		}
		return &c, nil
	}

	return nil, errors.New("BankAccount tidak ditemukan")

}

func (bankAccount BankAccount) RemoveByNo(bankAccountNo string) error {

	sql := "delete from bank_accounts where bank_account_no=$1"
	result, err := bankAccount.Exec(sql, bankAccountNo)
	if err != nil {
		return err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	log.Println("BankAccount.Remove.AffectedRows=", affectedRows)
	return nil
}

func (bankAccount BankAccount) UpdateIdentity(bankAccountNo string, bankAccountOwner string) error {

	sql := `update bank_accounts 
	set bank_account_owner=$2
	
	where bank_account_no=$1
	`
	result, err := bankAccount.Exec(sql, bankAccountNo, bankAccountOwner)
	if err != nil {
		return err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	log.Println("BankAccount.Update.AffectedRows=", affectedRows)

	return nil

}

func (bankAccount BankAccount) UpdateSaldo(bankAccountNo string, saldo float64) error {

	sql := `update bank_accounts 
	set saldo=$2

	where bank_account_no=$1
	`
	result, err := bankAccount.Exec(sql, bankAccountNo, saldo)
	if err != nil {
		return err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	log.Println("BankAccount.Update.AffectedRows=", affectedRows)

	return nil

}

func (bankAccount BankAccount) Truncate() error {

	sql := `delete from bank_accounts`
	result, err := bankAccount.Exec(sql)
	if err != nil {
		return err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	log.Println("BankAccount.Truncate.AffectedRows=", affectedRows)
	return nil

}

func (bankAccount BankAccount) Migrate() error {

	sqlDropTable := "DROP TABLE IF EXISTS public.bank_accounts CASCADE"

	sqlCreateTable := `
CREATE TABLE public.bank_accounts
(
    bank_account_no character varying(40) NOT NULL,
    bank_account_owner character varying(60) NOT NULL,
	saldo double precision NOT NULL,
    created_at timestamp with time zone NOT NULL,
    CONSTRAINT bank_accounts_pkey PRIMARY KEY (bank_account_no)
)

TABLESPACE pg_default;

ALTER TABLE public.bank_accounts
    OWNER to postgres;
	`

	_, err := bankAccount.DB.Exec(sqlDropTable)
	if err != nil {
		return err
	}

	_, err = bankAccount.DB.Exec(sqlCreateTable)
	if err != nil {
		return err
	}

	return nil

}
