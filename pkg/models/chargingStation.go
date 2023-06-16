package models

import "github.com/guregu/null"

type ChargingStation struct {
	StationID         int       `json:"stationID" gorm:"column:station_id;unique"`
	EnergyOutput      string    `json:"energyOutput" gorm:"column:energy_output"`
	Type              string    `json:"type" gorm:"column:type"`
	Occupied          bool      `json:"occupied" gorm:"column:occupied"`
	ChargingStartTime null.Time `json:"chargingStartTime" gorm:"column:charging_start_time"`
}
