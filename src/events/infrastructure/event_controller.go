package infrastructure

import (
	"encoding/json"
	"fmt"
	"minimulti/src/core/rabbit/infrastructureR"
	"minimulti/src/events/application"
	"net/http"
	//"github.com/streadway/amqp"
)

type EventController struct {
	CreateUseCase          *application.CreateEvent
	GetAllUseCase          *application.GetAllEvents
	DeleteAllEventsUseCase *application.DeletEvents
	RabbitClient           *infrastructureR.RabbitMQ
}

func NewEventController(
	create *application.CreateEvent,
	getAll *application.GetAllEvents,
	deleteAll *application.DeletEvents,
	rabbitClient *infrastructureR.RabbitMQ,
) *EventController {
	return &EventController{
		CreateUseCase:          create,
		GetAllUseCase:          getAll,
		DeleteAllEventsUseCase: deleteAll,
		RabbitClient:           rabbitClient,
	}
}

func (ec *EventController) CreateNewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var espInput struct {
		Device_name string `json:"Device_name"`
	}

	err := json.NewDecoder(r.Body).Decode(&espInput)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al leer datos: %v", err), http.StatusBadRequest)
		return
	}
	fmt.Printf("Datos recibidos: %v\n", espInput)

	err = ec.CreateUseCase.Run(espInput.Device_name)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al registrar evento: %v", err), http.StatusInternalServerError)
		return
	}

	eventNotification := map[string]interface{}{
		"device_name": espInput.Device_name,
		"message":     "Evento registrado desde dispositivo",
	}

	eventNotificationJSON, err := json.Marshal(eventNotification)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al convertir mensaje a JSON: %v", err), http.StatusInternalServerError)
		return
	}

	err = ec.RabbitClient.PublishMessage("event_queue", eventNotificationJSON)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al publicar mensaje en RabbitMQ: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("Evento registrado desde: '%s'", espInput.Device_name)))
}

func (ec *EventController) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	events, err := ec.GetAllUseCase.Execute()

	if err != nil {

		if err.Error() == "no existen registros" {
			http.Error(w, "No existen registros", http.StatusNotFound)
			return
		}

		http.Error(w, fmt.Sprintf("Error al obtener eventos: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "Application/json")
	json.NewEncoder(w).Encode(events)
}

func (ec *EventController) DeleteAllHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	err := ec.DeleteAllEventsUseCase.Run()

	if err != nil {
		http.Error(w, fmt.Sprintf("Error al eliminar registro de eventos: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	msg := ("registro de eventos eliminado con exito")
	w.Write([]byte(msg))
}

/*
func (ec *EventController) CreateNewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var espInput struct {
		Device_name string `json:"Device_name"`
	}

	err := json.NewDecoder(r.Body).Decode(&espInput)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al leer datos: %v", err), http.StatusBadRequest)
		return
	}
	fmt.Printf("Datos recibidos: %v\n", espInput)

	err = ec.CreateUseCase.Run(espInput.Device_name)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al registrar evento: %v", err), http.StatusInternalServerError)
		return
	}

	eventNotification := map[string]interface{}{
		"device_name": espInput.Device_name,
		"message":     "Evento registrado desde dispositivo",
	}

	eventNotificationJSON, err := json.Marshal(eventNotification)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al convertir mensaje a JSON: %v", err), http.StatusInternalServerError)
		return
	}

	// 🔹 Se abre un nuevo canal antes de publicar
	ch, err := ec.RabbitClient.conn.Channel()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al abrir canal: %v", err), http.StatusInternalServerError)
		return
	}
	defer ch.Close() // 🔹 Se asegura de cerrar el canal después de usarlo

	// 🔹 Se declara la cola para evitar problemas si no existe
	_, err = ch.QueueDeclare(
		"event_queue",
		true,  // Durable
		false, // Auto-delete
		false, // Exclusivo
		false, // No-wait
		nil,   // Argumentos adicionales
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al declarar la cola: %v", err), http.StatusInternalServerError)
		return
	}

	// 🔹 Se publica el mensaje utilizando el canal recién creado
	err = ch.Publish(
		"",            // Exchange
		"event_queue", // Routing Key
		false,         // Mandatory
		false,         // Immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        eventNotificationJSON,
		},
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al publicar mensaje en RabbitMQ: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("Evento registrado desde: '%s'", espInput.Device_name)))
} */
