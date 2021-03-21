package handlers

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/3d0c/sample-api/pkg/rpc"

	"github.com/3d0c/sample-api/api/models"
	"github.com/3d0c/sample-api/pkg/helpers"
)

func TestFlightFlow(t *testing.T) {
	TestCreateUser(t)

	id := testCreateFlight(t)
	testUpdateFlight(t, id)
	testSearchFlight(t)
	testRemoveFlight(t, id)
}

func testCreateFlight(t *testing.T) uint {
	endpoint := "http://" + listenOn + "/flights"

	expected := models.Flight{
		Name:   "test",
		Number: "AB123",
		// Scheduled: time.Now(),
		// Arrival:     time.Now().Add(time.Minute),
		// Departure:   time.Now().Add(time.Minute),
		Destination: "Moscow",
		Fare:        100,
		Duration:    60,
	}

	payload, err := json.MarshalIndent(expected, "", "    ")
	if err != nil {
		t.Fatalf("Error marshalling struct - %s\n", err)
	}

	r, err := rpc.Request("POST", endpoint, payload, rpcCfg)
	if err != nil {
		t.Fatalf("Error requesting %s - %s\n", endpoint, err)
	}

	if r.StatusCode != 200 {
		t.Fatalf("\nExpected status code: %d\nObtained: %d\n", 200, r.StatusCode)
	}

	obtained := models.Flight{}

	if err := helpers.Decode(r.Body, &obtained); err != nil {
		t.Fatalf("Unexpected error - %s\n", err)
	}

	if obtained.ID == 0 {
		t.Fatalf("Expected non 0 flight id\n")
	} else {
		expected.ID = obtained.ID
	}

	if !reflect.DeepEqual(expected, obtained) {
		t.Fatalf("\nExpected: %v\nObtained: %v\n", expected, obtained)
	}

	return obtained.ID
}

func testUpdateFlight(t *testing.T, flightID uint) {
	endpoint := fmt.Sprintf("http://%s/flights/%d", listenOn, flightID)
	payload := `{"duration": 70}`

	r, err := rpc.Request("PUT", endpoint, []byte(payload), rpcCfg)
	if err != nil {
		t.Fatalf("Error requesting %s - %s\n", endpoint, err)
	}

	if r.StatusCode != 200 {
		t.Fatalf("\nExpected status code: %d\nObtained: %d\n", 200, r.StatusCode)
	}
}

func testSearchFlight(t *testing.T) {
	endpoint := fmt.Sprintf("http://%s/flights?destination=Moscow", listenOn)

	r, err := rpc.Request("GET", endpoint, nil, rpcCfg)
	if err != nil {
		t.Fatalf("Error requesting %s - %s\n", endpoint, err)
	}

	obtained := []models.Flight{}

	if err := helpers.Decode(r.Body, &obtained); err != nil {
		t.Fatalf("Unexpected error - %s\n", err)
	}

	if r.StatusCode != 200 {
		t.Fatalf("\nExpected status code: %d\nObtained: %d\n", 200, r.StatusCode)
	}

	if len(obtained) == 0 {
		t.Fatalf("Expected non 0 records\n")
	}

	if obtained[0].Destination != "Moscow" {
		t.Fatalf("\nExpected destination: %s\nObtained destination: %s\n", "Moscow", obtained[0].Destination)
	}
}

func testRemoveFlight(t *testing.T, flightID uint) {
	endpoint := fmt.Sprintf("http://%s/flights/%d", listenOn, flightID)

	r, err := rpc.Request("DELETE", endpoint, nil, rpcCfg)
	if err != nil {
		t.Fatalf("Error requesting %s - %s\n", endpoint, err)
	}

	if r.StatusCode != 200 {
		t.Fatalf("\nExpected status code: %d\nObtained: %d\n", 200, r.StatusCode)
	}
}
