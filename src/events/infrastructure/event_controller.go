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
	MQTTClient             *infrastructureR.MQTTClient
}

func NewEventController(
	create *application.CreateEvent,
	getAll *application.GetAllEvents,
	deleteAll *application.DeletEvents,
	mqttClient *infrastructureR.MQTTClient,
) *EventController {
	return &EventController{
		CreateUseCase:          create,
		GetAllUseCase:          getAll,
		DeleteAllEventsUseCase: deleteAll,
		MQTTClient:             mqttClient,
	}
}

func (ec *EventController) CreateNewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var espInput struct {
		Serie       string  `json:"serie"`
		Title       string  `json:"title"`
		Description float32 `json:"description"`
		Emitter     string  `json:"emitter"`
		Topic       string  `json:"topic"`
	}

	err := json.NewDecoder(r.Body).Decode(&espInput)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al leer datos: %v", err), http.StatusBadRequest)
		return
	}
	fmt.Printf("Datos recibidos en controller: %v\n", espInput)

	eventID, createdAt, err := ec.CreateUseCase.Run(espInput.Serie, espInput.Title, espInput.Description, espInput.Emitter, espInput.Topic)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al registrar evento: %v", err), http.StatusInternalServerError)
		return
	}

	eventNotification := map[string]interface{}{
		"id":          eventID,
		"createdAt":   createdAt,
		"serie":       espInput.Serie,
		"title":       espInput.Title,
		"description": espInput.Description,
		"emitter":     espInput.Emitter,
		"topic":       espInput.Topic,
	}

	eventNotificationJSON, err := json.Marshal(eventNotification)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al convertir mensaje a JSON: %v", err), http.StatusInternalServerError)
		return
	}

	err = ec.MQTTClient.PublishMessage(espInput.Topic, eventNotificationJSON)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al publicar mensaje en RabbitMQ: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("%s : %s .Registrado desde: '%s' con ID: %d y TOPIC: %s ", espInput.Title, espInput.Description, espInput.Serie, eventID, espInput.Topic)))
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
