package models

import (
	"encoding/base64"
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

type TokenDetails struct {
	GUID uuid.UUID	`json:"GUID,omitempty" bson:"GUID,omitempty"`
	Token    string `json:"token,omitempty" bson:"token,omitempty"`
	Refresh string `json:"refresh,omitempty" bson:"refresh,omitempty"`
}

func (tkn *TokenDetails) CreateTokens(GUID uuid.UUID) *TokenDetails {
	token := jwt.New(jwt.SigningMethodHS512)
	claims := token.Claims.(jwt.MapClaims)
	claims["UUID"] = GUID
	claims["exp"] = time.Now().Add(time.Minute * 5).Unix()
	tokenString,_:= token.SignedString([]byte("secret"))
	tkn.Token=tokenString
	refreshToken := jwt.New(jwt.SigningMethodHS512)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["UUID"] = GUID
	rtClaims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	refreshTokenString,_:= refreshToken.SignedString([]byte("secret"))
	tkn.Refresh=base64.StdEncoding.EncodeToString([]byte(refreshTokenString))
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


	table.Client().UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		_, err := table.Collection("token").UpdateMany(
			ctx,
			bson.M{"GUID":GUID},
			bson.D{
				{"$set",bson.D{{"refresh",tkn.Refresh}}},
				{"$set",bson.D{{"token",tkn.Token}}}})
		if err != nil {
			sessionContext.AbortTransaction(sessionContext)
			fmt.Println(err)
			return err
		}
		return nil
	})

	return tkn}


func (tkn *TokenDetails) DeleteRefreshToken(refreshToken string) {
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

	table.Client().UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		_, err := table.Collection("token").DeleteOne(ctx, bson.M{"refresh": refreshToken})
		if err != nil {
			sessionContext.AbortTransaction(sessionContext)
			fmt.Println(err)
			return err
		}
		return nil
	})
}


func (tkn *TokenDetails) DeleteAllRefreshTokenByGUID(GUID uuid.UUID) {	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017/?readPreference=primary&appname=MongoDB%20Compass&ssl=false&retryWrites=false"))
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

	table.Client().UseSession(ctx,func(sessionContext mongo.SessionContext) error {
		_, err := table.Collection("token").DeleteMany(ctx, bson.M{"GUID": GUID})
		if err != nil {
			sessionContext.AbortTransaction(sessionContext)
			fmt.Println(err)
			return err
		}
		return nil
	})
}

func (tkn *TokenDetails) refreshingToken(token string) *TokenDetails {
	stringGUID :=tkn.ParseToken(token)
	GUID, _ :=uuid.Parse(stringGUID)
	return tkn.CreateTokens(GUID)
}


func(tkn *TokenDetails) ParseToken(token string) string {
	tokenString := token
	claims := jwt.MapClaims{}
	 jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})

	for key, val := range claims {
		if key == "UUID"{
			str:=fmt.Sprint(val)
			return str
		}
	}
	return "not parse"
}

