package main

import (
	"github.com/BenWhiting/sqlboiler/v4/drivers"
	"github.com/BenWhiting/sqlboiler/v4/drivers/sqlboiler-sqlite3/driver"
)

func main() {
	drivers.DriverMain(&driver.SQLiteDriver{})
}
