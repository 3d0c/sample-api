package handlers

import (
	"net/http"
	"testing"

	"github.com/3d0c/sample-api/api/models"
	"github.com/3d0c/sample-api/pkg/helpers"
	"github.com/3d0c/sample-api/pkg/rpc"
)

func TestCreateUser(t *testing.T) {
	endpoint := "http://" + listenOn + "/users"
	payload := `{"name": "test", "password": "test"}`

	r, err := rpc.Request("POST", endpoint, []byte(payload), nil)
	if err != nil {
		t.Fatalf("Error requesting %s - %s\n", endpoint, err)
	}

	if r.StatusCode != 200 {
		t.Fatalf("\nExpected status code: %d\nObtained: %d\n", 200, r.StatusCode)
	}

	result := models.User{}

	if err := helpers.Decode(r.Body, &result); err != nil {
		t.Fatalf("Unexpected error - %s\n", err)
	}

	if result.ID == 0 {
		t.Fatalf("Expected non 0 user id\n")
	}

	if result.Name != "test" {
		t.Fatalf("\nExpected user name: %s\nObtained: %s\n", "test", result.Name)
	}

	testLoginUser(t, result.ID)
}

func testLoginUser(t *testing.T, userID uint) {
	endpoint := "http://" + listenOn + "/users/login"
	payload := `{"name": "test", "password": "test"}`

	r, err := rpc.Request("POST", endpoint, []byte(payload), nil)
	if err != nil {
		t.Fatalf("Error requesting %s - %s\n", endpoint, err)
	}

	if r.StatusCode != 200 {
		t.Fatalf("\nExpected status code: %d\nObtained: %d\n", 200, r.StatusCode)
	}

	result := models.JWTToken{}

	if err := helpers.Decode(r.Body, &result); err != nil {
		t.Fatalf("Unexpected error - %s\n", err)
	}

	if len(result.Token) == 0 {
		t.Fatalf("Expected non 0 length token\n")
	}

	u := &models.User{ID: userID}

	if err := u.Delete(); err != nil {
		t.Fatalf("Error remove temoporary user - %s\n", err)
	}

	rpcCfg = &rpc.Config{Headers: make(http.Header)}
	rpcCfg.Headers.Set("Authorization", "Bearer "+result.Token)
}
