package connection

import (
	"testing"
	"virtual-account/config"
)

func TestConnection_OpenConnection(t *testing.T) {

	db := OpenConnection(config.DEV)
	defer db.Close()

	err := db.Ping()
	if err != nil {
		t.Fatal("ping error", err.Error())
	}

}
