package response

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Json(w http.ResponseWriter, status int, data interface{}) error {
	b, err := json.Marshal(&data)
	if err != nil {
		return fmt.Errorf("error traying converting data")
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(b)
	if err != nil {
		return fmt.Errorf("error writing data")
	}
	return nil
}
