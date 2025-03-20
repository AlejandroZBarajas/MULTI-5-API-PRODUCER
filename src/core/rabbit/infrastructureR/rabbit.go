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
	exchange   string
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
	exchangeName := os.Getenv("EXCHANGE_NAME")

	if rabbitURL == "" || exchangeName == "" {
		return nil, fmt.Errorf("variables de entorno RABBITMQ indefinidas")
	}

	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		return nil, fmt.Errorf("error al conectar con RabbitMQ: %w", err)
	}
	fmt.Println("conectado a rabbit")

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("error al abrir el canal de RabbitMQ: %w", err)
	}

	err = ch.ExchangeDeclare(
		exchangeName, // Nombre del exchange
		"topic",      // Tipo de exchange
		true,         // Durable
		false,        // Auto-delete
		false,        // Internal
		false,        // NoWait
		nil,          // Argumentos
	)
	/* _, err = ch.QueueDeclare(
		queueName, // Nombre de la cola
		true,      // Durable
		false,     // Auto-delete
		false,     // Exclusive
		false,     // NoWait
		nil,       // Argumentos adicionales
	) */
	if err != nil {
		return nil, fmt.Errorf("error al declarar la cola: %w", err)
	}

	return &RabbitMQ{
		connection: conn,
		channel:    ch,
		exchange:   exchangeName,
	}, nil
}

func (client *RabbitMQ) PublishMessage(routingKey string, message []byte) error {
	err := client.channel.Publish(
		client.exchange, // Exchange
		routingKey,
		false, // Mandatory
		false, // Immediate
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

func (client *RabbitMQ) DeclareQueue(queueName, routingKey string) error {

	_, err := client.channel.QueueDeclare(
		queueName, // Nombre de la cola
		true,      // Durable
		false,     // Auto-delete
		false,     // Exclusive
		false,     // NoWait
		nil,       // Argumentos adicionales
	)
	if err != nil {
		return fmt.Errorf("error al declarar la cola %s: %w", queueName, err)
	}

	err = client.channel.QueueBind(
		queueName,       // Nombre de la cola
		routingKey,      // Routing key (ej: "notificaciones.*")
		client.exchange, // Exchange
		false,           // NoWait
		nil,             // Argumentos
	)
	if err != nil {
		return fmt.Errorf("error al enlazar la cola %s con el routing key %s: %w", queueName, routingKey, err)
	}

	fmt.Printf("✅ Cola '%s' enlazada a exchange '%s' con routing key '%s'\n", queueName, client.exchange, routingKey)
	return nil
}
