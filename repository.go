package libservice

import (
	"errors"
	"github.com/google/uuid"
	"reflect"
)

const (
	StructTagName = "name"
	StructTagSQL  = "sql"
)

var (
	FieldNotFound = errors.New("field was not found")
)

type DefaultEntity struct {
	ID uuid.UUID `persistence:"id"`
}

func (e *DefaultEntity) Index() uuid.UUID {
	if e == nil {
		*e = DefaultEntity{}
	}
	return e.ID
}

func (e *DefaultEntity) SetIndex(u uuid.UUID) {
	if e == nil {
		*e = DefaultEntity{}
	}
	e.ID = u
}

type Entity interface {
	Index() uuid.UUID
	SetIndex(u uuid.UUID)
	Name() string
	Type() reflect.Type
	Value(key string) (interface{}, error)
}

type Repository interface {
	CreateTable(i Entity, ifNotExist bool) error
	DropTable(i Entity, ifExists bool) error
	Save(i Entity, fields ...string) error
	Update(i Entity, fields ...string) error
	Get(i uuid.UUID) (Entity, error)
	List() []Entity
	Delete(i Entity) (bool, error)
}
