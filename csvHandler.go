package main

import (
	"encoding/csv"
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type CsvHandler struct {
	filename string
}
type Point struct {
	latitude  float64
	longitude float64
}

var data [][]string
var DS100_INDEX = 1
var LONGITUDE_INDEX = 5
var LATITUDE_INDEX = 6
var STATIONNAME_INDEX = 3

func (handler *CsvHandler) LoadData() error {
	file, err := os.Open(handler.filename)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	data = records
	return nil
}

func (handler *CsvHandler) GetPoint(DS100 string) (Point, error, int) {
	var coor Point
	var err error = nil
	var found bool = false

	for _, record := range data {

		if record[DS100_INDEX] == DS100 {
			coor.latitude, err = strconv.ParseFloat(strings.Replace(record[LATITUDE_INDEX], ",", ".", 1), 64)
			if err != nil {
				return coor, err, http.StatusInternalServerError
			}
			coor.longitude, err = strconv.ParseFloat(strings.Replace(record[LONGITUDE_INDEX], ",", ".", 1), 64)
			if err != nil {
				return coor, err, http.StatusInternalServerError
			}
			found = true
			break
		}
	}

	if !found {
		return coor, errors.New("Invalid GS100 code"), http.StatusNotFound
	}

	return coor, err, http.StatusAccepted
}

func (handler *CsvHandler) GetStationName(DS100 string) (string, error, int) {
	var err error = nil
	var found bool = false
	var stationName string = ""

	for _, record := range data {

		if record[DS100_INDEX] == DS100 {
			stationName = record[STATIONNAME_INDEX]
			found = true
			break
		}
	}

	if !found {
		return "", errors.New("invalid GS100 code"), http.StatusNotFound
	}

	return stationName, err, http.StatusAccepted
}
