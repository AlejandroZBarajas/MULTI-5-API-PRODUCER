package main

import (
	"fmt"
	"log"
	"minimulti/src/core/mysql/infrastructureC"
	"minimulti/src/core/rabbit/infrastructureR"
	"minimulti/src/events/application"
	"minimulti/src/events/infrastructure"
	"net/http"
)

func main() {
	infrastructureC.ConnectDB()
	db := infrastructureC.GetDB()

	eventRepo := infrastructure.NewEventRepository(db)

	createUseCase := application.NewCreateEvent(eventRepo)
	getAllEventsUseCase := application.NewGetAllEvents(eventRepo)
	deleteAllEvents := application.NewDeletEvents(eventRepo)

	rabbitClient, err := infrastructureR.NewMQTTClient()
	if err != nil {
		log.Fatalf("Error al conectar a RabbitMQ: %v", err)
	}
	//defer rabbitClient.Close()

	eventController := infrastructure.NewEventController(
		createUseCase,
		getAllEventsUseCase,
		deleteAllEvents,
		rabbitClient,
	)

	infrastructureC.SetRoutes(eventController)

	fmt.Println("Servidor corriendo en el puerto 8080...")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
