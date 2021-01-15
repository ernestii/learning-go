// From tutorial on https://www.youtube.com/watch?v=oW7PMHEYiSk

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/ernestii/learning-go/1-first-api-with-mongo/controllers"
	"github.com/ernestii/learning-go/1-first-api-with-mongo/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var peopleCollection *mongo.Collection

var personRepo models.PersonRepository


func CreatePersonEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("context-type", "application/json") 
	var person models.Person
	json.NewDecoder(request.Body).Decode(&person)
	result, err := personRepo.CreateOne(&person)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(result)
}

func GetPersonEndpoint(response http.ResponseWriter, request *http.Request) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	response.Header().Add("context-type", "application/json")
	params := mux.Vars(request)

	log.Print("get person - ", params["id"])
	id, _ := primitive.ObjectIDFromHex(params["id"])

	var person models.Person
	err := peopleCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&person)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	json.NewEncoder(response).Encode(person)
}

func GetPeopleWithName(response http.ResponseWriter, request *http.Request) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	response.Header().Add("context-type", "application/json")

	nameQuery := request.URL.Query().Get("name")

	var results []models.Person
	filter := bson.M{
		"firstname": primitive.Regex{
			Pattern: nameQuery,
			Options: "i",
		},
	}
	cursor, _ := peopleCollection.Find(ctx, filter)
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var person models.Person
		cursor.Decode(&person)
		results = append(results, person)
	}

	json.NewEncoder(response).Encode(results)
}

func main() {
	port := 8081
	fmt.Println("Starting the application on port " + strconv.Itoa(port) + "...")

	clientOptions := options.Client().ApplyURI("mongodb://ernest:secret@localhost:27017/mydb")
	client, _ = mongo.NewClient(clientOptions)

	// Init
	personRepo = models.NewPersonRepository(client)
	personController := controllers.NewPersonController(&personRepo)

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	err := client.Connect(ctx) 
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	peopleCollection = client.Database("mydb").Collection("people")

	router := mux.NewRouter()
	router.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("am ok!"))
	})
	router.HandleFunc("/person", CreatePersonEndpoint).Methods("POST")
	router.HandleFunc("/person/{id}", GetPersonEndpoint).Methods("GET")
	router.HandleFunc("/query", GetPeopleWithName).Methods("GET")
	router.HandleFunc("/people", personController.GetPeopleEndpoint).Methods("GET")

	http.ListenAndServe(":" + strconv.Itoa(port), router)
}
