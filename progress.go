package progress

import (
	"github.com/satori/go.uuid"
	"sync"
)

var (
	bars         sync.Map
	impls        sync.Map
	updateChan   chan update
	doUpdateChan chan interface{}
	WG           sync.WaitGroup
)

func init() {
	go func() {
		for u := range updateChan {
			prev, exists := bars.Load(u.ID())
			if !exists {
				prev = 0
			}
			bars.Store(u.ID(), u.NewVal(prev.(float32)))
			doUpdateChan <- struct{}{}
		}
	}()

	go func() {
		for range doUpdateChan {
			impls.Range(func(id interface{}, bar interface{}) bool {
				val, _ := bars.Load(id)
				bar.(ProgressBar).Update(val.(float32))
				return true
			})
		}
	}()
}

func UpdateProgressFixed(id uuid.UUID, progress float32) {
	WG.Add(1)
	go func(i uuid.UUID, p float32) {
		defer WG.Done()

		updateChan <- newFixedUpdate(i, p)
	}(id, progress)
}

func UpdateProgressDelta(id uuid.UUID, progress float32) {
	WG.Add(1)
	go func(i uuid.UUID, p float32) {
		defer WG.Done()

		updateChan <- newDeltaUpdate(i, p)
	}(id, progress)
}
