package repositories

import (
	"fmt"
	"github.com/4thel00z/libservice/v1"
	"github.com/google/uuid"
	"github.com/iancoleman/strcase"
	"github.com/jmoiron/sqlx"
	"reflect"
	"strings"
)

var (
	defaultEntityType = reflect.TypeOf(&v1.DefaultEntity{})
)

type PostgresRepository struct {
	DB *sqlx.DB
}

func (r *PostgresRepository) CreateTable(e v1.Entity, ifNotExist bool) error {
	var (
		template string
	)
	// For now we expect well behaved entities and don't use reflection magic for default types
	fields, _, sqlProperties := GetFields(e)

	tableEntries := make([]string, len(fields))
	for i, field := range fields {
		tableEntries[i] = fmt.Sprintf("%s %s", field, sqlProperties[i])
	}
	if ifNotExist {
		template = "CREATE TABLE IF NOT EXISTS %s (%s);"
	} else {
		template = "CREATE TABLE %s (%s);"
	}
	query := fmt.Sprintf(template, e.Name(), strings.Join(tableEntries, ","))
	_, err := r.DB.Exec(query)
	return err
}

func (r *PostgresRepository) DropTable(i v1.Entity, ifExists bool) error {
	panic("implement me")
}

func (r *PostgresRepository) Open(dataSourceName string) error {
	db, err := sqlx.Connect("postgres", dataSourceName)
	if err != nil {
		return err
	}
	r.DB = db
	return nil
}

func (r *PostgresRepository) Save(e v1.Entity, fields ...string) error {
	session, err := r.DB.Beginx()
	if err != nil {
		return err
	}

	if e.Index() == [16]byte{0} {
		newUUID, err := uuid.NewUUID()
		if err != nil {
			return err
		}
		e.SetIndex(newUUID)
	}

	if len(fields) == 0 {
		fields, _, _ = GetFields(e)
	}

	tableName := e.Name()
	values := make([]interface{}, len(fields))
	placeHolders := make([]string, len(fields))
	for i, field := range fields {
		val, err := e.Value(field)
		if err != nil {
			return err
		}
		values[i] = val
		placeHolders[i] = fmt.Sprintf("$%d", i+1)
		i++
	}
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s);", tableName, strings.Join(fields, ","), strings.Join(placeHolders, ","))
	fmt.Println(query)
	_, err = session.Exec(query, values...)
	return err
}

func GetFields(e v1.Entity) ([]string, []reflect.Type, []string) {
	t := e.Type()
	fields := make([]string, t.NumField())
	types := make([]reflect.Type, t.NumField())
	sqlProperties := make([]string, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Type == defaultEntityType {
			fields[i] = "id"
			types[i] = reflect.TypeOf(uuid.UUID{})
			sqlProperties[i] = "BIGSERIAL PRIMARY KEY"
			continue
		}
		fieldName := strcase.ToSnake(field.Name)
		tag, ok := field.Tag.Lookup(v1.StructTagName)
		if !ok {
			fieldName = tag
		}
		fields[i] = fieldName
		types[i] = field.Type
		sqlProperties[i] = field.Tag.Get(v1.StructTagSQL)

	}
	return fields, types, sqlProperties
}

func (r *PostgresRepository) Update(i v1.Entity, fields ...string) error {
	panic("implement me")
}

func (r *PostgresRepository) Get(i uuid.UUID) (v1.Entity, error) {
	panic("implement me")
}

func (r *PostgresRepository) List() []v1.Entity {
	panic("implement me")
}

func (r *PostgresRepository) Delete(i v1.Entity) (bool, error) {
	panic("implement me")
}
