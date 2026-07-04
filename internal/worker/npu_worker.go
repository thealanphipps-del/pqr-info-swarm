package worker

type NPUBackend interface {
	LoadModel(path string) error
	Infer(input []float32) ([]float32, error)
}

type NPUWorker struct {
	device string // e.g. "npu", "cpu-fallback"
	// later: real onnxruntime session
}

func NewNPUWorker(device string) (*NPUWorker, error) {
	return &NPUWorker{device: device}, nil
}

func (w *NPUWorker) LoadModel(path string) error {
	// TODO: wire to onnxruntime-directml when available
	return nil
}

func (w *NPUWorker) Infer(input []float32) ([]float32, error) {
	// Stub: simple passthrough or tiny transform
	// This keeps the mesh wiring and scheduler logic testable.
	return input, nil
}
