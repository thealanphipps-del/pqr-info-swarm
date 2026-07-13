package model
import "time"
type VitalityAnchor struct {
    Slope float64
}
func (v *VitalityAnchor) CalculateSlope(diff int, duration time.Duration) {
    v.Slope = float64(diff) / duration.Minutes()
}
