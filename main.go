package main

import (
	"fmt"
	"log"
	"minimulti/src/core/infrastructureC"
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

	eventController := infrastructure.NewEventController(
		createUseCase,
		getAllEventsUseCase,
		deleteAllEvents,
	)

	infrastructureC.SetRoutes(eventController)

	fmt.Println("Servidor corriendo en el puerto 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
