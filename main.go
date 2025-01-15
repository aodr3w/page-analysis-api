package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/aodr3w/extractor-api/common"
	"github.com/aodr3w/extractor-api/llm"
)

var (
	client = llm.NewClient()
)

func main() {
	ss := flag.Bool("server", false, "supply to start server")
	flag.Parse()

	if *ss {
		server := NewServer()
		log.Println("server listening on port 3000...")
		if err := http.ListenAndServe("localhost:3000", server); err != nil {
			log.Fatal(err)
		}
	}

}

type FindRequest struct {
	Url        string   `json:"url"`
	Attributes []string `json:"attributes"`
}

func decodeRequest(request *http.Request) (*FindRequest, error) {
	findRequest := FindRequest{}
	if err := json.NewDecoder(request.Body).Decode(&findRequest); err != nil {
		return nil, fmt.Errorf("invalid request data please provide url && attributes")
	} else {
		return &findRequest, nil
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
	mux := http.ServeMux{}
	methods := Methods{}
	mux.HandleFunc("/find", methods.POST(func(w http.ResponseWriter, r *http.Request) {
		bytesTokenLimit := 9216
		//parse the product information from the url
		req, err := decodeRequest(r)
		if err != nil {
			//write response
			common.EncodeResponse(err, w, 400)
			return
		}
		respBytes, err := fetchHTML(req.Url)
		if err != nil {
			//write response
			common.EncodeResponse(err, w, 500)
		}
		if len(respBytes) > bytesTokenLimit {
			respBytes = respBytes[:bytesTokenLimit]
		}
		prefix := "extract the main product details from the `page_text` and return them as key value pairs in json format"

		client.SendMsg(fmt.Sprintf("%s, page_text: %v", prefix, string(respBytes)), w)

	}))

	return &mux
}

type Methods struct{}

func checkMethod(method string, f func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.Error(w, "Method not allowed", http.StatusBadRequest)
			return
		}
		f(w, r)
	}
}
func (m Methods) POST(f func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return checkMethod("POST", f)
}
