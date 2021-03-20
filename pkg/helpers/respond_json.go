package helpers

import (
	"encoding/json"
	"log"
	"net/http"
)

type jsonResponder struct {
	w http.ResponseWriter
}

func NewJsonResponder(w http.ResponseWriter) jsonResponder {
	return jsonResponder{w: w}
}

func (j jsonResponder) Encode(v interface{}) ([]byte, error) {
	var (
		b   []byte
		err error
	)

	if b, err = json.MarshalIndent(v, "", "    "); err != nil {
		return nil, err
	}

	return b, nil
}

func (j jsonResponder) Write(v interface{}) {
	var (
		b   []byte
		err error
	)

	if b, err = j.Encode(v); err != nil {
		http.Error(j.w, "", http.StatusInternalServerError)
		return
	}

	if _, err = j.w.Write(b); err != nil {
		log.Printf("Error writing response - %s\n", err)
	}

	return
}
