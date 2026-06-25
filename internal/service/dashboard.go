package service

import (
	"context"
	_ "embed"
	"encoding/json"
	"net/http"
	"sort"
	"time"

	"github.com/thealanphipps-del/pqr/internal/infrastructure/db"
)

type DashboardService struct {
	db     *db.CockroachRepository
	goback *GobackService
}

func NewDashboardService(repo *db.CockroachRepository, goback *GobackService) *DashboardService {
	return &DashboardService{
		db:     repo,
		goback: goback,
	}
}

func (d *DashboardService) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/health", d.handleHealth)
	mux.HandleFunc("/fixes", d.handleFixes)
	mux.HandleFunc("/snapshots", d.handleSnapshots)
	mux.HandleFunc("/journal", d.handleJournal)
	mux.HandleFunc("/timeline", d.handleTimeline)
	mux.HandleFunc("/diff", d.handleDiff)
	mux.HandleFunc("/revert", d.handleRevert)
	mux.HandleFunc("/knowledge", d.handleKnowledge)
	mux.HandleFunc("/knowledge/search", d.handleKnowledgeSearch)
	mux.HandleFunc("/", d.handleIndex)
}

func (d *DashboardService) Start() {
	mux := http.NewServeMux()
	d.RegisterRoutes(mux)
	go http.ListenAndServe("127.0.0.1:7777", mux)
}

//go:embed index.html
var indexHTML []byte

func (d *DashboardService) handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write(indexHTML)
}

func (d *DashboardService) handleHealth(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	
	// Assess Goback health
	snapshots, _ := d.db.GetSnapshots(ctx)
	journal, _ := d.db.GetJournal(ctx)
	
	var lastSnapshot time.Time
	var genesisPresent bool
	if len(snapshots) > 0 {
		lastSnapshot = snapshots[0].Timestamp
		for _, s := range snapshots {
			if s.IsGenesis {
				genesisPresent = true
				break
			}
		}
	}
	
	var lastJournal time.Time
	if len(journal) > 0 {
		lastJournal = journal[0].Timestamp
	}

	health := map[string]interface{}{
		"status": "healthy",
		"goback": map[string]interface{}{
			"last_snapshot":    lastSnapshot,
			"last_journal":     lastJournal,
			"genesis_present":  genesisPresent,
			"snapshot_count":   len(snapshots),
			"journal_count":    len(journal),
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(health)
}

func (d *DashboardService) handleFixes(w http.ResponseWriter, r *http.Request) {
	fixes, err := d.db.ListTopFixes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fixes)
}

func (d *DashboardService) handleSnapshots(w http.ResponseWriter, r *http.Request) {
	snapshots, err := d.db.GetSnapshots(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(snapshots)
}

func (d *DashboardService) handleJournal(w http.ResponseWriter, r *http.Request) {
	journal, err := d.db.GetJournal(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(journal)
}

type TimelineEvent struct {
	Timestamp time.Time              `json:"timestamp"`
	Type      string                 `json:"type"` // "snapshot", "journal", "error_history"
	Data      map[string]interface{} `json:"data"`
}

func (d *DashboardService) handleTimeline(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var events []TimelineEvent

	// 1. system_snapshots
	snapshots, _ := d.db.GetSnapshots(ctx)
	for _, s := range snapshots {
		var fp map[string]interface{}
		if len(s.HostFingerprint) > 0 {
			json.Unmarshal(s.HostFingerprint, &fp)
		}
		events = append(events, TimelineEvent{
			Timestamp: s.Timestamp,
			Type:      "snapshot",
			Data: map[string]interface{}{
				"is_genesis":       s.IsGenesis,
				"proto_checksum":   s.ProtoChecksum,
				"host_fingerprint": fp,
			},
		})
	}

	// 2. action_journal
	journal, _ := d.db.GetJournal(ctx)
	for _, j := range journal {
		events = append(events, TimelineEvent{
			Timestamp: j.Timestamp,
			Type:      "journal",
			Data: map[string]interface{}{
				"id":     j.ID,
				"action": j.Action,
			},
		})
	}

	// 3. error_solution_history
	history, _ := d.db.GetErrorHistory(ctx)
	for _, h := range history {
		events = append(events, TimelineEvent{
			Timestamp: h.SnapshotTime,
			Type:      "error_history",
			Data: map[string]interface{}{
				"signature":    h.SignatureHash,
				"success_rate": h.SuccessRate,
			},
		})
	}

	// Sort chronologically (oldest to newest)
	sort.Slice(events, func(i, j int) bool {
		return events[i].Timestamp.Before(events[j].Timestamp)
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

func (d *DashboardService) handleDiff(w http.ResponseWriter, r *http.Request) {
	tsStr := r.URL.Query().Get("ts")
	if tsStr == "" {
		http.Error(w, "missing ts parameter", http.StatusBadRequest)
		return
	}
	ts, err := time.Parse(time.RFC3339Nano, tsStr)
	if err != nil {
		ts, err = time.Parse(time.RFC3339, tsStr)
		if err != nil {
			http.Error(w, "invalid ts parameter", http.StatusBadRequest)
			return
		}
	}

	snap, err := d.db.GetSnapshot(context.Background(), ts)
	if err != nil {
		http.Error(w, "snapshot not found", http.StatusNotFound)
		return
	}

	// Also get current live fingerprint
	currentFp := collectHostFingerprint()

	var fp map[string]interface{}
	if len(snap.HostFingerprint) > 0 {
		json.Unmarshal(snap.HostFingerprint, &fp)
	}

	var engine map[string]interface{}
	json.Unmarshal(snap.EngineState, &engine)

	var config map[string]interface{}
	json.Unmarshal(snap.ConfigState, &config)

	payload := map[string]interface{}{
		"timestamp": snap.Timestamp,
		"engine_state": engine,
		"config_state": config,
		"snapshot_fingerprint": fp,
		"current_fingerprint": currentFp,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payload)
}

func (d *DashboardService) handleRevert(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Timestamp string `json:"timestamp"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := d.goback.System(req.Timestamp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"reverted"}`))
}

func (d *DashboardService) handleKnowledge(w http.ResponseWriter, r *http.Request) {
	entries, err := d.db.GetRecentKnowledge(context.Background(), 50)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entries)
}

func (d *DashboardService) handleKnowledgeSearch(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if q == "" {
		http.Error(w, "missing q parameter", http.StatusBadRequest)
		return
	}
	entries, err := d.db.SearchKnowledge(context.Background(), q)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entries)
}
