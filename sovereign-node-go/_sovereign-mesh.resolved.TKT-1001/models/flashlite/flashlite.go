package flashlite

import (
    "fmt"
    "sync"
)

// Config holds configuration for the Flash Lite model.
type Config struct {
    // Path to the serialized model file (e.g., .bin, .pt, .onnx).
    ModelPath string `json:"model_path"`
    // Device to run on: "cpu" or "gpu". Defaults to "cpu".
    Device string `json:"device"`
    // Optional batch size for inference.
    BatchSize int `json:"batch_size"`
}

// Model is a lightweight wrapper around the Flash Lite inference engine.
// This implementation is a stub that pretends to load a model and perform
// inference. Replace the stub with actual calls to the underlying library
// (e.g., Gorgonia, TensorFlow Go bindings) when the real model is available.
type Model struct {
    cfg   Config
    mu    sync.Mutex // protects any internal state during concurrent calls
    ready bool
}

// New creates a new Model instance and loads the model according to cfg.
func New(cfg Config) (*Model, error) {
    m := &Model{cfg: cfg}
    if err := m.load(); err != nil {
        return nil, err
    }
    return m, nil
}

// load simulates model loading. In a real implementation, this would read the
// model file and initialise the inference runtime.
func (m *Model) load() error {
    m.mu.Lock()
    defer m.mu.Unlock()
    if m.cfg.ModelPath == "" {
        return fmt.Errorf("model_path is required")
    }
    // Stub: just mark as ready.
    m.ready = true
    return nil
}

// Predict runs inference on the provided input tensor and returns an output.
// The input is a slice of float32 representing a flattened tensor. The stub
// simply echoes the input.
func (m *Model) Predict(input []float32) ([]float32, error) {
    m.mu.Lock()
    defer m.mu.Unlock()
    if !m.ready {
        return nil, fmt.Errorf("model not loaded")
    }
    // Stub behavior: copy input to output.
    output := make([]float32, len(input))
    copy(output, input)
    return output, nil
}
