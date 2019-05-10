package progress

import uuid "github.com/satori/go.uuid"

func NewBarID() uuid.UUID {
	return uuid.Must(uuid.NewV4())
}
