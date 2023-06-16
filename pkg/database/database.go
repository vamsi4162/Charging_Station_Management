package database

import (
	"github.com/vamsi4162/Charging_Station_Management/pkg/models"
	"fmt"
	"log"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Initialize() {

	dbHost := viper.GetString("db.host")
	dbPort := viper.GetString("db.port")
	dbName := viper.GetString("db.name")
	dbUser := viper.GetString("db.username")
	dbPassword := viper.GetString("db.password")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPassword, dbHost, dbPort, dbName)
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	err = DB.AutoMigrate(&models.ChargingStation{}, &models.OccupiedChargingStation{})
	if err != nil {
		log.Fatal("Failed to migrate the database:", err)
	}
}
