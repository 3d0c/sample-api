package handlers

import (
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/3d0c/sample-api/api/models"
	"github.com/3d0c/sample-api/pkg/rpc"
)

const listenOn = "127.0.0.1:6677"

var (
	rpcCfg *rpc.Config
)

func TestMain(m *testing.M) {
	router := SetupRouter()

	if err := models.ConnectDatabase(); err != nil {
		log.Fatalf("Error connecting to database - %s\n", err)
	}

	go func() {
		log.Fatalln(
			http.ListenAndServe(listenOn, router),
		)
	}()

	os.Exit(m.Run())
}
