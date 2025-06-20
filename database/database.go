package database

// database/database.go
// Package database provides functionalities to manage database connections.
// It uses the crypton package for secure handling of encrypted passwords.
// It requires a configuration struct that contains database connection details and an encryption configuration.
// It is designed to be used with PostgreSQL and supports setting the search path for schemas.
// It provides a DBManager struct to handle the database connection and operations.

import (
	"database/sql"
	"fmt"

	"github.com/Shauanth/Singleton_Encription_ServiceGolang/crypton"
	_ "github.com/lib/pq"
)

// Config representa la estructura del JSON
type Config struct {
	Driver       string `json:"driver"`
	TipoDriver   string `json:"tipo_de_driver"`
	DBName       string `json:"base_de_datos"`
	Host         string `json:"nombre_de_host"`
	Port         string `json:"puerto"`
	User         string `json:"usuario"`
	DatabaseName string `json:"esquemabd"`
	Password     string `json:"contrasenha"` // Asumimos que está cifrada
}

// DBManager maneja la conexión a la BD
type DBManager struct {
	DB *sql.DB
}

func NuevoDBManager(config Config, configuracion crypton.Config) (*DBManager, error) {
	// Descifrar la contraseña usando el campo Password del struct Config
	password, err := crypton.Decrypt(config.Password, configuracion)
	if err != nil {
		return nil, fmt.Errorf("error al descifrar password: %v", err)
	}
	// Crear la cadena de conexión para PostgreSQL
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		config.Host,
		config.Port,
		config.User,
		password, // Contraseña descifrada
		config.DBName,
	)
	db, err := sql.Open(config.Driver, connStr)
	if err != nil {
		return nil, fmt.Errorf("error al conectar a PostgreSQL: %v", err)
	}
	_, err = db.Exec(fmt.Sprintf("SET search_path TO '%s', public@%s", config.DatabaseName, config.DBName))
	if err != nil {
		return nil, fmt.Errorf("error al configurar search_path: %v", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error al hacer ping a la BD: %v", err)
	}
	return &DBManager{DB: db}, nil
}

// Cerrar cierra la conexión
func (m *DBManager) Cerrar() {
	m.DB.Close()
}
