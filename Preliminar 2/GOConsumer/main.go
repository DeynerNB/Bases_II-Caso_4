package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Estructuras para almacenar datos JSON
type BusinessModel struct {
	Description string `json:"description"`
	Type        string `json:"type"`
}

type SocialNetwork struct {
	Instagram string `json:"instagram"`
	Facebook  string `json:"facebook"`
}

type AccessInfo struct {
	URL           string        `json:"url"`
	Phone         string        `json:"phone"`
	SocialNetwork SocialNetwork `json:"social_network"`
}

type Service struct {
	Name          string          `json:"name"`
	Logo          string          `json:"logo"`
	Category      string          `json:"category"`
	Description   string          `json:"description"`
	Score         float64         `json:"score"`
	BusinessModel []BusinessModel `json:"bussiness_model"`
	AccessInfo    AccessInfo      `json:"access_info"`
	Country       string          `json:"country"`
	OperationType string          `json:"operation_type"`
}


func main() {

	// Conectar a MongoDB
	mongoURI := "mongodb://192.168.192.57:27040"
	client, err := mongo.Connect(nil, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %s\n", err)
	}
	fmt.Println("Succesfull MongoDB connections!")

	defer client.Disconnect(nil)

	collection := client.Database("PruebaDB").Collection("services")


	jsonData := []byte(`{
        "service_name": "Lampara",
    	"country": "Costa Rica",
    	"fields" : {
        	"logo": "Si actualizo",
        	"category": "Si actualizo"
    	}
    }`)

    // Crear un mapa vacío para almacenar el JSON
    data := make(map[string]interface{})

    // Decodificar el JSON en el mapa
    error1 := json.Unmarshal(jsonData, &data)
    if error1 != nil {
        fmt.Println("Error al decodificar el JSON:", error1)
        return
    }

    fmt.Println(data["service_name"].(string))

	filter := bson.D{{"name", data["service_name"].(string)}}
	update := bson.D{{"$set", data["fields"] }}

	res, error2 := collection.UpdateOne(context.TODO(), filter, update)
	if error2 != nil {
		panic(error2)
	}
	fmt.Printf("Updated document: %v\n", res.UpsertedID)
}
