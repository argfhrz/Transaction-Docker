package data

import (
	"log"
	"testing"
	"virtual-account/config"
	"virtual-account/connection"
	"virtual-account/helpers"
)

func TestBankAccount(t *testing.T) {

	//test connection
	db := connection.OpenConnection(config.DEV)
	defer db.Close()

	bankAccountModel := CreateBankAccount(db)
	err := bankAccountModel.Migrate()
	if err != nil {
		t.Fatal(err.Error())
	}
	err = bankAccountModel.Truncate()
	if err != nil {
		t.Fatal(err.Error())
	}

	bankAccountOwner := "Owner1"
	saldo := float64(50000000)

	//add
	bankAccountNo, err := bankAccountModel.AddBankAccount(bankAccountOwner, saldo)
	if err != nil {
		t.Fatal(err.Error())
	}

	//get
	bankAccounts, err := bankAccountModel.GetListBankAccount()
	if err != nil {
		t.Fatal(err)
	}
	log.Println("Inserted", helpers.ToJson(bankAccounts))

	//update
	// saldo = float64(3000000)

	// err = bankAccountModel.Update(bankAccountNo, bankAccountOwner, saldo)
	// if err != nil {
	// 	t.Fatal(err.Error())
	// }

	// bankAccount, err := bankAccountModel.GetListBankAccount()
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// log.Println("Updated", helpers.ToJson(bankAccount))

	//find by no
	found, err := bankAccountModel.FindBankAccountByNo(bankAccountNo)
	if err != nil {
		t.Fatal(err)
	}
	if found == nil {
		t.Fatal("bank_account_not_found")
	}

	if found.BankAccountNo != bankAccountNo {
		t.Fatal("Expected=", bankAccountNo, "actual=", found.BankAccountNo)
	}

	if found.BankAccountOwner != bankAccountOwner {
		t.Fatal("Expected=", bankAccountOwner, "actual=", found.BankAccountOwner)
	}

	if found.Saldo != saldo {
		t.Fatal("Expected=", saldo, "actual=", found.Saldo)
	}

	// remove by no
	// err = bankAccountModel.RemoveByNo(bankAccountNo)
	// if err != nil {
	// 	t.Fatal(err.Error())
	// }

	// bankAccounts, err = bankAccountModel.GetListBankAccount()
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// log.Println("After_Deleted", helpers.ToJson(bankAccounts))

}
