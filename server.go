package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"

	"github.com/gorilla/mux"
)

type ResponseMessage struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Distance int    `json:"distance"`
	Unit     string `json:"unit"`
}

var handler CsvHandler

func main() {

	handler = CsvHandler{filename: "dataset.CSV"}
	handler.LoadData()

	r := mux.NewRouter()

	r.HandleFunc("/api/v1/distance/{from}/{to}", GetDistance).Methods(http.MethodGet)

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":8000", r))
}

func GetDistance(response http.ResponseWriter, r *http.Request) {

	fromDS100 := mux.Vars(r)["from"]
	toDS100 := mux.Vars(r)["to"]

	fromStationName, err, statusCode := handler.GetStationName(fromDS100)
	if err != nil {
		response.WriteHeader(statusCode)
		response.Write([]byte(err.Error()))
		return
	}
	toStationName, err, statusCode := handler.GetStationName(toDS100)
	if err != nil {
		response.WriteHeader(statusCode)
		response.Write([]byte(err.Error()))
		return
	}
	pointFrom, err, statusCode := handler.GetPoint(fromDS100)
	if err != nil {
		response.WriteHeader(statusCode)
		response.Write([]byte(err.Error()))
		return
	}
	pointTo, err, statusCode := handler.GetPoint(toDS100)
	if err != nil {
		response.WriteHeader(statusCode)
		response.Write([]byte(err.Error()))
		return
	}

	distance := calcDistance(pointFrom, pointTo)
	responseMessage := ResponseMessage{From: fromStationName, To: toStationName, Distance: distance, Unit: "km"}

	jsonResponse, jsonError := json.Marshal(responseMessage)
	if jsonError != nil {
		fmt.Println("Unable to encode JSON")
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	response.Write(jsonResponse)

}

func calcDistance(point1, point2 Point) int {
	const earthRadius = 6371 // Radius of Earth in kilometers

	// Convert degrees to radians
	lat1 := toRadians(point1.latitude)
	long1 := toRadians(point1.longitude)
	lat2 := toRadians(point2.latitude)
	long2 := toRadians(point2.longitude)

	// Haversine formula
	a := math.Sin(lat2-lat1) / 2.0
	b := math.Sin(long2-long1) / 2.0
	c := math.Cos(lat1) * math.Cos(lat2)
	d := a*a + c*b*b

	return int(earthRadius * 2 * math.Atan2(math.Sqrt(d), math.Sqrt(1-d)))
}

func toRadians(degrees float64) float64 {
	return degrees * math.Pi / 180.0
}
