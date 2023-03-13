package bak

import (
	"bufio"
	"bytes"
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type Choice struct {
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
	Index        int     `json:"index"`
}

type ErrorResp struct {
	Type string `json:"type"`
}

type ChatResponse struct {
	ID      string    `json:"id"`
	Object  string    `json:"object"`
	Created int       `json:"created"`
	Model   string    `json:"model"`
	Usage   Usage     `json:"usage"`
	Choices []Choice  `json:"choices"`
	Error   ErrorResp `json:"error"`
}

type APIResponse struct {
	Success bool        `json:"success"`
	Result  interface{} `json:"result"`
	Message string      `json:"message"`
}

var httpClient = http.Client{}

// var API_KEY = "Bearer " + os.Getenv("API_KEY")
var API_KEY = "Bearer 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"
var doneSequence = []byte("[DONE]")
var waitSequence = []byte("data: [WAIT]")

func StreamChatHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var mutex sync.Mutex
		received := false

		request := ChatRequest{}
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		messageLen := len(request.Messages)
		if messageLen > 10 {
			request.Messages = request.Messages[messageLen-10:]
		}
		request.Stream = true

		reqByte, err := json.Marshal(request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewReader(reqByte))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", API_KEY)

		resp, err := httpClient.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		reader := bufio.NewReader(resp.Body)
		defer resp.Body.Close()

		// diff way
		go func() {
			for {
				mutex.Lock()
				if !received {
					w.Write(waitSequence)
				} else {
					mutex.Unlock()
					return
				}
				mutex.Unlock()
				time.Sleep(time.Second * 3)
			}
		}()

		for {
			line, err := reader.ReadBytes('\n')
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if !received {
				mutex.Lock()
				received = true
				mutex.Unlock()
			}

			w.Write(line)

			if bytes.HasPrefix(line, doneSequence) {
				break
			}
		}

		//_, err = io.Copy(w, reader)
		//if err != nil {
		//	http.Error(w, err.Error(), http.StatusInternalServerError)
		//	return
		//}

		return

	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

}
