package service

import (
	"math"
	"time"
)

type HardwareMetrics struct {
	CPUPercent  float64
	RAMUsedGB   float64
	RAMTotalGB  float64
	GPU0Percent float64 // NVIDIA
	GPU1Percent float64 // Intel
	GPU2Percent float64 // Yoga NPU (logical "GPU2")
	Timestamp   time.Time
}

func CollectHardwareMetrics() (HardwareMetrics, error) {
	t := float64(time.Now().UnixNano() / 1e9)
	cpuLoad := 30 + 15*math.Sin(t/5)
	ramUsed := 16.0 + 2.0*math.Cos(t/20)

	return HardwareMetrics{
		CPUPercent:  cpuLoad,
		RAMUsedGB:   ramUsed,
		RAMTotalGB:  64.0,
		GPU0Percent: mockGPU0Percent(),
		GPU1Percent: mockGPU1Percent(),
		GPU2Percent: mockGPU2Percent(), // NPU mock
		Timestamp:   time.Time{},       // replaced by monitoring_service
	}, nil
}

func mockGPU0Percent() float64 {
	// Simulate a heavy workload for NVIDIA
	t := float64(time.Now().UnixNano() / 1e9)
	return 60 + 30*math.Sin(t/8)
}

func mockGPU1Percent() float64 {
	// Simulate a lighter/oscillating workload for Intel
	t := float64(time.Now().UnixNano() / 1e9)
	return 20 + 15*math.Cos(t/4)
}

func mockGPU2Percent() float64 {
	// Deterministic oscillation for NPU sparkline
	// Replace with real NPU utilization later.
	return 30.0
}
