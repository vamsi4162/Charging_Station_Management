# Charging_Station_Management

Created a Charging Station Management application which has the following functionalities

Add Charging Stations
Start Charging
Available Charging Stations
Occupied ChargingStations Along With Their Availability Time.

Steps to Run The code

1. Install GOLANG and set path correctly.
2. Create mod file with command " go mod init <module-name>
3. Install all the necessary dependencies with command "go mod tidy"
4.To run the code use : " go run main.go.

This will start the server in port 8080

Check whether the server is running or not by simply opening the new chrome tab and enter the following URL: "localhost:8080" it should give some response. Since we have not passed any data it should give "error 404".

You can use API testing tools like "POSTMAN" or "ThundeClient(vscode extension)" to check its functioning.
Sample Requests and Responses

Here I added 5 Charging Stations to database

1.Add Charging Station:

POST localhost:8080/charging-stations

Request Body: 
{
  "stationID": 1,
  "energyOutput": "100kWh",
  "type": "DC",
 "Occupied":false
}

Response: 
{
    Charging station added successfully
}

2.Start Charging:

POST localhost:8080/charging/start

Request Body: 
{
  "stationID":1,
  "vehicleBatteryCapacity": "500kWh",
  "currentVehicleCharge": "100kWh"
}

Response: 
{
  "chargingStartTime": "2023-06-16T15:52:21.24757+05:30",
  "message": "Charging started successfully",
  "stationAvailabilityTime": "2023-06-16T19:52:21.24757+05:30"
}

3.Get Available Charging Stations:

GET  http://localhost:8080/charging/available

Response: 
{
  "source": "database",
  "stations": [
    {
      "stationID": 2,
      "energyOutput": "100kWh",
      "type": "DC",
      "occupied": false,
      "chargingStartTime": null
    },
    {
      "stationID": 3,
      "energyOutput": "100kWh",
      "type": "DC",
      "occupied": false,
      "chargingStartTime": null
    },
    {
      "stationID": 4,
      "energyOutput": "100kWh",
      "type": "DC",
      "occupied": false,
      "chargingStartTime": null
    },
    {
      "stationID": 5,
      "energyOutput": "100kWh",
      "type": "DC",
      "occupied": false,
      "chargingStartTime": null
    }
  ]
}

4. Get Occupied Charging Stations:

GET  http://localhost:8080/charging/occupied

Response:
{
  "source": "cache",
  "stations": [
    {
      "stationID": 1,
      "vehicleBatteryCapacity": "500kWh",
      "currentVehicleCharge": "100kWh",
      "chargingStartTime": "2023-06-16T15:52:21.248+05:30",
      "stationAvailabilityTime": "2023-06-16T19:52:21.248+05:30"
    }
  ]
}
