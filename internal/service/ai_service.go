package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// AIService provides a unified interface for agentic model access
type AIService struct {
	GemmaURL    string
	LMStudioURL string
}

func NewAIService() *AIService {
	gemma := os.Getenv("GEMMA_ENDPOINT")
	if gemma == "" {
		gemma = "http://192.168.12.169:11434"
	}
	
	lm := os.Getenv("LMSTUDIO_ENDPOINT")
	if lm == "" {
		lm = "http://192.168.12.236:1234"
	}

	return &AIService{
		GemmaURL:    gemma,
		LMStudioURL: lm,
	}
}

// QuerySwarm attempts to resolve a query using the available AI mesh (Ollama -> LM Studio)
func (a *AIService) QuerySwarm(ctx context.Context, prompt string) (string, string, error) {
	// Try Primary (Ollama/Gemma)
	resp, err := a.QueryGemma(ctx, "gemma2:2b", prompt)
	if err == nil && resp != "" {
		return resp, "Ollama-169", nil
	}

	// Fallback to Local (LM Studio)
	resp, err = a.QueryLMStudio(ctx, "gemma-2-9b-it", prompt)
	if err == nil && resp != "" {
		return resp, "LM-Studio-Local", nil
	}

	return "", "", fmt.Errorf("all AI nodes are offline")
}

func (a *AIService) QueryGemma(ctx context.Context, model, prompt string) (string, error) {
	reqBody, _ := json.Marshal(map[string]interface{}{
		"model": model,
		"messages": []map[string]interface{}{
			{"role": "user", "content": prompt},
		},
		"stream": false,
	})

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Post(a.GemmaURL+"/api/chat", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	}
	json.NewDecoder(resp.Body).Decode(&result)
	return result.Message.Content, nil
}

func (a *AIService) QueryLMStudio(ctx context.Context, model, prompt string) (string, error) {
	// Using LM Studio v1 REST API
	reqBody, _ := json.Marshal(map[string]interface{}{
		"model": model,
		"messages": []map[string]interface{}{
			{"role": "user", "content": prompt},
		},
		"stream": false,
	})

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Post(a.LMStudioURL+"/api/v1/chat", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Parse the v1 API response
	var result struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode v1 API response: %v", err)
	}

	if result.Message.Content == "" {
		return "", fmt.Errorf("empty response from LM Studio v1 API")
	}

	return result.Message.Content, nil
}
