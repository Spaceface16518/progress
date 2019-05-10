package progress

import uuid "github.com/satori/go.uuid"

func NewBarID() uuid.UUID {
	return uuid.Must(newUUID())
}

func newUUID() (uuid.UUID, error) {
	return uuid.NewV4()
}
