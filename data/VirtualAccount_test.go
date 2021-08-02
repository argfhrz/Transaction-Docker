package data

import (
	"log"
	"testing"
	"virtual-account/config"
	"virtual-account/connection"
)

func TestVirtualAccount(t *testing.T) {

	//test connection
	db := connection.OpenConnection(config.DEV)
	defer db.Close()

	virtualAccountModel := CreateVirtualAccount(db)
	err := virtualAccountModel.Migrate()
	if err != nil {
		t.Fatal(err.Error())
	}
	err = virtualAccountModel.Truncate()
	if err != nil {
		t.Fatal(err.Error())
	}

	phoneNumber := "0812312312"
	email := "email@gmail.com"
	accountName := "account1"
	pin := "123456"
	password := "123"

	//add
	virtualAccountNo, err := virtualAccountModel.Add(phoneNumber, email, accountName, pin, password)
	if err != nil {
		t.Fatal(err.Error())
	}
	log.Println(virtualAccountNo)

	phoneNumber2 := "0812312313"
	email2 := "email2@gmail.com"
	accountName2 := "account2"
	pin2 := "654321"
	password2 := "123"

	//add
	virtualAccountNo2, err := virtualAccountModel.Add(phoneNumber2, email2, accountName2, pin2, password2)
	if err != nil {
		t.Fatal(err.Error())
	}
	log.Println(virtualAccountNo2)

	phoneNumber3 := "0812234567"
	email3 := "email3@gmail.com"
	accountName3 := "merchant1"
	pin3 := "98765"
	password3 := "123"

	//add
	virtualAccountNo3, err := virtualAccountModel.Add(phoneNumber3, email3, accountName3, pin3, password3)
	if err != nil {
		t.Fatal(err.Error())
	}
	log.Println(virtualAccountNo3)

	//get
	// virtualAccounts, err := virtualAccountModel.GetListVirtualAccount()
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// log.Println("Inserted", helpers.ToJson(virtualAccounts))

	//update
	// accountName = "account2"
	// pin = "654321"

	// err = virtualAccountModel.Update(virtualAccountNo, accountName, pin)
	// if err != nil {
	// 	t.Fatal(err.Error())
	// }

	// virtualAccount, err := virtualAccountModel.GetListVirtualAccount()
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// log.Println("Updated", helpers.ToJson(virtualAccount))

	// //find by no
	// found, err := virtualAccountModel.FindVirtualAccountByPhone(phoneNumber)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// if found == nil {
	// 	t.Fatal("virtual_account_not_found")
	// }

	// if found.PhoneNumber != phoneNumber {
	// 	t.Fatal("Expected=", phoneNumber, "actual=", found.PhoneNumber)
	// }

	// if found.Email != email {
	// 	t.Fatal("Expected=", email, "actual=", found.Email)
	// }

	// if found.AccountName != accountName {
	// 	t.Fatal("Expected=", accountName, "actual=", found.AccountName)
	// }

	// if found.Pin != pin {
	// 	t.Fatal("Expected=", pin, "actual=", found.Pin)
	// }

	// if found.Password != password {
	// 	t.Fatal("Expected=", password, "actual=", found.Password)
	// }

	// if found.Saldo != saldo {
	// 	t.Fatal("Expected=", saldo, "actual=", found.Saldo)
	// }

	// if found.SeqNo != seqNo {
	// 	t.Fatal("Expected=", seqNo, "actual=", found.SeqNo)
	// }

	// // remove by no
	// err = virtualAccountModel.RemoveByNo(virtualAccountNo)
	// if err != nil {
	// 	t.Fatal(err.Error())
	// }

	// virtualAccounts, err = virtualAccountModel.GetListVirtualAccount()
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// log.Println("After_Deleted", helpers.ToJson(virtualAccounts))

}
