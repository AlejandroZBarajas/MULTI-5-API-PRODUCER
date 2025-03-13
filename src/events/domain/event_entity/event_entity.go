package evententity

import "time"

type Event struct {
	Id          int32     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Emitter     string    `json:"emitter"`
	Created_at  time.Time `json:"created_at"`
}

func CreateEvent(title string, description string, emitter string) *Event {
	return &Event{Title: title, Description: description, Emitter: emitter}
}
