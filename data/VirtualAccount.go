package data

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"log"
	"strconv"
	"time"
	"virtual-account/helpers"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type VirtualAccount struct {
	BaseData

	VirtualAccountNo string     `json:"virtualAccountNo"`
	PhoneNumber      string     `json:"phoneNumber"`
	Email            string     `json:"email"`
	AccountName      string     `json:"accountName"`
	Pin              string     `json:"pin"`
	Password         string     `json:"password"`
	Saldo            float64    `json:"saldo"`
	SeqNo            int        `json:"seqNo"`
	CreatedAt        time.Time  `json:"createdAt"`
	UpdatedAt        *time.Time `json:"updatedAt"`
}

type VirtualAccountNoPassword struct {
	BaseData

	VirtualAccountNo string     `json:"virtualAccountNo"`
	PhoneNumber      string     `json:"phoneNumber"`
	Email            string     `json:"email"`
	AccountName      string     `json:"accountName"`
	CreatedAt        time.Time  `json:"createdAt"`
	UpdatedAt        *time.Time `json:"updatedAt"`
}

type SeqNo struct {
	SeqNo int `json:"seqNo"`
}

func CreateVirtualAccount(db *sql.DB) VirtualAccount {
	virtualAccount := VirtualAccount{}
	virtualAccount.DB = db
	return virtualAccount
}

func CreateVirtualAccountWithTransaction(transaction *sql.Tx) VirtualAccount {
	virtualAccount := VirtualAccount{}
	virtualAccount.Transaction = transaction
	virtualAccount.UseTransaction = true
	return virtualAccount
}

func (virtualAccount VirtualAccount) Add(phoneNumber string, email string,
	accountName string, pin string, password string) (string, error) {

	seqNo, err := virtualAccount.GetSeqNo()
	if err != nil {
		return "", err
	}
	number := seqNo.SeqNo + 1

	virtualAccountNo := ""

	if number < 10 {
		virtualAccountNo = phoneNumber + "-0000" + strconv.Itoa(number)
	} else if number < 100 {
		virtualAccountNo = phoneNumber + "-000" + strconv.Itoa(number)
	} else if number < 1000 {
		virtualAccountNo = phoneNumber + "-00" + strconv.Itoa(number)
	} else if number < 10000 {
		virtualAccountNo = phoneNumber + "-0" + strconv.Itoa(number)
	} else {
		virtualAccountNo = phoneNumber + "-" + strconv.Itoa(number)
	}

	accountNo, err := virtualAccount.AddVirtualAccount(virtualAccountNo, phoneNumber, email, accountName, pin, password, number)
	if err != nil {
		return "", err
	}

	return accountNo, nil

}

func (virtualAccount VirtualAccount) AddVirtualAccount(virtualAccountNo string, phoneNumber string, email string,
	accountName string, pin string, password string, seqNo int) (string, error) {
	sql := `insert into virtual_accounts (virtual_account_no, phone_number, email, 
		account_name, pin, password, saldo,
		seq_no, created_at) values 
	($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	pin = base64.StdEncoding.EncodeToString([]byte(pin))
	encryptedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	result, err := virtualAccount.Exec(sql, virtualAccountNo, phoneNumber, email, accountName, pin, string(encryptedPass), float64(0), seqNo, time.Now().UTC())
	if err != nil {
		return "", err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return "", err
	}

	log.Println("VirtualAccount.Add.AffectedRows=", affectedRows)
	return virtualAccountNo, nil

}

func (virtualAccount VirtualAccount) fetchRow(cursor *sql.Rows) (VirtualAccountNoPassword, error) {
	b := VirtualAccountNoPassword{}
	err := cursor.Scan(&b.VirtualAccountNo, &b.PhoneNumber, &b.Email, &b.AccountName, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		return VirtualAccountNoPassword{}, err
	}
	return b, nil
}

func (virtualAccount VirtualAccount) fetchRowAll(cursor *sql.Rows) (VirtualAccount, error) {
	b := VirtualAccount{}
	err := cursor.Scan(&b.VirtualAccountNo, &b.PhoneNumber, &b.Email, &b.AccountName, &b.Pin, &b.Password, &b.Saldo, &b.SeqNo, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		return VirtualAccount{}, err
	}
	return b, nil
}

func (virtualAccount VirtualAccount) fetchRowFind(cursor *sql.Rows) (VirtualAccount, error) {
	b := VirtualAccount{}
	err := cursor.Scan(&b.VirtualAccountNo, &b.Password)
	if err != nil {
		return VirtualAccount{}, err
	}
	return b, nil
}

func (virtualAccount VirtualAccount) selectQuery() string {
	sql := `select virtual_account_no, phone_number, email, account_name, created_at, updated_at from virtual_accounts`
	return sql
}

func (virtualAccount VirtualAccount) selectQueryAll() string {
	sql := `select virtual_account_no, phone_number, email, account_name, pin, password, saldo, seq_no, created_at, updated_at from virtual_accounts`
	return sql
}

func (virtualAccount VirtualAccount) GetListVirtualAccount() ([]VirtualAccountNoPassword, error) {

	sql := virtualAccount.selectQuery()
	cursor, err := virtualAccount.Query(sql)
	if err != nil {
		return nil, err
	}

	virtualAccounts := []VirtualAccountNoPassword{}

	for cursor.Next() {
		c, err := virtualAccount.fetchRow(cursor)
		if err != nil {
			return nil, err
		}
		virtualAccounts = append(virtualAccounts, c)
	}

	return virtualAccounts, nil
}

func (virtualAccount VirtualAccount) GetSeqNo() (*SeqNo, error) {
	sql := `select count(virtual_account_no)as seq_no from virtual_accounts`
	cursor, err := virtualAccount.Query(sql)
	if err != nil {
		return nil, err
	}
	if cursor.Next() {
		c := SeqNo{}
		err := cursor.Scan(&c.SeqNo)
		if err != nil {
			return nil, err
		}
		return &c, nil
	}
	return nil, nil
}

func (virtualAccount VirtualAccount) FindVirtualAccountByNo(virtualAccountNo string) (*VirtualAccount, error) {

	if helpers.Empty(virtualAccountNo) {
		return nil, errors.New("invalid_virtual_account")
	}

	sql := virtualAccount.selectQueryAll()
	sql += ` where virtual_account_no=$1`

	cursor, err := virtualAccount.Query(sql, virtualAccountNo)
	if err != nil {
		return nil, err
	}

	if cursor.Next() {
		c, err := virtualAccount.fetchRowAll(cursor)
		if err != nil {
			return nil, err
		}
		return &c, nil
	}

	return nil, errors.New("VirtualAccount tidak ditemukan")

}

func (virtualAccount VirtualAccount) FindVirtualAccountByPhone(phoneNumber string) (*VirtualAccount, error) {

	if helpers.Empty(phoneNumber) {
		return nil, errors.New("invalid_virtual_account")
	}

	sql := `select virtual_account_no, password from virtual_accounts`
	sql += ` where phone_number=$1`

	cursor, err := virtualAccount.Query(sql, phoneNumber)
	if err != nil {
		return nil, err
	}

	if cursor.Next() {
		c, err := virtualAccount.fetchRowFind(cursor)
		if err != nil {
			return nil, err
		}
		return &c, nil
	}

	return nil, errors.New("VirtualAccount tidak ditemukan")

}

func (virtualAccount VirtualAccount) RemoveByNo(virtualAccountNo string) error {

	sql := "delete from virtual_accounts where virtual_account_no=$1"
	result, err := virtualAccount.Exec(sql, virtualAccountNo)
	if err != nil {
		return err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	log.Println("VirtualAccount.Remove.AffectedRows=", affectedRows)
	return nil
}

func (virtualAccount VirtualAccount) Update(virtualAccountNo string, accountName string) error {

	sql := `update virtual_accounts 
	set account_name=$2,
	updated_at=$3

	where virtual_account_no=$1
	`
	result, err := virtualAccount.Exec(sql, virtualAccountNo, accountName, time.Now().UTC())
	if err != nil {
		return err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	log.Println("VirtualAccount.Update.AffectedRows=", affectedRows)

	return nil

}

func (virtualAccount VirtualAccount) UpdatePin(virtualAccountNo string, pin string) error {

	sql := `update virtual_accounts 
	set pin=$2,
	updated_at=$3

	where virtual_account_no=$1
	`

	pin = base64.StdEncoding.EncodeToString([]byte(pin))
	result, err := virtualAccount.Exec(sql, virtualAccountNo, pin, time.Now().UTC())
	if err != nil {
		return err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	log.Println("VirtualAccount.Update.AffectedRows=", affectedRows)

	return nil

}

func (virtualAccount VirtualAccount) UpdateSaldo(virtualAccountNo string, saldo float64) error {

	sql := `update virtual_accounts 
	set saldo=$2,
	updated_at=$3

	where virtual_account_no=$1
	`
	result, err := virtualAccount.Exec(sql, virtualAccountNo, saldo, time.Now().UTC())
	if err != nil {
		return err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	log.Println("VirtualAccount.Update.AffectedRows=", affectedRows)

	return nil

}

func (virtualAccount VirtualAccount) Truncate() error {

	sql := `delete from virtual_accounts`
	result, err := virtualAccount.Exec(sql)
	if err != nil {
		return err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	log.Println("VirtualAccount.Truncate.AffectedRows=", affectedRows)
	return nil

}

func (virtualAccount VirtualAccount) Migrate() error {

	sqlDropTable := "DROP TABLE IF EXISTS public.virtual_accounts CASCADE"

	sqlCreateTable := `
CREATE TABLE public.virtual_accounts
(
    virtual_account_no character varying(40) NOT NULL,
    phone_number character varying(20) UNIQUE NOT NULL,
    email character varying(20) NOT NULL,
	account_name character varying(60) NOT NULL,
	pin character varying(60) NOT NULL,
	password character varying(60) NOT NULL,
	saldo double precision NOT NULL,
	seq_no int NOT NULL, 
    created_at timestamp with time zone NOT NULL,
	updated_at timestamp with time zone,
    CONSTRAINT virtual_accounts_pkey PRIMARY KEY (virtual_account_no)
)

TABLESPACE pg_default;

ALTER TABLE public.virtual_accounts
    OWNER to postgres;
	`

	_, err := virtualAccount.DB.Exec(sqlDropTable)
	if err != nil {
		return err
	}

	_, err = virtualAccount.DB.Exec(sqlCreateTable)
	if err != nil {
		return err
	}

	return nil

}
