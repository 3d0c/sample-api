package helpers

import (
	"net/http"
)

type Responder interface {
	Encode(v interface{}) []byte
	Write(w http.ResponseWriter)
}
