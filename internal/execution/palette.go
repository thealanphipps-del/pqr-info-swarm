package execution

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/thealanphipps-del/pqr/internal/domain"
)

// Define interface matching db.CockroachRepository capabilities
type MemoryRepository interface {
	ListTopFixes() ([]domain.ErrorSolutionRecord, error)
}

func StartPalette(engine *ExecutionEngine, memory MemoryRepository, handleGoback func()) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println()
		fmt.Println("=== SWEN Command Palette ===")
		fmt.Println("1) Start Emulator Workflow")
		fmt.Println("2) Diagnose GCP Auth")
		fmt.Println("3) View Known Fixes")
		fmt.Println("4) Goback (Time Machine)")
		fmt.Println("5) Exit")
		fmt.Print("Select option: ")

		line, _ := reader.ReadString('\n')
		choice := strings.TrimSpace(line)

		switch choice {
		case "1":
			runEmulatorWorkflow(engine)
		case "2":
			runGcpDiagnosis(engine)
		case "3":
			showKnownFixes(memory)
		case "4":
			if handleGoback != nil {
				handleGoback()
			}
		case "5":
			fmt.Println("Exiting palette.")
			return
		default:
			fmt.Println("Unknown option.")
		}
	}
}

func runEmulatorWorkflow(engine *ExecutionEngine) {
	fmt.Println("[palette] launching emulator workflow...")
	// You can later prompt for AVD name, etc.
	res, err := engine.LaunchEmulator(
		ctxBackground(),
		&WorkflowDescriptor{
			Kind:    WorkflowKindEmulator,
			Command: "emulator",
			Args:    []string{"-list-avds"},
		},
	)
	printResult(res, err)
}

func runGcpDiagnosis(engine *ExecutionEngine) {
	fmt.Println("[palette] diagnosing GCP auth...")
	res, err := engine.DiagnoseGCP(ctxBackground())
	printResult(res, err)
}

func showKnownFixes(memory MemoryRepository) {
	fmt.Println("[palette] known fixes:")
	fixes, err := memory.ListTopFixes()
	if err != nil {
		fmt.Printf("error listing fixes: %v\n", err)
		return
	}
	for _, f := range fixes {
		fmt.Printf("- %s (success_rate=%.2f)\n", f.SignatureHash, f.SuccessRate)
	}
}

func printResult(res *ExecutionResult, err error) {
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	for _, line := range res.Logs {
		fmt.Println(line)
	}
}

func ctxBackground() context.Context {
	return context.Background()
}
