package tests

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/magiconair/properties/assert"
)

var ip string
var port string
var prefix string

type GetDistanceBody struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Distance int    `json:"distance"`
	Unit     string `json:"unit"`
}

func init() {
	ip = "127.0.0.1"
	port = "8000"
	prefix = "http://" + ip + ":" + port
}

func TestGetDistance(t *testing.T) {
	//test get distance between two not existing DS100 code
	fromN := "7ABCI" // not existing DS100 code
	toN := "9KIL"    // not existing DS100 code
	_, statusCode := RequestGetDistance(fromN, toN)
	assert.Equal(t, statusCode, http.StatusNotFound)

	//test get distance between one existing and one not exiting DS100 code
	from := "TKL" //  existing DS100 code
	toN = "7ABCI" // not existing DS100 code
	_, statusCode = RequestGetDistance(from, toN)
	assert.Equal(t, statusCode, http.StatusNotFound)

	//test get distance from two existing DS100 code
	from = "FF" //  existing DS100 code
	to := "BLS" // not existing DS100 code
	body, statusCode := RequestGetDistance(from, to)

	assert.Equal(t, statusCode, http.StatusOK)

	assert.Equal(t, body.From, "Frankfurt(Main)Hbf")
	assert.Equal(t, body.To, "Berlin Hbf")
	assert.Equal(t, body.Distance, 423)
	assert.Equal(t, body.Unit, "km")

}

func RequestGetDistance(from, to string) (GetDistanceBody, int) {
	url := fmt.Sprintf("%s/api/v1/distance/%s/%s", prefix, from, to)
	println(url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		println(err.Error())
	}
	///resp, err := http.Get(url)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		println(err.Error())
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		println(err.Error())
	}
	var body GetDistanceBody
	err = json.Unmarshal(b, &body)
	if err != nil {
		println(err.Error())
	}
	return body, resp.StatusCode
}
