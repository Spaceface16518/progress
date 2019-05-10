package progress

import "github.com/satori/go.uuid"

type update interface {
	ID() uuid.UUID
	NewVal(prev float32) float32
}

type fixedUpdate struct {
	id       uuid.UUID
	progress float32
}

func newFixedUpdate(id uuid.UUID, progress float32) *fixedUpdate {
	if progress > 1 {
		progress = 1
	} else if progress < 0 {
		progress = 0
	}
	return &fixedUpdate{id: id, progress: progress}
}

func (update *fixedUpdate) ID() uuid.UUID {
	return update.id
}

func (update *fixedUpdate) NewVal(_ float32) float32 {
	return update.progress
}

type deltaUpdate struct {
	id       uuid.UUID
	progress float32
}

func newDeltaUpdate(id uuid.UUID, progress float32) *deltaUpdate {
	if progress > 1 {
		progress = 1
	} else if progress < 0 {
		progress = 0
	}
	return &deltaUpdate{id: id, progress: progress}
}

func (update *deltaUpdate) ID() uuid.UUID {
	return update.id
}

func (update *deltaUpdate) NewVal(prev float32) float32 {
	n := prev + update.progress
	if n > 1 {
		return 1
	} else if n < 0 {
		return 0
	} else {
		return n
	}
}
