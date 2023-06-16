package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/vamsi4162/Charging_Station_Management/pkg/database"
	"github.com/vamsi4162/Charging_Station_Management/pkg/models"
	"github.com/vamsi4162/Charging_Station_Management/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/guregu/null"
)

func AddChargingStation(c *gin.Context) {
	var payload models.ChargingStation
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	var existingStation models.ChargingStation
	result := database.DB.First(&existingStation, "station_id = ?", payload.StationID)
	if result.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Charging station already exists"})
		return
	}

	chargingStation := models.ChargingStation{
		StationID:    payload.StationID,
		EnergyOutput: payload.EnergyOutput,
		Type:         payload.Type,
		Occupied:     false,
	}

	result = database.DB.Create(&chargingStation)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add the charging station"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Charging station added successfully"})
}

func StartCharging(c *gin.Context) {
	var payload struct {
		StationID              int    `json:"stationID"`
		VehicleBatteryCapacity string `json:"vehicleBatteryCapacity"`
		CurrentVehicleCharge   string `json:"currentVehicleCharge"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	var station models.ChargingStation
	result := database.DB.First(&station, "station_id = ?", payload.StationID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Charging station not found"})
		return
	}

	if station.Occupied {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Charging station is already occupied"})
		return
	}

	chargingStartTime := time.Now().Local()

	availabilityTime, err := calculateAvailabilityTime(payload.VehicleBatteryCapacity, payload.CurrentVehicleCharge, chargingStartTime, station.EnergyOutput)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to calculate availability time"})
		return
	}

	occupiedStation := models.OccupiedChargingStation{
		StationID:               payload.StationID,
		VehicleBatteryCapacity:  payload.VehicleBatteryCapacity,
		CurrentVehicleCharge:    payload.CurrentVehicleCharge,
		ChargingStartTime:       null.TimeFrom(chargingStartTime),
		StationAvailabilityTime: availabilityTime,
	}

	result = database.DB.Create(&occupiedStation)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start charging"})
		return
	}

	result = database.DB.Model(&models.ChargingStation{}).Where("station_id = ? AND occupied = ? AND charging_start_time IS NULL", payload.StationID, false).Updates(map[string]interface{}{
		"occupied":            true,
		"charging_start_time": null.TimeFrom(chargingStartTime),
	})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update the charging station status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Charging started successfully", "chargingStartTime": chargingStartTime, "stationAvailabilityTime": availabilityTime})
}

func GetAvailableChargingStations(c *gin.Context) {
	// Check cache first
	cacheKey := "available_charging_stations"
	if cacheData, found := utils.Cache.Get(cacheKey); found {
		if stations, ok := cacheData.([]models.ChargingStation); ok {
			fmt.Println("Retrieved available charging stations from cache")
			c.JSON(http.StatusOK, gin.H{"stations": stations, "source": "cache"})
			return
		}
	}

	var stations []models.ChargingStation
	result := database.DB.Find(&stations, "occupied = ?", false)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve available charging stations"})
		return
	}

	// Update cache
	utils.Cache.Set(cacheKey, stations, utils.CacheTTL)
	fmt.Println("Retrieved available charging stations from the database")

	c.JSON(http.StatusOK, gin.H{"stations": stations, "source": "database"})
}

func GetOccupiedChargingStations(c *gin.Context) {
	// Check cache first
	cacheKey := "occupied_charging_stations"
	if cacheData, found := utils.Cache.Get(cacheKey); found {
		if occupiedStations, ok := cacheData.([]models.OccupiedChargingStation); ok {
			fmt.Println("Retrieved occupied charging stations from cache")
			c.JSON(http.StatusOK, gin.H{"stations": occupiedStations, "source": "cache"})
			return
		}
	}

	var occupiedStations []models.OccupiedChargingStation
	result := database.DB.Find(&occupiedStations)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve occupied charging stations"})
		return
	}

	// Update cache
	utils.Cache.Set(cacheKey, occupiedStations, utils.CacheTTL)
	fmt.Println("Retrieved occupied charging stations from the database")

	c.JSON(http.StatusOK, gin.H{"stations": occupiedStations, "source": "database"})
}


func calculateAvailabilityTime(vehicleBatteryCapacity, currentVehicleCharge string, chargingStartTime time.Time, energyOutput string) (time.Time, error) {
	batteryCap, err := parseEnergyValue(vehicleBatteryCapacity)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse vehicleBatteryCapacity: %v", err)
	}

	currentCharge, err := parseEnergyValue(currentVehicleCharge)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse currentVehicleCharge: %v", err)
	}

	energyOut, err := parseEnergyValue(energyOutput)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse energyOutput: %v", err)
	}

	if energyOut <= 0 {
		return time.Time{}, errors.New("failed to calculate availability time: energy output must be greater than zero")
	}

	remainingEnergy := batteryCap - currentCharge
	if remainingEnergy <= 0 {
		return time.Time{}, errors.New("failed to calculate availability time: current vehicle charge exceeds or equals vehicle battery capacity")
	}

	chargingDuration := time.Duration(remainingEnergy / energyOut * float64(time.Hour))
	availabilityTime := chargingStartTime.Add(chargingDuration)

	return availabilityTime, nil
}

func parseEnergyValue(value string) (float64, error) {
	var quantity float64
	unit := "kWh"
	_, err := fmt.Sscanf(value, "%f%s", &quantity, &unit)
	if err != nil {
		return 0, err
	}
	if unit != "kWh" {
		return 0, errors.New("invalid energy unit")
	}
	return quantity, nil
}
