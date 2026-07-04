package worker

type InferenceTask struct {
	Model  string
	Input  []float32
	Device string // "intel", "nvidia", "npu", "auto"
}

type Scheduler struct {
	intel  InferenceBackend
	nvidia InferenceBackend
	npu    NPUBackend
}

func NewScheduler(intel, nvidia InferenceBackend, npu NPUBackend) *Scheduler {
	return &Scheduler{
		intel:  intel,
		nvidia: nvidia,
		npu:    npu,
	}
}

func (s *Scheduler) Dispatch(t InferenceTask) ([]float32, error) {
	switch t.Device {
	case "intel":
		return s.intel.Infer(t.Input)
	case "nvidia":
		return s.nvidia.Infer(t.Input)
	case "npu":
		return s.npu.Infer(t.Input)
	case "auto":
		// Simple policy for now: prefer NPU for light tasks
		return s.npu.Infer(t.Input)
	default:
		return s.npu.Infer(t.Input)
	}
}
