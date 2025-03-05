package infrastructure

import (
	"encoding/json"
	"fmt"
	"minimulti/src/events/application"
	"net/http"
)

type EventController struct {
	CreateUseCase          *application.CreateEvent
	GetAllUseCase          *application.GetAllEvents
	DeleteAllEventsUseCase *application.DeletEvents
}

func NewEventController(
	create *application.CreateEvent,
	getAll *application.GetAllEvents,
	deleteAll *application.DeletEvents,
) *EventController {
	return &EventController{
		CreateUseCase:          create,
		GetAllUseCase:          getAll,
		DeleteAllEventsUseCase: deleteAll,
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
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("Evento registrado desde: '%s'", espInput.Device_name)))
}

func (ec *EventController) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	events, err := ec.GetAllUseCase.Run()

	if err != nil {
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
