package progress

import uuid "github.com/satori/go.uuid"

type ProgressBar interface {
	Update(val float32)
}

func Register(bar *ProgressBar) uuid.UUID {
	id := uuid.Must(uuid.NewV4())

	impls.Store(id, bar)

	return id
}
