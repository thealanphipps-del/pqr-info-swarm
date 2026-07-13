package controllers
import "fmt"

type HUD struct {
    VerbalEnabled bool
}

func (h *HUD) Coach(message string) {
    if h.VerbalEnabled {
        fmt.Printf("[*] HUD Voice Coach: \"%s\"\n", message)
    }
}
