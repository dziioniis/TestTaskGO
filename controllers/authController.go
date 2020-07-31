package controllers

import (
	"encoding/json"
	"fmt"
	_ "go.mongodb.org/mongo-driver/bson"
	"net/http"
	"project/models"
	"project/utils"
)


var CreateAccount = func(w http.ResponseWriter, r *http.Request) {
	account := &models.Account{}
	json.NewDecoder(r.Body).Decode(account)
	account.Create()
}
var HelloWorld = func(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode("hello")
}

var CreateTokens= func(w http.ResponseWriter, r *http.Request){
tkn := &models.TokenDetails{}
json.NewDecoder(r.Body).Decode(tkn)
resp:=tkn.CreateTokens(tkn.GUID)
	utils.Respond(w,resp)
}

var DeleteRefreshToken = func(w http.ResponseWriter, r *http.Request) {
	tkn := &models.TokenDetails{}
	json.NewDecoder(r.Body).Decode(tkn)
	tkn.DeleteRefreshToken(tkn.Refresh)
}


var DeleteAllRefreshTokenById = func(w http.ResponseWriter, r *http.Request) {
	tkn := &models.TokenDetails{}
	json.NewDecoder(r.Body).Decode(tkn)
	fmt.Println(tkn.GUID)
	tkn.DeleteAllRefreshTokenByGUID(tkn.GUID)
}

var RefreshingToken = func(w http.ResponseWriter, r *http.Request) {
	tkn := &models.TokenDetails{}
	json.NewDecoder(r.Body).Decode(tkn)
	tkn.DeleteRefreshToken(tkn.Refresh)
}


var ParseToken = func(w http.ResponseWriter, r *http.Request) {
	tkn := &models.TokenDetails{}
	json.NewDecoder(r.Body).Decode(tkn)
	tkn.ParseToken(tkn.Refresh)
}


