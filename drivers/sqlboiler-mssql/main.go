package main

import (
	"github.com/BenWhiting/sqlboiler/v4/drivers"
	"github.com/BenWhiting/sqlboiler/v4/drivers/sqlboiler-mssql/driver"
)

func main() {
	drivers.DriverMain(&driver.MSSQLDriver{})
}
