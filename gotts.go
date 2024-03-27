package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// JSONModel represents each model in the JSON response
type JSONModel struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	OwnedBy string `json:"owned_by"`
}

type TextToSpeechRequest struct {
	Model string  `json:"model"`
	Input string  `json:input`
	Voice string  `json:voice`
	Speed float32 `json:speed`
}

func getAvailableModels(openAIAPIKey string) ([]JSONModel, error) {
	req, err := http.NewRequest("GET", "https://api.openai.com/v1/models", nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Add("Authorization", "Bearer "+openAIAPIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	fmt.Println("Response Status Code:", resp.StatusCode)

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	// Unmarshal JSON into slice of Model structs
	var modelsData struct {
		Data []JSONModel `json:"data"`
	}
	if err := json.Unmarshal(body, &modelsData); err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	return modelsData.Data, nil
}

func convertTextToSpeech(apiKey string, text string, filePath string) error {
	const (
		ttsURL            = "https://api.openai.com/v1/audio/speech"
		ttsModel          = "tts-1"
		ttsVoice          = "onyx"
		ttsSpeed          = 1.0
		ttsResponseFormat = "mp3"
	)
	reqBody := TextToSpeechRequest{
		Model: ttsModel,
		Input: text,
		Voice: ttsVoice,
		Speed: ttsSpeed,
	}
	// Marshall JSON Body
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("Could not marshall request body to JSON: %v", err)
	}
	req, err := http.NewRequest("POST", ttsURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return fmt.Errorf("Could not create request: %v", err)
	}

	req.Header.Add("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	// Send HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Unexpected Status Code: %v", resp.StatusCode)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("Could not create file: %v", err)
	}
	defer file.Close()

	// Handle response here if needed

	return nil

}

func main() {
	openAIAPIKey := os.Getenv("OPENAI_API_KEY")
	if openAIAPIKey == "" {
		fmt.Println("OPENAI_API_KEY environment variable is not set.")
		return
	}

	models, err := getAvailableModels(openAIAPIKey)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print the models
	for _, model := range models {
		fmt.Printf("Model ID: %s, Object: %s, Created: %d, OwnedBy: %s\n", model.ID, model.Object, model.Created, model.OwnedBy)
	}

}
