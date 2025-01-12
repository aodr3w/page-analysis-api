package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func main() {
	server := NewServer()
	if err := http.ListenAndServe("localhost:5000", server); err != nil {
		log.Fatal(err)
	}
}

type FindRequest struct {
	Url        string   `json:"url"`
	Attributes []string `json:"attributes"`
}

func decodeRequest(request *http.Request) (*FindRequest, error) {
	findRequest := FindRequest{}
	if err := json.NewDecoder(request.Body).Decode(&findRequest); err != nil {
		return nil, err
	} else {
		return &findRequest, nil
	}
}

func encodeResponse(payload interface{}, w http.ResponseWriter, status int) {
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
func fetchHTML(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}
	return body, nil
}

func NewServer() *http.ServeMux {
	m := http.ServeMux{}
	m.HandleFunc("/find", func(w http.ResponseWriter, r *http.Request) {
		bytesTokenLimit := 9216
		//parse the product information from the url
		req, err := decodeRequest(r)
		if err != nil {
			//write response
			encodeResponse(err, w, 500)
			return
		}
		respBytes, err := fetchHTML(req.Url)
		if err != nil {
			//write response
			encodeResponse(err, w, 500)
		}
		if len(respBytes) > bytesTokenLimit {
			respBytes = respBytes[:bytesTokenLimit]
		}
		encodeResponse(string(respBytes), w, 200)

		//fetch the html from the web site
		//if the site is not available respond accordingly
		//if it is available, forward the content to an LLM with a large enough context window
	})

	return &m
}
