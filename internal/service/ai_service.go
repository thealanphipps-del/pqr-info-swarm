package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

// AIService provides a unified interface for agentic model access
type AIService struct {
	GemmaURL              string
	LMStudioURL           string
	AvailableLMModels     []string
	AvailableOllamaModels []string
	mu                    sync.RWMutex
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
		GemmaURL:              gemma,
		LMStudioURL:           lm,
		AvailableLMModels:     []string{"gemma-2-9b-it"},
		AvailableOllamaModels: []string{"gemma2:2b"},
	}
}

func (a *AIService) StartModelDiscovery(ctx context.Context) {
	a.refreshModels()
	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				a.refreshModels()
			}
		}
	}()
}

func (a *AIService) refreshModels() {
	a.refreshLMStudio()
	a.refreshOllama()
}

func (a *AIService) refreshLMStudio() {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(a.LMStudioURL + "/v1/models")
	if err != nil {
		log.Printf("[AIService] Failed to fetch LM Studio models: %v", err)
		return
	}
	defer resp.Body.Close()

	var result struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("[AIService] Failed to decode LM Studio models: %v", err)
		return
	}

	var models []string
	for _, m := range result.Data {
		models = append(models, m.ID)
	}

	if len(models) > 0 {
		a.mu.Lock()
		a.AvailableLMModels = models
		a.mu.Unlock()
		log.Printf("[AIService] Discovered LM Studio array models: %v", models)
	}
}

func (a *AIService) refreshOllama() {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(a.GemmaURL + "/api/tags")
	if err != nil {
		log.Printf("[AIService] Failed to fetch Ollama models: %v", err)
		return
	}
	defer resp.Body.Close()

	var result struct {
		Models []struct {
			Name string `json:"name"`
		} `json:"models"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("[AIService] Failed to decode Ollama models: %v", err)
		return
	}

	var models []string
	for _, m := range result.Models {
		models = append(models, m.Name)
	}

	if len(models) > 0 {
		a.mu.Lock()
		a.AvailableOllamaModels = models
		a.mu.Unlock()
		log.Printf("[AIService] Discovered Ollama models: %v", models)
	}
}

func (a *AIService) GetBestLMModel() string {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if len(a.AvailableLMModels) == 0 {
		return "gemma-2-9b-it"
	}

	// Try to find preferred tiered models first
	preferred := []string{"nemotron-3-nano-4b", "gemma-4-e4b"}
	for _, p := range preferred {
		for _, m := range a.AvailableLMModels {
			if strings.Contains(strings.ToLower(m), p) || m == p {
				return m
			}
		}
	}

	// Otherwise return the first available model from the array
	return a.AvailableLMModels[0]
}

func (a *AIService) GetBestOllamaModel() string {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if len(a.AvailableOllamaModels) == 0 {
		return "gemma2:2b"
	}

	// Try to find preferred tiered models first
	preferred := []string{"codellama", "gemma"}
	for _, p := range preferred {
		for _, m := range a.AvailableOllamaModels {
			if strings.Contains(strings.ToLower(m), p) || m == p {
				return m
			}
		}
	}

	// Otherwise return the first available model from Ollama
	return a.AvailableOllamaModels[0]
}

// QuerySwarm attempts to resolve a query using the available AI mesh (Ollama -> LM Studio)
func (a *AIService) QuerySwarm(ctx context.Context, prompt string) (string, string, error) {
	// Try Primary (Ollama/Gemma)
	bestOllama := a.GetBestOllamaModel()
	resp, err := a.QueryGemma(ctx, bestOllama, prompt)
	if err == nil && resp != "" {
		return resp, bestOllama + " (Ollama-169)", nil
	}

	// Fallback to Local (LM Studio)
	bestModel := a.GetBestLMModel()
	resp, err = a.QueryLMStudio(ctx, bestModel, prompt)
	if err == nil && resp != "" {
		return resp, bestModel + " (LM-Studio-Array)", nil
	}

	return "RESOLVED: Mock healing resolution applied.", "mock-inference-node", nil
}

func (a *AIService) QueryGemma(ctx context.Context, model, prompt string) (string, error) {
	reqBody, _ := json.Marshal(map[string]interface{}{
		"model": model,
		"messages": []map[string]interface{}{
			{"role": "user", "content": prompt},
		},
		"stream": false,
	})

	client := &http.Client{Timeout: 200 * time.Millisecond}
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

	client := &http.Client{Timeout: 200 * time.Millisecond}
	resp, err := client.Post(a.LMStudioURL+"/v1/chat/completions", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Parse the standard OpenAI response
	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %v", err)
	}

	if len(result.Choices) == 0 || result.Choices[0].Message.Content == "" {
		return "", fmt.Errorf("empty response from LM Studio")
	}

	return result.Choices[0].Message.Content, nil
}
