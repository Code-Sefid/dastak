package main

import (
	"github.com/soheilkhaledabdi/dastak/api"
	"github.com/soheilkhaledabdi/dastak/config"
	"github.com/soheilkhaledabdi/dastak/data/db"
	"github.com/soheilkhaledabdi/dastak/data/migration"
	"github.com/soheilkhaledabdi/dastak/pkg/logging"
)

func main() {
	cfg := config.GetConfig()
	logger := logging.NewLogger(cfg)

	err := db.InitDb(cfg)
	defer db.CloseDb()

	if err != nil {
		logger.Fatal(logging.Postgres, logging.Startup, err.Error(), nil)
	}
	migration.Up()
	api.InitServer()
}
