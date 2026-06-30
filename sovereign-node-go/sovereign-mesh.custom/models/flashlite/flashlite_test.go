package flashlite

import (
    "reflect"
    "testing"
)

func TestNewModelLoadsSuccessfully(t *testing.T) {
    cfg := Config{ModelPath: "dummy.bin", Device: "cpu", BatchSize: 1}
    m, err := New(cfg)
    if err != nil {
        t.Fatalf("expected no error, got %v", err)
    }
    if m == nil {
        t.Fatalf("model instance is nil")
    }
    if !m.ready {
        t.Fatalf("model not marked as ready after load")
    }
}

func TestPredictEchoesInput(t *testing.T) {
    cfg := Config{ModelPath: "dummy.bin", Device: "cpu", BatchSize: 1}
    m, err := New(cfg)
    if err != nil {
        t.Fatalf("failed to create model: %v", err)
    }
    input := []float32{0.1, 0.2, 0.3, 0.4}
    output, err := m.Predict(input)
    if err != nil {
        t.Fatalf("predict returned error: %v", err)
    }
    if !reflect.DeepEqual(output, input) {
        t.Fatalf("expected output %v, got %v", input, output)
    }
}
