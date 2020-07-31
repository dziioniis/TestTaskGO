package main

import (
	"context"
	"encoding/json"
	_ "encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	_"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"net/http"
	"project/controllers"
	"time"
)

	var client *mongo.Client

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/account", controllers.CreateAccount).Methods("POST")
	r.HandleFunc("/delete", controllers.DeleteRefreshToken).Methods("POST")
	r.HandleFunc("/deleteAllRF", controllers.DeleteAllRefreshTokenById).Methods("POST")
	r.HandleFunc("/createTokens", controllers.CreateTokens).Methods("POST")
	r.HandleFunc("/parse", controllers.ParseToken).Methods("POST")
	corsWrapper := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"Content-Type", "Origin", "Accept", "*"},
	})

	http.ListenAndServe(":8080", corsWrapper.Handler(r))
}

type Person struct {
	ID primitive.ObjectID
	username string `json:"username,omitempty" bson:"username,omitempty"`
	password string `json:"password,omitempty" bson:"password,omitempty"`


}

func CreatePerson(response http.ResponseWriter, request *http.Request){
	response.Header().Add("content-type","application-json")
	var person Person
	json.NewDecoder(request.Body).Decode(&person)
	fmt.Println("val"+request.Form.Encode())
	collection := client.Database("user").Collection("person")
	ctx,_:= context.WithTimeout(context.Background(),10*time.Second)
	result,_:= collection.InsertOne(ctx,person)
	fmt.Println(result.InsertedID)
	json.NewEncoder(response).Encode(result)
}

/*func GetPerson(response http.ResponseWriter, request *http.Request){
	response.Header().Add("content-type","application-json")
	var persons []Person
	collection := client.Database("user").Collection("Person")
	ctx,_:= context.WithTimeout(context.Background(),10*time.Second)
	curser,_:= collection.Find(ctx,bson.M{})
	json.NewEncoder(response).Encode(result)
}
*/

var NotImplemented = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("Not Implemented"))
})

var StatusHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("API is up and running"))
})

var ProductsHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
	client,err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017/?readPreference=primary&appname=MongoDB%20Compass&ssl=false"))
	if err!= nil {
		log.Fatal(err)
	}
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
	database :=client.Database("user");
	podcastsCollection:=database.Collection("pesron")

	podcasts, _ := podcastsCollection.InsertOne(ctx,bson.D{{"username","dziioniis"},{"password","qwerty"}})
	payload, _ := json.Marshal(podcastsCollection.FindOne(ctx,bson.M{"username":"dziioniis"}))
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(payload))
	fmt.Println(podcasts)
	fmt.Println(podcastsCollection.FindOne(ctx,bson.M{"username":"dziioniis"}))
})



