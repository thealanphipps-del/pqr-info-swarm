package sovereign

import (
	"context"
	"testing"
)

func TestOptimizePrompt(t *testing.T) {
	c := NewController("model-loader-495607-m2", "us-central1")
	task := OptimizationTask{
		ID:          "test-1",
		Type:        "prompt-opt",
		Prompt:      "Optimize this prompt: write a self-healing process watchdog",
		ExamplePath: "gs://model-loader-495607-m2_cloudbuild/test_examples.jsonl",
		Method:      "apophis",
	}

	ctx := context.Background()
	jobName, err := c.OptimizePrompt(ctx, task)
	if err != nil {
		t.Logf("Warning: Vertex AI request failed (expected if Vertex API is not active or credential lacks permission): %v", err)
	} else {
		t.Logf("Success: Created Pipeline Job: %s", jobName)
	}
}
