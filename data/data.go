package data

import (
	"encoding/json"
	"log"
	"net/http"
)

func EncodeResponse(payload interface{}, w http.ResponseWriter, status int) {
	resp := make(map[string]interface{}, 0)
	if 200 <= status && status <= 300 {
		resp["data"] = payload
	} else {
		if err, ok := payload.(error); ok {
			resp["error"] = err.Error()
		} else {
			resp["error"] = payload
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding json response: %v\n", err)
	}
}
