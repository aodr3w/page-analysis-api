package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/aodr3w/page-analysis-api/data"
)

type LLMClient struct {
	model string
}

type LLMResponse struct {
	Model     string            `json:"model"`
	CreatedAt time.Time         `json:"created_at"`
	Message   map[string]string `json:"message"`
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
		data.EncodeResponse([]byte(fmt.Sprintf("error marshalling body %v", err)), w, http.StatusInternalServerError)
		return
	}
	response, err := http.Post(
		"http://localhost:11434/api/chat",
		"application/json",
		bytes.NewBuffer(payload),
	)
	if err != nil {
		data.EncodeResponse(fmt.Sprintf("error sending llm request: %v\n", err), w, http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	var result strings.Builder
	decoder := json.NewDecoder(response.Body)
	for {
		llmResponse := LLMResponse{}
		if err := decoder.Decode(&llmResponse); err == io.EOF {
			break
		} else if err != nil {
			data.EncodeResponse(fmt.Sprintf("Error decoding LLM response: %v", err), w, http.StatusInternalServerError)
			return
		}
		if content, exists := llmResponse.Message["content"]; exists {
			result.WriteString(content)
		}
	}
	data.EncodeResponse(result.String(), w, http.StatusOK)
}
