package main

import (
	"log"

	"github.com/Shauanth/Singleton_Encription_ServiceGolang.git/database"
)

func main() {
	var NuevaConfiguracion database.Config
	NuevaConfiguracion.Driver = "postgres"
	NuevaConfiguracion.TipoDriver = "PostgreSQL"
	NuevaConfiguracion.DBName = "postgres"
	NuevaConfiguracion.User = "postgres"
	NuevaConfiguracion.Host = "localhost"
	NuevaConfiguracion.Port = "5432"
	NuevaConfiguracion.DatabaseName = "esquemabd"
	NuevaConfiguracion.Password = "a5i3aJtCcU0P56OTDmXSGb/kfkZY1/lEGdh5eVsbomGgL6ss7Q=="
	// Inicializar el gestor de base de datos
	dbManager, err := database.NuevoDBManager(NuevaConfiguracion)
	if err != nil {
		log.Fatalf("Error al conectar a la BD: %v", err)
	}
	defer dbManager.Cerrar()
}
