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

func (repo *EventRepository) Create(event *evententity.Event) (int, string, error) {
	query := "INSERT INTO notifications (title, description, emitter, topic) VALUES(?, ?, ?, ?)"
	fmt.Printf("%s : %s .Registrado desde: '%s' GUARDADO EN BASE DE DATOS (event repo)", event.Title, event.Description, event.Emitter)

	result, err := repo.db.Exec(query, event.Title, event.Description, event.Emitter, event.Topic)

	if err != nil {
		return 0, "", err
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return 0, "", err
	}

	var createdAt string
	err = repo.db.QueryRow("SELECT created_at FROM notifications WHERE id = ?", lastInsertID).Scan(&createdAt)

	if err != nil {
		return 0, "", err
	}

	return int(lastInsertID), createdAt, nil
}

func (repo *EventRepository) GetAll() ([]*evententity.Event, error) {
	query := "SELECT id, title, description, emitter, created_at, topic FROM notifications"
	rows, err := repo.db.Query(query)

	if err != nil {
		return nil, fmt.Errorf("error al obtener datos %w", err)
	}
	defer rows.Close()

	var events []*evententity.Event

	for rows.Next() {
		var event evententity.Event
		var createdAtRaw []uint8

		if err := rows.Scan(&event.Id, &event.Title, &event.Description, &event.Emitter, &createdAtRaw, &event.Topic); err != nil {
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
	query := "DELETE FROM notifications"
	_, err := repo.db.Exec(query)
	return err
}
