package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/TheYaminaK/Go-Binance-OrderBook/middleware"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

func main() {
	fmt.Println("Iniciando servicios...")

	timeEnvStr := os.Getenv("TIME")
	timeEnv, _ := strconv.Atoi(timeEnvStr)

	db := middleware.GetConnection()
	defer db.Close()

	err := middleware.CreateTables(db)
	if err != nil {
		fmt.Println(errors.Wrap(err, "Error al verificar y crear las tablas en la base de datos"))
		return
	}

	fmt.Println("Verificación de base de datos completada. Tablas creadas correctamente.")

	for {
		err := middleware.BinanceData()
		if err != nil {
			fmt.Println(errors.Wrap(err, "Error al obtener datos de Binance"))
		}

		time.Sleep(time.Duration(timeEnv) * time.Second)
	}
}
