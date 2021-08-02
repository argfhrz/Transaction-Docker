package services

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"log"
	"time"
	"virtual-account/data"
	"virtual-account/nosql"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type LoginService struct {
	BaseService
}

func CreateLoginService(mongoClient *mongo.Client, db *sql.DB) LoginService {
	loginService := LoginService{}
	loginService.MongoClient = mongoClient
	loginService.DB = db
	return loginService
}

func (service LoginService) Login(ctx context.Context, phoneNumber string, password string) (string, error) {

	//GET Virtual Account
	virtualAccountModel := data.CreateVirtualAccount(service.DB)
	virtualAccount, err := virtualAccountModel.FindVirtualAccountByPhone(phoneNumber)
	if err != nil {
		log.Println(err)
		return "", errors.New("invalid_login")
	}

	//CEK PASSWORD
	err = bcrypt.CompareHashAndPassword([]byte(virtualAccount.Password), []byte(password))
	if err != nil {
		log.Println(err)
		return "", errors.New("invalid_password")
	}

	//CREATE TOKEN
	id := uuid.New().String()
	now := time.Now().UTC().Format(time.RFC3339)
	token := id + phoneNumber + now

	encryptedToken, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	encodedToken := base64.StdEncoding.EncodeToString(encryptedToken)

	//SAVE LOGIN
	login := nosql.CreateLoginNoSql(service.MongoClient)
	err = login.AddLogin(ctx, encodedToken, virtualAccount.VirtualAccountNo)
	if err != nil {
		return "", err
	}
	return encodedToken, nil

}

func (service LoginService) ParseToken(ctx context.Context, token string) (*nosql.LoginNoSql, error) {

	login := nosql.CreateLoginNoSql(service.MongoClient)
	err := login.FindOneByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	if login.IsExpired {
		return nil, errors.New("login_is_expired")
	}

	return &login, nil

}

func (service LoginService) Logout(ctx context.Context, accountNo string) error {
	logout := nosql.CreateLoginNoSql(service.MongoClient)
	err := logout.UpdateExpired(ctx, accountNo)
	if err != nil {
		return err
	}

	return nil
}
