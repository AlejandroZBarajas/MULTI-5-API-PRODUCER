package evententity

import "time"

type Event struct {
	Id          int32     `json:"id"`
	Serie       string    `json:"serie"`
	Title       string    `json:"title"`
	Description float32   `json:"description"`
	Emitter     string    `json:"emitter"`
	Topic       string    `json:"topic"`
	Created_at  time.Time `json:"created_at"`
}

func CreateEvent(serie string, title string, description float32, emitter string, topic string) *Event {
	return &Event{Serie: serie, Title: title, Description: description, Emitter: emitter, Topic: topic}
}
