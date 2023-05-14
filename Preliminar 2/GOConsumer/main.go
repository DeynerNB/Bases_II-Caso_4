package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
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
	// Configurar el consumidor de Kafka
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "10.0.0.51:9092",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		log.Fatalf("Failed to create consumer: %s\n", err)
	}

	defer consumer.Close()

	fmt.Print("Consumer iniciado...\n")

	// Suscribirse al topic
	err = consumer.SubscribeTopics([]string{"nuevo_Servicio", "actualizar_Servicio"}, nil)
	if err != nil {
		log.Fatalf("Failed to subscribe to topic: %s\n", err)
	}

	// Conectar a MongoDB
	mongoURI := "mongodb://192.168.192.57:27040"
	client, err := mongo.Connect(nil, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %s\n", err)
	}

	defer client.Disconnect(nil)

	collection := client.Database("PruebaDB").Collection("services")

	for {
		// Leer mensajes de Kafka
		msg, err := consumer.ReadMessage(-1)
		if err != nil {
			log.Printf("Failed to read message: %s\n", err)
			continue
		}

		// Decodificar el mensaje JSON
		var company Service
		err = json.Unmarshal(msg.Value, &company)
		if err != nil {
			log.Printf("Failed to decode JSON: %s\n", err)
			continue
		}

		// Insertar datos en MongoDB
		res, err := collection.InsertOne(nil, bson.M{
			"name":            company.Name,
			"logo":            company.Logo,
			"category":        company.Category,
			"description":     company.Description,
			"score":           company.Score,
			"bussiness_model": company.BusinessModel,
			"access_info":     company.AccessInfo,
			"country":         company.Country,
			"operation_type":  company.OperationType,
		})
		if err != nil {
			log.Printf("Failed to insert data into MongoDB: %s\n", err)
			continue
			break
		}

		fmt.Printf("Inserted document with ID: %v\n", res.InsertedID)
	}
	consumer.Close()
	fmt.Printf("\nConsumer cerrado!\n")
}
