package req

import (
	"encoding/json"
	"net/http"
)

type RequestObject interface {
	Validate() error
}

// DecodeRequest decodes request object
func DecodeRequest(r *http.Request, dst RequestObject) error {
	dec := json.NewDecoder(r.Body)

	err := dec.Decode(&dst)
	if err != nil {
		return err
	}

	err = dst.Validate()

	return err
}
