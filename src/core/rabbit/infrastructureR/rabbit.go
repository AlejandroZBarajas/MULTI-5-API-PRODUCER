package infrastructureR

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	connection *amqp.Connection
	channel    *amqp.Channel
}

type RabbitMQConfig struct {
	URL       string
	QueueName string
}

func NewRabbitMQ() (*RabbitMQ, error) {

	err := godotenv.Load()
	if err != nil {
		log.
			Fatalf("Error al cargar el archivo .env: %v", err)
	}
	rabbitURL := os.Getenv("RABBITMQ_URL")
	queueName := os.Getenv("QUEUE_NAME")

	if rabbitURL == "" || queueName == "" {
		return nil, fmt.Errorf("variables de entorno RABBITMQ indefinidas")
	}

	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		return nil, fmt.Errorf("error al conectar con RabbitMQ: %w", err)
	}
	fmt.Println("conectado a rabbit")

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("error al abrir el canal de RabbitMQ: %w", err)
	}

	_, err = ch.QueueDeclare(
		queueName, // Nombre de la cola
		true,      // Durable
		false,     // Auto-delete
		false,     // Exclusive
		false,     // NoWait
		nil,       // Argumentos adicionales
	)
	if err != nil {
		return nil, fmt.Errorf("error al declarar la cola: %w", err)
	}

	return &RabbitMQ{
		connection: conn,
		channel:    ch,
	}, nil
}

func (client *RabbitMQ) PublishMessage(queueName string, message []byte) error {
	err := client.channel.Publish(
		"",        // Exchange
		queueName, // Routing key (nombre de la cola)
		false,     // Mandatory
		false,     // Immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)
	if err != nil {
		return fmt.Errorf("error al publicar mensaje: %w", err)
	}
	return nil
}

func (client *RabbitMQ) Close() error {
	if err := client.channel.Close(); err != nil {
		return err
	}
	return client.connection.Close()
}
