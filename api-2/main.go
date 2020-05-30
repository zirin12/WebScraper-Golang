package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DBResponseStruct has the record structure as in mongodb
type DBResponseStruct struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	URL       string             `json:"url" bson:"url"`
	CreatedAt string             `json:"created_at" bson:"created_at"`
	Product   ProductInfo        `json:"product" bson:"product"`
	UpdatedAt string             `json:"updated_at" bson:"updated_at"`
}

// ResponseStruct has the final response fields
type ResponseStruct struct {
	URL     string      `json:"url"`
	Product ProductInfo `json:"product"`
}

// ProductInfo has product related info
type ProductInfo struct {
	Title            string `json:"title"`
	ImageURL         string `json:"imageURL"`
	ShortDescription string `json:"description"`
	Rating           string `json:"rating"`
	Price            string `json:"price"`
	TotalReviews     string `json:"totalReviews"`
}

var products *mongo.Collection

func main() {
	// Connect to mongo
	clientOptions := options.Client().ApplyURI("mongodb://mongo:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Connected to MongoDB!")

	// Get products collection
	products = client.Database("go_rest_api").Collection("products")

	r := mux.NewRouter()
	r.HandleFunc("/products", createProduct).Methods("POST")
	r.HandleFunc("/products", getProducts).Methods("GET")
	if err := http.ListenAndServe(":8080", cors.AllowAll().Handler(r)); err != nil {
		log.Fatal(err)
	}
}

// Add the product info received from api-1 to mongoDB
func createProduct(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var product ResponseStruct
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		log.Fatalln(err.Error())
		http.Error(w, "Invalid Request body", http.StatusInternalServerError)
		return
	}

	// create or update the record in the DB
	// Use update with upsert set as true
	// Using bson map filter and update document
	// Find id the product url already exists , if it does update the fields and update time
	// if it doesn't upsert condition will make sure it inserts as a new document with create time
	timestamp := time.Now().UTC().String()
	opts := options.Update().SetUpsert(true)
	filter := bson.M{"url": bson.M{"$eq": product.URL}}
	update := bson.M{
		"$setOnInsert": bson.M{
			"created_at": timestamp,
		},
		"$set": bson.M{
			"url":        product.URL,
			"updated_at": timestamp,
			"product":    product.Product,
		},
	}

	result, err := products.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		log.Fatal(err)
		responseError(w, err.Error(), http.StatusInternalServerError)
	}

	responseJSON(w, result)
}

// Get all the products added to the DB
func getProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	var DBresponse []DBResponseStruct
	cursor, err := products.Find(context.TODO(), bson.M{})
	if err != nil {
		responseError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var DBrecord DBResponseStruct

		err := cursor.Decode(&DBrecord)
		if err != nil {
			log.Fatalln(err)
		}
		DBresponse = append(DBresponse, DBrecord)
	}

	if err := cursor.Err(); err != nil {
		log.Fatalln(err)
	}

	responseJSON(w, DBresponse)

}

func responseError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func responseJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
