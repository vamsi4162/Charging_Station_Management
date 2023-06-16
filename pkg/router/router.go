package router

import (
	"github.com/vamsi4162/Charging_Station_Management/pkg/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	r.POST("/charging-stations", controllers.AddChargingStation)
	r.POST("/charging/start", controllers.StartCharging)
	r.GET("/charging/available", controllers.GetAvailableChargingStations)
	r.GET("/charging/occupied", controllers.GetOccupiedChargingStations)
}
