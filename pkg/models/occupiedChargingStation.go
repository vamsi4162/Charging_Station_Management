package models

import (
	"time"

	"github.com/guregu/null"
)

type OccupiedChargingStation struct {
	StationID               int       `json:"stationID"`
	VehicleBatteryCapacity  string    `json:"vehicleBatteryCapacity"`
	CurrentVehicleCharge    string    `json:"currentVehicleCharge"`
	ChargingStartTime       null.Time `json:"chargingStartTime" gorm:"column:charging_start_time"`
	StationAvailabilityTime time.Time `json:"stationAvailabilityTime"`
}
