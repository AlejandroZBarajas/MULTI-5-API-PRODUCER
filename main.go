package main

import (
	"minimulti/src/core/infrastructureC"
	"minimulti/src/events/infrastructure"
)

func main() {
	infrastructureC.ConnectDB()
	db := infrastructureC.GetDB()
	EventRepo := infrastructure.NewEventRepository(db)

	infrastructureC.SetRoutes(EventController)

}
