package infrastructureR

import (
	"fmt"
	"log"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/joho/godotenv"
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

	if topic == "" {
		return fmt.Errorf("el tópico MQTT está vacío")
	}

	token := client.client.Publish(topic, 0, false, message)
	token.Wait()

	if token.Error() != nil {
		return fmt.Errorf("error al publicar mensaje en MQTT: %w", token.Error())
	}
	return nil
}
