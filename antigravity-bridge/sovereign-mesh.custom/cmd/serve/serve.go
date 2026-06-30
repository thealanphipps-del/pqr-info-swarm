package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "time"

    "github.com/pqr-info/sovereign-mesh/models/flashlite"
)

// RequestPayload defines the JSON payload for inference requests.
type RequestPayload struct {
    Input []float32 `json:"input"`
}

// ResponsePayload defines the JSON response from the model.
type ResponsePayload struct {
    Output []float32 `json:"output"`
    Error  string   `json:"error,omitempty"`
}

func main() {
    // Load configuration from environment variables (fallback to defaults).
    cfg := flashlite.Config{
        ModelPath: getEnv("FLASHLITE_MODEL_PATH", "./model/flashlite.bin"),
        Device:    getEnv("FLASHLITE_DEVICE", "cpu"),
        BatchSize: getEnvInt("FLASHLITE_BATCH_SIZE", 1),
    }

    model, err := flashlite.New(cfg)
    if err != nil {
        log.Fatalf("failed to initialize FlashLite model: %v", err)
    }
    log.Printf("FlashLite model loaded (path=%s, device=%s)", cfg.ModelPath, cfg.Device)

    http.HandleFunc("/api/v1/flashlite/predict", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
            return
        }
        var payload RequestPayload
        if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
            http.Error(w, "invalid JSON payload", http.StatusBadRequest)
            return
        }
        // Perform inference.
        out, err := model.Predict(payload.Input)
        resp := ResponsePayload{Output: out}
        if err != nil {
            resp.Error = err.Error()
        }
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(resp)
    })

    // Simple health endpoint.
    http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "ok")
    })

    // Server configuration.
    port := getEnv("FLASHLITE_PORT", "8080")
    srv := &http.Server{
        Addr:         ":" + port,
        ReadTimeout:  5 * time.Second,
        WriteTimeout: 10 * time.Second,
        IdleTimeout:  120 * time.Second,
    }

    log.Printf("FlashLite inference server listening on %s", srv.Addr)
    if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        log.Fatalf("server error: %v", err)
    }
}

// getEnv reads an environment variable or returns a default.
func getEnv(key, def string) string {
    if v := os.Getenv(key); v != "" {
        return v
    }
    return def
}

// getEnvInt reads an environment variable as int or returns a default.
func getEnvInt(key string, def int) int {
    if v := os.Getenv(key); v != "" {
        var i int
        if _, err := fmt.Sscanf(v, "%d", &i); err == nil {
            return i
        }
    }
    return def
}
