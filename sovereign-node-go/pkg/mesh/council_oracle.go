package mesh

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"
)

// -------------------------------
// Oracle Request / Response Types
// -------------------------------

type CouncilAuditRequest struct {
    CouncilMembers []int            `json:"council_members"`
    Wealth         []float64        `json:"wealth"`
    Confidence     []float64        `json:"confidence"`
    Phenotypes     []string         `json:"phenotypes"`
    TopologyDelta  float64          `json:"topology_divergence"`
    TurnoverRate   float64          `json:"turnover_rate"`
}

type CouncilAuditResponse struct {
    Stability      string `json:"stability"`      // "stable" | "unstable"
    Recommendation string `json:"recommendation"` // "approve" | "delay" | "reshuffle"
    Notes          string `json:"notes"`
}

// -------------------------------
// Internal HTTP Client
// -------------------------------

var oracleHTTP = &http.Client{
    Timeout: 2500 * time.Millisecond,
}

const oracleURL = "http://localhost:11434/api/generate"

// -------------------------------
// Core Function: AskOracleForCouncilAudit
// -------------------------------

func AskOracleForCouncilAudit(req CouncilAuditRequest) CouncilAuditResponse {
    // Marshal request JSON
    payload, err := json.Marshal(map[string]interface{}{
        "model":  "sovereign-oracle",
        "prompt": buildCouncilAuditPrompt(req),
        "format": "json",
    })
    if err != nil {
        return fallbackCouncilAudit("marshal_error")
    }

    // Send request
    httpReq, err := http.NewRequest("POST", oracleURL, bytes.NewBuffer(payload))
    if err != nil {
        return fallbackCouncilAudit("request_error")
    }
    httpReq.Header.Set("Content-Type", "application/json")

    resp, err := oracleHTTP.Do(httpReq)
    if err != nil {
        return fallbackCouncilAudit("timeout_or_network_error")
    }
    defer resp.Body.Close()

    // Read response
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return fallbackCouncilAudit("read_error")
    }

    // Parse JSON
    var parsed struct {
        Response string `json:"response"`
    }
    if err := json.Unmarshal(body, &parsed); err != nil {
        return fallbackCouncilAudit("unmarshal_error")
    }

    // Parse the Oracle's JSON inside the "response" field
    var audit CouncilAuditResponse
    if err := json.Unmarshal([]byte(parsed.Response), &audit); err != nil {
        return fallbackCouncilAudit("schema_error")
    }

    // Validate fields
    if audit.Stability == "" || audit.Recommendation == "" {
        return fallbackCouncilAudit("missing_fields")
    }

    return audit
}

// -------------------------------
// Core Function: AskOracleForCouncilRecommendation
// -------------------------------

func AskOracleForCouncilRecommendation(req CouncilAuditRequest) string {
    audit := AskOracleForCouncilAudit(req)
    return audit.Recommendation
}

// -------------------------------
// Prompt Builder
// -------------------------------

func buildCouncilAuditPrompt(req CouncilAuditRequest) string {
    b, _ := json.Marshal(req)
    return fmt.Sprintf(`
You are the Sovereign Mesh Oracle.

Evaluate the following Council governance snapshot and return STRICT JSON ONLY.

%s

Return JSON:
{
  "stability": "stable | unstable",
  "recommendation": "approve | delay | reshuffle",
  "notes": "short machine-readable rationale"
}
`, string(b))
}

// -------------------------------
// Fallback Logic
// -------------------------------

func fallbackCouncilAudit(reason string) CouncilAuditResponse {
    return CouncilAuditResponse{
        Stability:      "stable",
        Recommendation: "approve",
        Notes:          "fallback_" + reason,
    }
}
