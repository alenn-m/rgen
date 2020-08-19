package resp

import (
	"encoding/json"
	"net/http"
)

func ReturnSuccess(w http.ResponseWriter, data interface{}) {
	r := make(map[string]interface{})

	r["code"] = http.StatusOK
	r["data"] = data

	result, _ := json.Marshal(r)

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}

func ReturnError(w http.ResponseWriter, message string, code int) {
	r := make(map[string]interface{})

	r["code"] = code
	r["error"] = message

	result, _ := json.Marshal(r)

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}
