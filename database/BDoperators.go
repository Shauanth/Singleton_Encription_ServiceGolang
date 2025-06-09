package database

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

// CRUD proporciona operaciones genéricas para cualquier tabla
type CRUD struct {
	DB *sql.DB
}

// NuevoCRUD crea una nueva instancia de CRUD
func NuevoCRUD(db *sql.DB) *CRUD {
	return &CRUD{DB: db}
}

// Insertar inserta un registro en la tabla especificada
func (c *CRUD) Insertar(tabla string, datos interface{}) error {
	// Obtener el tipo y valor de la estructura
	v := reflect.ValueOf(datos)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()

	// Preparar los nombres de las columnas y los valores
	var columns []string
	var placeholders []string
	var values []interface{}

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i).Interface()

		// Obtener el nombre de la columna (usando el tag `db` o el nombre del campo)
		column := strings.ToLower(field.Name)
		if dbTag := field.Tag.Get("db"); dbTag != "" {
			column = strings.Split(dbTag, ",")[0]
		}

		columns = append(columns, column)
		placeholders = append(placeholders, fmt.Sprintf("$%d", i+1))
		values = append(values, value)
	}

	// Construir la consulta SQL
	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		tabla,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
	)

	// Ejecutar la consulta
	_, err := c.DB.Exec(query, values...)
	if err != nil {
		return fmt.Errorf("error al insertar en la tabla %s: %v", tabla, err)
	}

	return nil
}

// Actualizar actualiza registros en la tabla especificada con una condición flexible
func (c *CRUD) Actualizar(tabla string, datos interface{}, whereClause string, whereArgs ...interface{}) error {
	// Obtener el tipo y valor de la estructura
	v := reflect.ValueOf(datos)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()

	// Preparar las columnas y los valores
	var updates []string
	var values []interface{}

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i).Interface()

		// Obtener el nombre de la columna (usando el tag `db` o el nombre del campo)
		column := field.Name
		if dbTag := field.Tag.Get("db"); dbTag != "" {
			column = strings.Split(dbTag, ",")[0]
		}
		updates = append(updates, fmt.Sprintf("%s = $%d", column, i+1))
		values = append(values, value)
	}
	// Agregar los argumentos del WHERE al final
	values = append(values, whereArgs...)
	// Construir la consulta SQL
	query := fmt.Sprintf(
		"UPDATE %s SET %s WHERE %s",
		tabla,
		strings.Join(updates, ", "),
		whereClause,
	)

	// Ejecutar la consulta
	_, err := c.DB.Exec(query, values...)
	if err != nil {
		return fmt.Errorf("error al actualizar en la tabla %s: %v", tabla, err)
	}

	return nil
}

// Eliminar elimina un registro de la tabla especificada
func (c *CRUD) Eliminar(tabla string, id string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", tabla)
	_, err := c.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error al eliminar de la tabla %s: %v", tabla, err)
	}
	return nil
}

func (c *CRUD) Seleccionar(tabla string, columnas []string, whereClause string, args ...interface{}) (*sql.Rows, error) {
	// Construir la lista de columnas
	cols := "*"
	if len(columnas) > 0 {
		cols = strings.Join(columnas, ", ")
	}
	// Construir la consulta SQL
	query := fmt.Sprintf("SELECT %s FROM %s", cols, tabla)

	// Agregar WHERE si se especificó
	if whereClause != "" {
		query += " WHERE " + whereClause
	}
	// Ejecutar la consulta
	rows, err := c.DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error al realizar select en tabla %s: %v", tabla, err)
	}
	return rows, nil
}
