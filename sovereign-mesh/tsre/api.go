package tsre

import (
	"fmt"
	"net/http"

	pb "github.com/pqr-info/sovereign-mesh/proto"
)

// TSREAPI exposes the MIDI-to-Tensor Translation interface for the React frontend
type TSREAPI struct {
	Engine *TemporalEngine
}

func NewTSREAPI(engine *TemporalEngine) *TSREAPI {
	return &TSREAPI{
		Engine: engine,
	}
}

// StartHTTP starts a simple HTTP interface for intercepting UI MIDI buffers
func (api *TSREAPI) StartHTTP(port int) {
	mux := http.NewServeMux()
	
	mux.HandleFunc("/api/v1/tsre/mutate", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "POST required", http.StatusMethodNotAllowed)
			return
		}

		// In a production environment, this parses incoming JSON/Protobuf into a QStateMutation
		// For this implementation, we simulate an incoming MIDI CC event from the React Mixer
		mockMutation := &pb.QStateMutation{
			SequenceId:          "0x1B22-2C4A",
			CoordinateVector:    []float32{1.0, 0.5, -2.1, 4.4, 0.99}, // 5D Tensor coordinate
			MidiCcValue:         127,
			HitlReplayTimestamp: 0, // No HITL override by default
			MutationType:        "MIDI_SYNC_EQ",
		}

		api.Engine.IngestMutation(mockMutation)

		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("{\"status\": \"MUTATION_ACCEPTED\", \"timeline_stability\": 0.98}"))
	})

	fmt.Printf("[TSRE-API] MIDI-to-Tensor Translation API listening on :%d\n", port)
	go http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
}
