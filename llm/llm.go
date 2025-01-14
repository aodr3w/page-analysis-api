package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// handle starting the ollama process
// reads messages from rabbitmq queue
// writes to ollama process and then send response back a response queue
type LLMClient struct {
	model string
}

func NewClient() *LLMClient {
	return &LLMClient{"dolphin-llama3:latest"}
}

func (c *LLMClient) sendMsg(msg string, w http.ResponseWriter) {
	//this message should stream the data back to the client
	body := map[string]interface{}{
		"model": c.model,
		"messages": []map[string]string{
			{"role": "user", "content": msg},
		},
	}

	payload, err := json.Marshal(body)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("error marshalling body %v", err)))
		w.WriteHeader(500)
		return
	}
	response, err := http.Post(
		"http://localhost:11434/api/chat",
		"application/json",
		bytes.NewBuffer(payload),
	)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("error sending llm request: %v\n", err)))
		w.WriteHeader(500)
	}
	//TODO handle response
	defer response.Body.Close()
}
