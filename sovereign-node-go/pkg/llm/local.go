package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Config holds the local LM Studio / gateway configuration
type Config struct {
	BaseURL        string
	CompletionName string
	EmbeddingName  string
}

// DefaultConfig points to a standard local LM Studio setup
var DefaultConfig = Config{
	BaseURL:        "http://127.0.0.1:1234/v1", // Default LM Studio API port
	CompletionName: "gemma-4-e4b",
	EmbeddingName:  "text-embedding-nomic-embed-text-v1.5",
}

// LocalGateway provides methods to interface with the local LLM stack
type LocalGateway struct {
	client *http.Client
	config Config
}

func NewLocalGateway(cfg Config) *LocalGateway {
	return &LocalGateway{
		client: &http.Client{Timeout: 300 * time.Second},
		config: cfg,
	}
}

// CompletionRequest represents an OpenAI-compatible chat completion request
type CompletionRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Temp     float64   `json:"temperature,omitempty"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// CompletionResponse represents an OpenAI-compatible chat completion response
type CompletionResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

// GenerateSummary uses the local completion model to process swarm telemetry and tickets
func (g *LocalGateway) GenerateSummary(ctx context.Context, prompt string) (string, error) {
	// First-layer: Attempt to use LM Studio (configured Local LLM)
	reqBody := CompletionRequest{
		Model: g.config.CompletionName,
		Messages: []Message{
			{Role: "system", Content: "You are an autonomous hyperdevelopment agent summarizing Sovereign Node swarm data. Be concise and forensic."},
			{Role: "user", Content: prompt},
		},
		Temp: 0.2,
	}

	jsonData, err := json.Marshal(reqBody)
	if err == nil {
		req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/chat/completions", g.config.BaseURL), bytes.NewBuffer(jsonData))
		if err == nil {
			req.Header.Set("Content-Type", "application/json")
			resp, err := g.client.Do(req)
			if err == nil {
				defer resp.Body.Close()
				if resp.StatusCode == http.StatusOK {
					var compResp CompletionResponse
					if err := json.NewDecoder(resp.Body).Decode(&compResp); err == nil && len(compResp.Choices) > 0 {
						return compResp.Choices[0].Message.Content, nil
					}
				}
			}
		}
	}

	// Second-layer fallback: Attempt to use local Ollama with sovereign-oracle
	ollamaURL := "http://127.0.0.1:11434/v1/chat/completions"
	reqBodyOllama := CompletionRequest{
		Model: "sovereign-oracle",
		Messages: []Message{
			{Role: "system", Content: "You are an autonomous hyperdevelopment agent summarizing Sovereign Node swarm data. Be concise and forensic."},
			{Role: "user", Content: prompt},
		},
		Temp: 0.2,
	}

	jsonDataOllama, err := json.Marshal(reqBodyOllama)
	if err == nil {
		req, err := http.NewRequestWithContext(ctx, "POST", ollamaURL, bytes.NewBuffer(jsonDataOllama))
		if err == nil {
			req.Header.Set("Content-Type", "application/json")
			resp, err := g.client.Do(req)
			if err == nil {
				defer resp.Body.Close()
				if resp.StatusCode == http.StatusOK {
					var compResp CompletionResponse
					if err := json.NewDecoder(resp.Body).Decode(&compResp); err == nil && len(compResp.Choices) > 0 {
						return compResp.Choices[0].Message.Content, nil
					}
				}
			}
		}
	}

	return "", fmt.Errorf("both primary LM Studio and fallback Ollama inference layers failed")
}

// EmbeddingRequest represents an OpenAI-compatible embedding request
type EmbeddingRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

// EmbeddingResponse represents an OpenAI-compatible embedding response
type EmbeddingResponse struct {
	Data []struct {
		Embedding []float32 `json:"embedding"`
	} `json:"data"`
}

// GenerateEmbedding uses the specialized Nomic model to generate vector embeddings for semantic search
func (g *LocalGateway) GenerateEmbedding(ctx context.Context, text string) ([]float32, error) {
	reqBody := EmbeddingRequest{
		Model: g.config.EmbeddingName,
		Input: text,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/embeddings", g.config.BaseURL), bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("local embedding error: status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var embResp EmbeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&embResp); err != nil {
		return nil, err
	}

	if len(embResp.Data) > 0 {
		return embResp.Data[0].Embedding, nil
	}

	return nil, fmt.Errorf("no embedding data returned")
}
