package infrastructureR

import (
	"encoding/json"
	"fmt"
	"minimulti/src/core/rabbit/applicationR"
	"net/http"
)

type MessageController struct {
	messageService *applicationR.MessageService
}

func NewMessageController(service *applicationR.MessageService) *MessageController {
	return &MessageController{messageService: service}
}

func (controller *MessageController) PublishMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var msg struct {
		Content string `json:"content"`
	}

	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al decodificar el mensaje: %v", err), http.StatusBadRequest)
		return
	}

	err = controller.messageService.PublishMessage("api_queue", msg.Content)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al publicar el mensaje: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Mensaje publicado con éxito"))
}
