package adapter

import (
	"context"
	"net/http"
	"time"
)

type SwenClient struct {
	httpClient   *http.Client
	swarmAPIURL  string
	ticketAPIURL string
}

func NewSwenClient(swarmAPIURL, ticketAPIURL string) *SwenClient {
	return &SwenClient{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		swarmAPIURL:  swarmAPIURL,
		ticketAPIURL: ticketAPIURL,
	}
}

type ProxyEvent struct {
	Proxy    string `json:"proxy"`
	URL      string `json:"url"`
	Method   string `json:"method"`
	ReqBody  string `json:"reqBody"`
	RespBody string `json:"respBody"`
	Status   int    `json:"status"`
}

type ChatCompletion struct {
	Choices []struct {
		Message struct {
			Role             string `json:"role"`
			Content          string `json:"content"`
			ReasoningContent string `json:"reasoning_content"`
		} `json:"message"`
	} `json:"choices"`
}

// TODO: flesh these out with real REST 2.0 calls.

func (c *SwenClient) FetchSwarmStream(ctx context.Context) ([]string, error) {
	return []string{"swarm stream placeholder"}, nil
}

func (c *SwenClient) FetchTimeline(ctx context.Context) ([]string, error) {
	return []string{"timeline placeholder"}, nil
}

func (c *SwenClient) SendAgentMessage(ctx context.Context, agentID, msg string) error {
	return nil
}

func (c *SwenClient) ExecuteCommand(ctx context.Context, cmd string) (string, error) {
	return "command executed: " + cmd, nil
}
