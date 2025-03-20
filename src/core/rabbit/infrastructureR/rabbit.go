package infrastructureR

import (
	"fmt"
	"log"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/joho/godotenv"
	//"github.com/streadway/amqp"
)

type MQTTClient struct {
	client    mqtt.Client
	topic     string
	brokerURL string
}

type MQTTConfig struct {
	URL   string
	Topic string
}

func NewMQTTClient() (*MQTTClient, error) {

	err := godotenv.Load()
	if err != nil {
		log.
			Fatalf("Error al cargar el archivo .env: %v", err)
	}

	brokerURL := os.Getenv("RABBITMQ_URL")
	topic := os.Getenv("MQTT_TOPIC")

	if brokerURL == "" || topic == "" {
		return nil, fmt.Errorf("variables de entorno RABBITMQ indefinidas")
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(brokerURL)
	opts.SetClientID("api-client")
	opts.SetCleanSession(true)
	opts.SetWill(topic, "API desconectada", 0, false)

	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("error al conectar con el broker MQTT: %w", token.Error())
	}
	fmt.Println("Conectado al broker MQTT")

	return &MQTTClient{
		client:    client,
		topic:     topic,
		brokerURL: brokerURL,
	}, nil
}

func (client *MQTTClient) PublishMessage(topic string, message []byte) error {
	token := client.client.Publish(topic, 0, false, message)
	token.Wait()

	if token.Error() != nil {
		return fmt.Errorf("error al publicar mensaje en MQTT: %w", token.Error())
	}
	return nil
}

/*
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
*/
