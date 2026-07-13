package controller
import (
    "fmt"
    "gemini_testify/internal/model"
    "time"
)
type VitalityMonitor struct {
    Anchor     model.VitalityAnchor
    LastCount  int
    StartTime  time.Time
}
func (m *VitalityMonitor) Start(getCommitCount func() int) {
    m.StartTime = time.Now()
    ticker := time.NewTicker(1 * time.Minute)
    for range ticker.C {
        currentCount := getCommitCount()
        diff := currentCount - m.LastCount
        m.LastCount = currentCount
        m.Anchor.CalculateSlope(diff, 1*time.Minute)
        if m.Anchor.Slope < 5.0 {
            msg := fmt.Sprintf("[VITALITY_ALERT] SLOPE_DECAY: %.2f CPM", m.Anchor.Slope)
            fmt.Println(msg)
        }
    }
}
