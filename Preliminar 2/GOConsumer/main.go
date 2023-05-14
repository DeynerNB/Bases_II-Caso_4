package main

import (
	"context"
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
}

func configurarKafkaConsumer() (*kafka.Consumer, error) {
	//Se define la configuracion del consumer
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "10.0.0.26:9092",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})
	return consumer, err
}

func conectarMongoDB() (*mongo.Client, error) {
	//Se define la conexion al router con las bases de datos
	mongoURI := "mongodb://192.168.192.14:27021"
	client, err := mongo.Connect(nil, options.Client().ApplyURI(mongoURI))
	return client, err
}

func dataInsertion(collection *mongo.Collection, company Service) (*mongo.InsertOneResult, error) {
	res, err := collection.InsertOne(nil, bson.M{
		"name":            company.Name,
		"logo":            company.Logo,
		"category":        company.Category,
		"description":     company.Description,
		"score":           company.Score,
		"bussiness_model": company.BusinessModel,
		"access_info":     company.AccessInfo,
		"country":         company.Country,
	})
	return res, err
}

func main() {
	topic_addService := "nuevo_Servicio"
	topic_updateService := "actualizar_Servicio"

	consumer, err := configurarKafkaConsumer()
	if err != nil {
		log.Fatalf("Failed to create consumer: %s\n", err)
	}
	defer consumer.Close()

	fmt.Println("Consumer Iniciado...")
	err = consumer.SubscribeTopics([]string{"nuevo_Servicio", "actualizar_Servicio"}, nil)
	if err != nil {
		log.Fatalf("Failed to subscribe to topic: %s\n", err)
	}

	client, err := conectarMongoDB()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %s\n", err)
	}
	defer client.Disconnect(nil)

	fmt.Println("Conexion exitosa...")

	collection := client.Database("AlBulbDB").Collection("services")

	for {
		// Leer mensajes de Kafka
		msg, err := consumer.ReadMessage(-1)
		if err != nil {
			log.Printf("Failed to read message: %s\n", err)
			continue
		}

		if *msg.TopicPartition.Topic == topic_addService {
			// Decodificar el mensaje JSON
			var company Service
			err = json.Unmarshal(msg.Value, &company)
			if err != nil {
				log.Printf("Failed to decode JSON: %s\n", err)
				continue
			}

			// Insertar datos en MongoDB
			res, err := dataInsertion(collection, company)
			if err != nil {
				log.Printf("Failed to insert data into MongoDB: %s\n", err)
				continue
				break
			}
			fmt.Printf("Inserted document with ID: %v\n", res.InsertedID)

		} else if *msg.TopicPartition.Topic == topic_updateService {
			//Crear un mapa vacío para almacenar el JSON
			data := make(map[string]interface{})
			println("Se ha creado el map vacio...")

			// Decodificar el JSON en el mapa
			err := json.Unmarshal(msg.Value, &data)
			if err != nil {
				fmt.Println("Error al decodificar el JSON:", err)
				return
			}

			//Filtramos la el request para saber donde hacer el update
			filter := bson.D{{"name", data["service_name"].(string)}, {"country", data["country"].(string)}}

			//update := bson.D{{"$set", bson.D{data["fields"].(map[string]interface{})}}}

			fieldsMap := data["fields"].(map[string]interface{})

			fieldsBson := bson.D{}

			for key, value := range fieldsMap {
				fieldsBson = append(fieldsBson, bson.E{Key: key, Value: value})
			}
			update := bson.D{{"$set", fieldsBson}}

			res, err := collection.UpdateOne(context.TODO(), filter, update)
			if err != nil {
				panic(err)
			}
			fmt.Printf("Updated document: %v\n", res)
		}
	}
	//consumer.Close()
	//fmt.Printf("\nConsumer cerrado!\n")
}
