package evententity

import "time"

type Event struct {
	Id          int32     `json:"id"`
	Device_name string    `json:"device_name"`
	Created_at  time.Time `json:"created_at"`
}

func CreateEvent(Device_name string) *Event {
	return &Event{Device_name: Device_name}
}
