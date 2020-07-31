package models

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"golang.org/x/crypto/bcrypt"
	"log"

	"time"
)



type Account struct {
	GUID  uuid.UUID	`json:"GUID,omitempty" bson:"GUID,omitempty"`
	Name 	string `json:"username,omitempty" bson:"username,omitempty"`
	Password string `json:"password,omitempty" bson:"password,omitempty"`
	Token    string `json:"token,omitempty" bson:"token,omitempty"`
	Refresh string `json:"refresh,omitempty" bson:"refresh,omitempty"`

}



func (account *Account) Create() (map[string]interface{}) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017/?readPreference=primary&appname=MongoDB%20Compass&ssl=false&retryWrites=false"))
	ctx,_:= context.WithTimeout(context.Background(),10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx,readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	table:=client.Database("user")
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	account.Password = string(hashedPassword)
	account.GUID= uuid.New()
	fmt.Println(account.GUID)
	table.Client().UseSession(ctx,func(sessionContext mongo.SessionContext) error{
	_,err:=table.Collection("account").InsertOne(sessionContext, bson.M{"GUID":account.GUID,"name":account.Name,
		"password":account.Password})
		if err != nil {
			fmt.Println(err)
			return err
		}
		_, err = table.Collection("token").InsertOne(sessionContext, bson.M{"GUID": account.GUID, "refresh": "", "token": ""})
		if err != nil {
			sessionContext.AbortTransaction(sessionContext)
			fmt.Println(err)
			return err
		} else {
			sessionContext.CommitTransaction(sessionContext)
		}
	return nil
	})
return nil
}





