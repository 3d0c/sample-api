package helpers

import (
	"encoding/json"
	"io"
)

func Decode(r io.ReadCloser, v interface{}) error {
	defer r.Close()
	if err := json.NewDecoder(r).Decode(v); err != nil {
		return err
	}

	return nil
}
