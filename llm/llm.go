package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/aodr3w/extractor-api/common"
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

func (c *LLMClient) SendMsg(msg string, w http.ResponseWriter) {
	//this message should stream the data back to the client
	body := map[string]interface{}{
		"model": c.model,
		"messages": []map[string]string{
			{"role": "user", "content": msg},
		},
	}

	payload, err := json.Marshal(body)
	if err != nil {
		common.EncodeResponse([]byte(fmt.Sprintf("error marshalling body %v", err)), w, http.StatusInternalServerError)
		return
	}
	response, err := http.Post(
		"http://localhost:11434/api/chat",
		"application/json",
		bytes.NewBuffer(payload),
	)
	if err != nil {
		common.EncodeResponse(fmt.Sprintf("error sending llm request: %v\n", err), w, http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	var buffer bytes.Buffer
	_, err = io.Copy(&buffer, response.Body)
	if err != nil {
		common.EncodeResponse(fmt.Sprintf("error reading LLM response: %v", err), w, http.StatusInternalServerError)
		return
	}

	common.EncodeResponse(buffer.String(), w, http.StatusOK)
}
