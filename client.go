package pqr

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// Client provides a simple HTTP client for agents to interact with the ticketing system
type Client struct {
	BaseURL string
	Client  *http.Client
}

// NewClient creates a new ticketing client
func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// CreateTicket creates a new ticket (agent memory container)
func (c *Client) CreateTicket(ctx context.Context, subject, queue, content string, agentID string, intent map[string]interface{}) (string, error) {
	payload := map[string]interface{}{
		"Subject": subject,
		"Queue":   queue,
		"Text":    content,
		"AgentID": agentID,
		"Layer":   2,
		"Intent":  intent,
	}

	data, _ := json.Marshal(payload)
	req, _ := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/REST/2.0/ticket", c.BaseURL), bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("failed to create ticket: %v", result)
	}

	return result["id"].(string), nil
}

// StoreMemory stores agent context/memory for a ticket
func (c *Client) StoreMemory(ctx context.Context, agentID string, ticketID string, memType string, data map[string]interface{}, relevance float64) error {
	payload := map[string]interface{}{
		"memory_type":      memType,
		"data":             data,
		"relevance_score":  relevance,
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequestWithContext(ctx, "POST", 
		fmt.Sprintf("%s/REST/2.0/agent/%s/memory/%s", c.BaseURL, agentID, ticketID), 
		bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to store memory: %s", string(body))
	}

	return nil
}

// GetMemory retrieves agent memory for a ticket
func (c *Client) GetMemory(ctx context.Context, agentID string, ticketID string, memType string) (map[string]interface{}, error) {
	req, _ := http.NewRequestWithContext(ctx, "GET",
		fmt.Sprintf("%s/REST/2.0/agent/%s/memory/%s?type=%s", c.BaseURL, agentID, ticketID, memType),
		nil)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("memory not found")
	}

	var data map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&data)
	return data, nil
}

// GetContext retrieves all context tickets for an agent
func (c *Client) GetContext(ctx context.Context, agentID string) ([]map[string]interface{}, error) {
	req, _ := http.NewRequestWithContext(ctx, "GET",
		fmt.Sprintf("%s/REST/2.0/agent/%s/context", c.BaseURL, agentID),
		nil)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get context")
	}

	tickets := result["context_tickets"].([]interface{})
	var ticketList []map[string]interface{}
	for _, t := range tickets {
		ticketList = append(ticketList, t.(map[string]interface{}))
	}

	return ticketList, nil
}

// GetTicket retrieves a ticket and its content
func (c *Client) GetTicket(ctx context.Context, ticketID string) (map[string]interface{}, error) {
	req, _ := http.NewRequestWithContext(ctx, "GET",
		fmt.Sprintf("%s/REST/2.0/ticket/%s", c.BaseURL, ticketID),
		nil)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ticket not found")
	}

	var ticket map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&ticket)
	return ticket, nil
}

// LinkTickets creates a relationship between two tickets
func (c *Client) LinkTickets(ctx context.Context, parentID, childID string, relationType string, agentID string) error {
	payload := map[string]interface{}{
		"relationship_type": relationType,
		"agent_id":          agentID,
	}

	data, _ := json.Marshal(payload)
	req, _ := http.NewRequestWithContext(ctx, "POST",
		fmt.Sprintf("%s/REST/2.0/ticket/%s/link/%s", c.BaseURL, parentID, childID),
		bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to link tickets: %s", string(body))
	}

	return nil
}

// UpdateTicket updates a ticket status or title
func (c *Client) UpdateTicket(ctx context.Context, ticketID string, status string, title string) error {
	payload := map[string]interface{}{
		"Status": status,
		"Title":  title,
	}

	data, _ := json.Marshal(payload)
	req, _ := http.NewRequestWithContext(ctx, "PUT",
		fmt.Sprintf("%s/REST/2.0/ticket/%s", c.BaseURL, ticketID),
		bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to update ticket: %s", string(body))
	}

	return nil
}

// GetAuditTrail retrieves the audit trail for a ticket
func (c *Client) GetAuditTrail(ctx context.Context, ticketID string) ([]map[string]interface{}, error) {
	req, _ := http.NewRequestWithContext(ctx, "GET",
		fmt.Sprintf("%s/REST/2.0/ticket/%s/audit", c.BaseURL, ticketID),
		nil)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get audit trail")
	}

	trail := result["audit_trail"].([]interface{})
	var auditList []map[string]interface{}
	for _, entry := range trail {
		auditList = append(auditList, entry.(map[string]interface{}))
	}

	return auditList, nil
}

// Health checks if the service is running
func (c *Client) Health(ctx context.Context) (bool, error) {
	req, _ := http.NewRequestWithContext(ctx, "GET",
		fmt.Sprintf("%s/REST/2.0/health", c.BaseURL),
		nil)

	resp, err := c.Client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK, nil
}

// InitSchema initializes the database schema
func (c *Client) InitSchema(ctx context.Context) error {
	req, _ := http.NewRequestWithContext(ctx, "POST",
		fmt.Sprintf("%s/REST/2.0/init", c.BaseURL),
		nil)

	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to init schema: %s", string(body))
	}

	return nil
}

// AgentSession provides a high-level interface for agents
type AgentSession struct {
	client  *Client
	agentID string
	tickets []uuid.UUID
}

// NewAgentSession creates a session for an agent
func NewAgentSession(baseURL, agentID string) *AgentSession {
	return &AgentSession{
		client:  NewClient(baseURL),
		agentID: agentID,
	}
}

// CreateMemory creates a ticket and stores initial memory
func (as *AgentSession) CreateMemory(ctx context.Context, subject string, content map[string]interface{}) (string, error) {
	ticketID, err := as.client.CreateTicket(ctx, subject, "DEFAULT", "", as.agentID, content)
	if err != nil {
		return "", err
	}

	// Store the memory
	if err := as.client.StoreMemory(ctx, as.agentID, ticketID, "context", content, 1.0); err != nil {
		return "", err
	}

	return ticketID, nil
}

// RecallMemory retrieves memory for this agent
func (as *AgentSession) RecallMemory(ctx context.Context, ticketID string) (map[string]interface{}, error) {
	return as.client.GetMemory(ctx, as.agentID, ticketID, "context")
}

// GetAllMemories retrieves all context tickets for this agent
func (as *AgentSession) GetAllMemories(ctx context.Context) ([]map[string]interface{}, error) {
	return as.client.GetContext(ctx, as.agentID)
}


