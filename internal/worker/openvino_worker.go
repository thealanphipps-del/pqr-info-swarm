package worker

type InferenceBackend interface {
	LoadModel(path string) error
	Infer(input []float32) ([]float32, error)
}

type OpenVINOWorker struct {
	device string
	// core *ov.Core // real field when go-openvino is available
}

func NewOpenVINOWorker(device string) (*OpenVINOWorker, error) {
	return &OpenVINOWorker{device: device}, nil
}

func (w *OpenVINOWorker) LoadModel(path string) error {
	// TODO: wire to ov.Core when dependency is available
	return nil
}

func (w *OpenVINOWorker) Infer(input []float32) ([]float32, error) {
	// TODO: real inference; for now, echo or simple transform
	return input, nil
}
