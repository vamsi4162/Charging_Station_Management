package main

import (
	//"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/vamsi4162/Charging_Station_Management/pkg/database"
	"github.com/vamsi4162/Charging_Station_Management/pkg/router"
	"github.com/vamsi4162/Charging_Station_Management/pkg/utils"
)

func main() {

	err := utils.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	database.Initialize()
	utils.InitCache() 
	r := gin.Default()

	router.RegisterRoutes(r)

	err = r.Run(":8080")
	if err != nil {
		log.Fatal("Failed to start the server:", err)
	}
}
