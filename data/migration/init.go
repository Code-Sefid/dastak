package migration

import (
	"github.com/soheilkhaledabdi/dastak/config"
	"github.com/soheilkhaledabdi/dastak/data/db"
	"github.com/soheilkhaledabdi/dastak/data/models"
	"github.com/soheilkhaledabdi/dastak/pkg/logging"
	"gorm.io/gorm"
)

var logger = logging.NewLogger(config.GetConfig())

func Up() {
	database := db.GetDb()
	createTables(database)
	logger.Info(logging.Postgres, logging.Migration, "UP", nil)
}

func addNewTable(database *gorm.DB, model interface{}, tables []interface{}) []interface{} {
	if !database.Migrator().HasTable(model) {
		tables = append(tables, model)
	}
	return tables
}

func createTables(database *gorm.DB) {

	modelsList := []interface{}{
		&models.JWTToken{},
		&models.Users{},
		&models.Categories{},
		&models.Products{},
		&models.Customer{},
		&models.Factors{},
		&models.FactorProducts{},
		&models.FactorPayment{},
		&models.Wallet{},
		&models.BankAccounts{},
		&models.Transactions{},
	}

	tables := []interface{}{}

	for _, model := range modelsList {
		tables = addNewTable(database, model, tables)
	}

	err := database.Migrator().CreateTable(tables...)
	if err != nil {
		logger.Error(logging.Postgres, logging.Migration, err.Error(), nil)
	}
	logger.Info(logging.Postgres, logging.Migration, "tables created", nil)
}
