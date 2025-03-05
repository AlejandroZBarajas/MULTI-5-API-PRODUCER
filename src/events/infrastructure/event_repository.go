package infrastructure

import (
	"database/sql"
	"fmt"
	evententity "minimulti/src/events/domain/event_entity"
)

type EventRepository struct {
	db *sql.DB
}

func NewEventRepository(db *sql.DB) *EventRepository {
	return &EventRepository{db: db}
}

func (repo *EventRepository) Save(event *evententity.Event) error {
	query := "INSERT INTO pir_events (device_name) VALUES(?)"
	fmt.Printf("Evento creado desde dispositivo: %s \n", event.Device_name)
	_, err := repo.db.Exec(query, event.Device_name)
	if err != nil {
		return fmt.Errorf("error: %s", err)
	}
	return nil
}

func (repo *EventRepository) GetAllEvents() ([]*evententity.Event, error) {
	query := "SELECT * FROM pir_events"
	rows, err := repo.db.Query(query)

	if err != nil {
		return nil, fmt.Errorf("error al insertar datos %w", err)
	}
	defer rows.Close()

	var events []*evententity.Event

	for rows.Next() {
		var event evententity.Event

		if err := rows.Scan(&event.Id, &event.Device_name, &event.Created_at); err != nil {
			return nil, fmt.Errorf("error al obtener datos: %w", err)
		}
		events = append(events, &event)
	}
	return events, nil
}

func (repo *EventRepository) DeleteAllEvents() error {
	query := "DELETE FROM pir_events"
	_, err := repo.db.Exec(query)
	return err
}
