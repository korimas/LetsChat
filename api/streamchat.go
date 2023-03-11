package handler

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func StreamChatHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
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

		_, err = io.Copy(w, reader)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		return

	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

}
