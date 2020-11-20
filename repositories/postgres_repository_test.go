package repositories

import (
	"github.com/4thel00z/libservice"
	"github.com/jmoiron/sqlx"
	_ "github.com/proullon/ramsql/driver"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

// This is just an example entity, and shows you how you would create one yourself
type Mango struct {
	*libservice.DefaultEntity
	Color string `name:"color" sql:"VARCHAR(32)"`
}

func (m Mango) Name() string {
	return "mango"
}

func (m Mango) Type() reflect.Type {
	return reflect.TypeOf(m)
}

func (m Mango) Value(key string) (interface{}, error) {
	switch key {
	case "id":
		return m.ID, nil
	case "color":
		return m.Color, nil
	}
	return nil, libservice.FieldNotFound
}

func TestSave(t *testing.T) {
	repository := PostgresRepository{}
	db, err := sqlx.Open("ramsql", "TestSave")
	assert.Nil(t, err)
	repository.DB = db
	mango := Mango{
		Color: "#ffffff",
	}
	err = repository.CreateTable(mango, true)
	assert.Nil(t, err)
	mango.DefaultEntity = &libservice.DefaultEntity{}
	err = repository.Save(mango)
	assert.Nil(t, err)
}
