package infrastructure

import (
	"database/sql"
	"fmt"
	evententity "minimulti/src/events/domain/event_entity"
	"time"
)

type EventRepository struct {
	db *sql.DB
}

func NewEventRepository(db *sql.DB) *EventRepository {
	return &EventRepository{db: db}
}

func (repo *EventRepository) Create(event *evententity.Event) error {
	query := "INSERT INTO pir_events (device_name) VALUES(?)"
	fmt.Printf("Evento creado desde dispositivo: %s \n", event.Device_name)
	_, err := repo.db.Exec(query, event.Device_name)
	if err != nil {
		return fmt.Errorf("error: %s", err)
	}
	return nil
}

func (repo *EventRepository) GetAll() ([]*evententity.Event, error) {
	query := "SELECT * FROM pir_events"
	rows, err := repo.db.Query(query)

	if err != nil {
		return nil, fmt.Errorf("error al obtener datos %w", err)
	}
	defer rows.Close()

	var events []*evententity.Event

	for rows.Next() {
		var event evententity.Event
		var createdAtRaw []uint8

		if err := rows.Scan(&event.Id, &event.Device_name, &createdAtRaw); err != nil {
			return nil, fmt.Errorf("error al obtener datos: %w", err)
		}

		createdAtStr := string(createdAtRaw)
		parsedTime, err := time.Parse("2006-01-02 15:04:05", createdAtStr)

		if err != nil {
			return nil, fmt.Errorf("error al parsear created_at: %w", err)
		}

		event.Created_at = parsedTime

		events = append(events, &event)
	}

	if len(events) == 0 {
		return nil, fmt.Errorf("no existen registros")
	}
	return events, nil
}

func (repo *EventRepository) DeleteAll() error {
	query := "DELETE FROM pir_events"
	_, err := repo.db.Exec(query)
	return err
}
