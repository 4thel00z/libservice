package repositories

import (
	"fmt"
	"github.com/4thel00z/libservice"
	"github.com/google/uuid"
)

type MemoryRepository map[uuid.UUID]libservice.Entity

func (m MemoryRepository) Save(e libservice.Entity) error {
	e, ok := m[e.Index()]
	if ok {
		return fmt.Errorf("could not save %s, since it is already present", e)
	}
	u, err := uuid.NewUUID()
	if err != nil {
		return err
	}
	e.SetIndex(u)
	m[e.Index()] = e
	return nil
}

func (m MemoryRepository) Update(e libservice.Entity) error {
	if e.Index() == [16]byte{0} {
		return m.Save(e)
	}
	m[e.Index()] = e
	return nil
}

func (m MemoryRepository) Get(u uuid.UUID) (libservice.Entity, error) {
	indexable, ok := m[u]
	if !ok {
		return nil, fmt.Errorf("could not find %s", u)
	}
	return indexable, nil
}

func (m MemoryRepository) List() []libservice.Entity {
	entities := make([]libservice.Entity, len(m))
	i := 0
	for _, v := range m {
		entities[i] = v
		i++
	}
	return entities
}

func (m MemoryRepository) Delete(e libservice.Entity) (bool, error) {
	_, ok := m[e.Index()]
	delete(m, e.Index())
	return ok, nil
}
