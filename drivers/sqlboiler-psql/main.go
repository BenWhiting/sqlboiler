package main

import (
	"github.com/BenWhiting/sqlboiler/v4/drivers"
	"github.com/BenWhiting/sqlboiler/v4/drivers/sqlboiler-psql/driver"
)

func main() {
	drivers.DriverMain(&driver.PostgresDriver{})
}
