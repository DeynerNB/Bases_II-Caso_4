package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/Shopify/sarama"
)

func main() {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	// Especificar el grupo de consumidores y el topic al que se suscribirá el consumidor.
	consumer, err := sarama.NewConsumer([]string{"10.0.1.3:9092"}, config)
	if err != nil {
		log.Fatalln("Error al crear el consumidor de Kafka:", err)
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatalln("Error al cerrar el consumidor de Kafka:", err)
		}
	}()

	// Suscribirse al topic de Kafka.
	topic := "bases2_marvel"
	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalln("Error al suscribirse al topic de Kafka:", err)
	}
	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			log.Fatalln("Error al cerrar el consumidor de partición de Kafka:", err)
		}
	}()

	// Configurar una señal para manejar el cierre del programa.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// Consumir mensajes de Kafka hasta que se reciba una señal para detenerse.
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			fmt.Printf("Mensaje recibido en el topic %s, partición %d, offset %d: %s = %s\n", msg.Topic, msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
		case err := <-partitionConsumer.Errors():
			fmt.Println("Error al consumir mensaje de Kafka:", err)
		case <-signals:
			fmt.Println("Recibida señal para detener el consumo de mensajes de Kafka.")
			return
		}
	}
}
