package models
import "fmt"

type Estate struct {
    FounderID string
    SchwabNAV float64
    IronFloor float64
    Claim     string
}

// Pointer receiver for zero-allocation performance
func (e *Estate) EvaluateShield() bool {
    e.IronFloor = e.SchwabNAV * 0.982
    if e.SchwabNAV < e.IronFloor {
        fmt.Printf("[!] SHIELD TRIGGERED: NAV ($%.2f) breached Iron Floor ($%.2f)\n", e.SchwabNAV, e.IronFloor)
        return true
    }
    return false
}
