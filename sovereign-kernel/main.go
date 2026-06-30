package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/big"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// -----------------------------------------------------------------------------
// Core types & 5D Mapping
// -----------------------------------------------------------------------------

type Address5D struct {
	X, Y, Z, Phi, Lambda int `json:"x"`
}

type TesseractVertex struct {
	V0, V1, V2, V3 int `json:"v0"`
}

func IPv6To5D(ip net.IP) Address5D {
	ipv6 := ip.To16()
	if ipv6 == nil {
		return Address5D{}
	}

	val := new(big.Int).SetBytes(ipv6)

	v0 := new(big.Int).Rsh(val, 103).Uint64() & 0x1FFFFFF
	v1 := new(big.Int).Rsh(val, 78).Uint64() & 0x1FFFFFF
	v2 := new(big.Int).Rsh(val, 53).Uint64() & 0x1FFFFFF
	v3 := new(big.Int).Rsh(val, 28).Uint64() & 0x1FFFFFF
	v4 := new(big.Int).Uint64() & 0xFFFFFFF

	return Address5D{
		X:      (int(v0) % 150) - 75, 
		Y:      (int(v1) % 150) - 75,
		Z:      (int(v2) % 60) - 30,
		Phi:    (int(v3) % 40) - 20,
		Lambda: (int(v4) % 20) - 10,
	}
}

func Hash5D(a Address5D) string {
	return fmt.Sprintf("%d-%d-%d-%d-%d", a.X, a.Y, a.Z, a.Phi, a.Lambda)
}

func Distance5D(a, b Address5D) float64 {
	dx := float64(a.X - b.X)
	dy := float64(a.Y - b.Y)
	dz := float64(a.Z - b.Z)
	dp := float64(a.Phi - b.Phi)
	dl := float64(a.Lambda - b.Lambda)
	return math.Sqrt(dx*dx + dy*dy + dz*dz + dp*dp + dl*dl)
}

type CoordinateMapper struct {
	mu           sync.RWMutex
	IPv6Registry map[string]net.IP
	destTo5D     map[string]Address5D
	fiveDToDest  map[Address5D]string
}

func NewCoordinateMapper() *CoordinateMapper {
	return &CoordinateMapper{
		IPv6Registry: make(map[string]net.IP),
		destTo5D:     make(map[string]Address5D),
		fiveDToDest:  make(map[Address5D]string),
	}
}

func (m *CoordinateMapper) MapDestination(dest string, destType string) Address5D {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if addr, ok := m.destTo5D[dest]; ok {
		return addr
	}

	var addr Address5D
	if destType == "ipv6" {
		ip, _, err := net.ParseCIDR(dest)
		if err != nil {
			ip = net.ParseIP(dest)
		}
		if ip != nil {
			addr = IPv6To5D(ip)
			m.IPv6Registry[Hash5D(addr)] = ip
		}
	} else {
		hash := sha256.Sum256([]byte(dest))
		addr = Address5D{
			X:      int(binary.BigEndian.Uint32(hash[0:4])) % 5000,
			Y:      int(binary.BigEndian.Uint32(hash[4:8])) % 5000,
			Z:      int(binary.BigEndian.Uint32(hash[8:12])) % 5000,
			Phi:    int(binary.BigEndian.Uint32(hash[12:16])) % 5000,
			Lambda: int(binary.BigEndian.Uint32(hash[16:20])) % 5000,
		}
	}

	m.destTo5D[dest] = addr
	m.fiveDToDest[addr] = dest
	return addr
}

type AgentState struct {
	ID        string          `json:"id"`
	Address5D Address5D       `json:"address5d"`
	Vertex    TesseractVertex `json:"vertex"`
	Epoch     int64           `json:"epoch"`
	Neighbors []string        `json:"neighbors"`
}

type ExternalEndpoint struct {
	ASN       int             `json:"asn,omitempty"`
	Type      string          `json:"type"`
	Address   string          `json:"address"`
	Address5D Address5D       `json:"address5d"`
	Vertex    TesseractVertex `json:"vertex"`
	Epoch     int64           `json:"epoch"`
	Hostname  string          `json:"hostname,omitempty"`
	Status    string          `json:"status"`
}

type CircuitHop struct {
	AgentID    string          `json:"agentId,omitempty"`
	Onion      string          `json:"onion,omitempty"`
	Hostname   string          `json:"hostname,omitempty"`
	IPv6       string          `json:"ipv6,omitempty"`
	ASN        int             `json:"asn,omitempty"`
	Coordinate Address5D       `json:"coordinate"`
	Vertex     TesseractVertex `json:"vertex"`
}

type EventCircuit struct {
	ID    string       `json:"id"`
	Epoch int64        `json:"epoch"`
	Hops  []CircuitHop `json:"hops"`
}

type HealthMetrics struct {
	Epoch              int64   `json:"epoch"`
	CircuitDensity     float64 `json:"circuitDensity"`
	ContinuityVelocity float64 `json:"continuityVelocity"`
	AgentChurn         float64 `json:"agentChurn"`
}

type AnnouncePayload struct {
	Type       string    `json:"type"`
	ASN        int       `json:"asn"`
	IPv6Prefix string    `json:"ipv6Prefix"`
	Coordinate Address5D `json:"coordinate"`
	Hostname   string    `json:"hostname,omitempty"`
	Onion      string    `json:"onion,omitempty"`
	Timestamp  int64     `json:"timestamp"`
	Status     string    `json:"status,omitempty"`
}

// -----------------------------------------------------------------------------
// Event Model
// -----------------------------------------------------------------------------

type EventType string

const (
	EventContinuity  EventType = "continuity"
	EventLineage     EventType = "lineage"
	EventTesseract   EventType = "tesseract"
	EventAgentJoin   EventType = "agent_join"
	EventEpoch       EventType = "epoch"
	EventCircuitType EventType = "circuit"
	EventHealth      EventType = "health"
	Event5DAnnounce  EventType = "5d_announce"
	Event5DWithdraw  EventType = "5d_withdraw"
)

type KernelEvent struct {
	Type      EventType `json:"type"`
	Timestamp int64     `json:"timestamp"`
	Payload   any       `json:"payload"`
}

type LineageRecord struct {
	ClusterID string                   `json:"clusterId"`
	Epoch     int64                    `json:"epoch"`
	Event     string                   `json:"event"`
	Agents    map[string]AgentSnapshot `json:"agents"`
}

type AgentSnapshot struct {
	ID        string          `json:"id"`
	Address5D Address5D       `json:"address5d"`
	Vertex    TesseractVertex `json:"vertex"`
	Epoch     int64           `json:"epoch"`
}

type SovereignContinuityEvent struct {
	ClusterID  string `json:"clusterId"`
	Epoch      int64  `json:"epoch"`
	AgentCount int    `json:"agentCount"`
}

// -----------------------------------------------------------------------------
// WebSocket Hub & Bridge
// -----------------------------------------------------------------------------

type Client struct {
	Conn *websocket.Conn
	Send chan KernelEvent
}

type Hub struct {
	Clients    map[*Client]struct{}
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan KernelEvent
}

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.Register:
			h.Clients[c] = struct{}{}
		case c := <-h.Unregister:
			delete(h.Clients, c)
			close(c.Send)
		case evt := <-h.Broadcast:
			for c := range h.Clients {
				select {
				case c.Send <- evt:
				default:
					delete(h.Clients, c)
					close(c.Send)
				}
			}
		}
	}
}

func wireKernelToHub(kernel *ContinuityKernel, hub *Hub) {
	go func() {
		for evt := range kernel.Events {
			hub.Broadcast <- evt
		}
	}()
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func serveWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[WS] upgrade error: %v", err)
		return
	}

	client := &Client{
		Conn: conn,
		Send: make(chan KernelEvent, 256),
	}
	hub.Register <- client

	go writePump(client)
	go readPump(client, hub)
}

func writePump(c *Client) {
	for evt := range c.Send {
		if err := c.Conn.WriteJSON(evt); err != nil {
			break
		}
	}
	c.Conn.Close()
}

func readPump(c *Client, hub *Hub) {
	defer func() {
		hub.Unregister <- c
		c.Conn.Close()
	}()
	for {
		if _, _, err := c.Conn.ReadMessage(); err != nil {
			break
		}
	}
}

func startVisualizerBridge(kernel *ContinuityKernel) {
	hub := &Hub{
		Clients:    make(map[*Client]struct{}),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan KernelEvent, 1024),
	}
	go hub.Run()
	wireKernelToHub(kernel, hub)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWS(hub, w, r)
	})

	go func() {
		log.Printf("[HTTP] Visualizer WebSocket bridge listening on 127.0.0.3:1112")
		if err := http.ListenAndServe("127.0.0.3:1112", nil); err != nil {
			log.Printf("[HTTP] bridge error: %v", err)
		}
	}()
}

// -----------------------------------------------------------------------------
// Stubbed clients
// -----------------------------------------------------------------------------

type SubstrateClient interface {
	GetLastEpoch(clusterID string) int64
	StoreVertexMatrix(clusterID string, epoch int64, vertices map[string]TesseractVertex) error
	AppendLineage(record LineageRecord) error
}

type MeshQuorumClient interface {
	BroadcastEpoch(clusterID string, epoch int64) error
	EmitContinuity(event SovereignContinuityEvent) error
}

type SRRPClient interface {
	UpdateRoutingTables(agents map[string]*AgentState) error
	RebuildCircuits(epoch int64, agents map[string]*AgentState) error
}

type TesseractTracer interface {
	Assign5D(agentID string) Address5D
	InitialVertex(addr Address5D) TesseractVertex
	Evolve(addr Address5D, current TesseractVertex, epoch int64) TesseractVertex
}

// -----------------------------------------------------------------------------
// Continuity kernel
// -----------------------------------------------------------------------------

type ContinuityKernel struct {
	mu             sync.Mutex
	ClusterID      string
	Epoch          int64
	Agents         map[string]*AgentState
	Externals      map[string]*ExternalEndpoint
	Mapper         *CoordinateMapper
	Substrate      SubstrateClient
	MeshQuorum     MeshQuorumClient
	SRRP           SRRPClient
	Tesseract      TesseractTracer
	DeltaCycle     time.Duration
	Quit           chan struct{}
	Events         chan KernelEvent
	
	circuitsFormed int
}

func (k *ContinuityKernel) Start() {
	k.loadPersistedEpoch()
	k.loadPersistedAgents()
	if err := k.announceEpoch(); err != nil {
		log.Printf("[Kernel] failed to announce epoch: %v", err)
	}
	go k.runDeltaLoop()
	go k.runAgentAcceptLoop()
}

func (k *ContinuityKernel) loadPersistedEpoch() {
	last := k.Substrate.GetLastEpoch(k.ClusterID)
	if last == 0 {
		last = 1
	}
	k.Epoch = last
	log.Printf("[Kernel] loaded epoch %d", k.Epoch)
}

func (k *ContinuityKernel) loadPersistedAgents() {
	k.Agents = make(map[string]*AgentState)
	k.Externals = make(map[string]*ExternalEndpoint)
	log.Printf("[Kernel] initialized empty agent set")
}

func (k *ContinuityKernel) announceEpoch() error {
	err := k.MeshQuorum.BroadcastEpoch(k.ClusterID, k.Epoch)
	k.Events <- KernelEvent{
		Type:      EventEpoch,
		Timestamp: time.Now().UnixMilli(),
		Payload: map[string]any{
			"epoch":     k.Epoch,
			"clusterId": k.ClusterID,
		},
	}
	return err
}

func (k *ContinuityKernel) evolveTesseract() {
	k.mu.Lock()
	defer k.mu.Unlock()

	for _, agent := range k.Agents {
		agent.Vertex = k.Tesseract.Evolve(agent.Address5D, agent.Vertex, k.Epoch)
	}
	
	for _, ext := range k.Externals {
		ext.Vertex = k.Tesseract.Evolve(ext.Address5D, ext.Vertex, k.Epoch)
	}
	
	vertices := k.snapshotVertices()
	if err := k.Substrate.StoreVertexMatrix(k.ClusterID, k.Epoch, vertices); err != nil {
		log.Printf("[Kernel] failed to store vertex matrix: %v", err)
	}

	k.Events <- KernelEvent{
		Type:      EventTesseract,
		Timestamp: time.Now().UnixMilli(),
		Payload: struct {
			Epoch    int64                      `json:"epoch"`
			Vertices map[string]TesseractVertex `json:"vertices"`
		}{
			Epoch:    k.Epoch,
			Vertices: vertices,
		},
	}
}

func (k *ContinuityKernel) storeLineage(event string) {
	k.mu.Lock()
	defer k.mu.Unlock()

	record := LineageRecord{
		ClusterID: k.ClusterID,
		Epoch:     k.Epoch,
		Event:     event,
		Agents:    k.snapshotAgents(),
	}
	if err := k.Substrate.AppendLineage(record); err != nil {
		log.Printf("[Kernel] failed to append lineage: %v", err)
	}

	k.Events <- KernelEvent{
		Type:      EventLineage,
		Timestamp: time.Now().UnixMilli(),
		Payload:   record,
	}
}

func (k *ContinuityKernel) runNeighborDiscovery() {
	k.mu.Lock()
	defer k.mu.Unlock()

	vertices := k.snapshotVertices()
	for id, agent := range k.Agents {
		agent.Neighbors = Compute8NN(id, vertices)
	}
	if err := k.SRRP.UpdateRoutingTables(k.Agents); err != nil {
		log.Printf("[Kernel] failed to update routing tables: %v", err)
	}
}

func (k *ContinuityKernel) maintainSRRP() {
	k.mu.Lock()
	defer k.mu.Unlock()

	if err := k.SRRP.RebuildCircuits(k.Epoch, k.Agents); err != nil {
		log.Printf("[Kernel] failed to rebuild circuits: %v", err)
	}
	
	var hops []CircuitHop
	for id, agent := range k.Agents {
		hops = append(hops, CircuitHop{
			AgentID:    id,
			Coordinate: agent.Address5D,
			Vertex:     agent.Vertex,
		})
		if len(hops) == 2 { break }
	}
	
	if len(hops) > 1 {
		count := 0
		for _, ext := range k.Externals {
			hop := CircuitHop{
				Coordinate: ext.Address5D,
				Vertex:     ext.Vertex,
				ASN:        ext.ASN,
			}
			if ext.Type == "ipv6" {
				hop.IPv6 = ext.Address
			} else if ext.Type == "onion" {
				hop.Onion = ext.Address
			}
			hops = append(hops, hop)
			count++
			if count == 2 { break } // Create AS-to-AS backbone transit
		}

		circuit := EventCircuit{
			ID:    "srrp-circuit-" + time.Now().Format("150405"),
			Epoch: k.Epoch,
			Hops:  hops,
		}
		k.Events <- KernelEvent{
			Type:      EventCircuitType,
			Timestamp: time.Now().UnixMilli(),
			Payload:   circuit,
		}
		k.circuitsFormed++
		log.Printf("[Kernel] emitted 5D-agnostic SRRP circuit with %d hops", len(hops))
	}
}

type AgentJoinInfo struct {
	ID string
}

func (k *ContinuityKernel) runAgentAcceptLoop() {
	for {
		select {
		case <-k.Quit:
			return
		default:
			info := k.waitForAgentJoin()
			k.registerAgent(info)
		}
	}
}

func (k *ContinuityKernel) waitForAgentJoin() AgentJoinInfo {
	time.Sleep(5 * time.Second)
	return AgentJoinInfo{ID: "agent-auto-" + time.Now().Format("150405")}
}

func (k *ContinuityKernel) registerAgent(info AgentJoinInfo) {
	k.mu.Lock()
	defer k.mu.Unlock()
	
	if len(k.Agents) > 20 {
		return
	}

	addr := k.Tesseract.Assign5D(info.ID)
	vertex := k.Tesseract.InitialVertex(addr)
	k.Agents[info.ID] = &AgentState{
		ID:        info.ID,
		Address5D: addr,
		Vertex:    vertex,
		Epoch:     k.Epoch,
	}
	log.Printf("[Kernel] registered agent %s at epoch %d", info.ID, k.Epoch)
	k.storeLineage("agent_join:" + info.ID)

	snapshot := AgentSnapshot{
		ID:        info.ID,
		Address5D: addr,
		Vertex:    vertex,
		Epoch:     k.Epoch,
	}
	k.Events <- KernelEvent{
		Type:      EventAgentJoin,
		Timestamp: time.Now().UnixMilli(),
		Payload:   snapshot,
	}
}

func (k *ContinuityKernel) runDeltaLoop() {
	ticker := time.NewTicker(k.DeltaCycle)
	defer ticker.Stop()

	var prevVertices map[string]TesseractVertex
	var prevAgentsCount int

	for {
		select {
		case <-k.Quit:
			return
		case <-ticker.C:
			k.Epoch++
			log.Printf("[Kernel] δ-cycle tick, advancing to epoch %d", k.Epoch)
			if err := k.announceEpoch(); err != nil {
				log.Printf("[Kernel] epoch announce failed: %v", err)
			}
			
			k.evolveTesseract()
			k.runNeighborDiscovery()
			k.maintainSRRP()
			k.storeLineage("delta_cycle")
			k.emitContinuityEvent()

			k.mu.Lock()
			agentsCount := len(k.Agents)
			circuits := k.circuitsFormed
			k.circuitsFormed = 0
			k.mu.Unlock()

			newAgents := agentsCount - prevAgentsCount
			if newAgents < 0 { newAgents = 0 }
			churn := float64(newAgents) / k.DeltaCycle.Seconds()
			
			density := 0.0
			if agentsCount > 0 {
				density = float64(circuits) / float64(agentsCount)
			}
			
			velocity := 0.0
			currentVertices := k.snapshotVertices()
			if prevVertices != nil {
				var totalDelta float64
				var count int
				for id, v := range currentVertices {
					if pv, ok := prevVertices[id]; ok {
						delta := math.Abs(float64(v.V0-pv.V0)) + math.Abs(float64(v.V1-pv.V1)) + math.Abs(float64(v.V2-pv.V2)) + math.Abs(float64(v.V3-pv.V3))
						totalDelta += delta
						count++
					}
				}
				if count > 0 {
					velocity = totalDelta / float64(count)
				}
			}

			metrics := HealthMetrics{
				Epoch:              k.Epoch,
				CircuitDensity:     density,
				ContinuityVelocity: velocity,
				AgentChurn:         churn,
			}
			k.Events <- KernelEvent{
				Type:      EventHealth,
				Timestamp: time.Now().UnixMilli(),
				Payload:   metrics,
			}

			prevVertices = currentVertices
			prevAgentsCount = agentsCount
		}
	}
}

func (k *ContinuityKernel) emitContinuityEvent() {
	k.mu.Lock()
	defer k.mu.Unlock()

	event := SovereignContinuityEvent{
		ClusterID:  k.ClusterID,
		Epoch:      k.Epoch,
		AgentCount: len(k.Agents),
	}
	if err := k.MeshQuorum.EmitContinuity(event); err != nil {
		log.Printf("[Kernel] failed to emit continuity event: %v", err)
	}

	k.Events <- KernelEvent{
		Type:      EventContinuity,
		Timestamp: time.Now().UnixMilli(),
		Payload:   event,
	}
}

func (k *ContinuityKernel) snapshotVertices() map[string]TesseractVertex {
	k.mu.Lock()
	defer k.mu.Unlock()
	out := make(map[string]TesseractVertex, len(k.Agents) + len(k.Externals))
	for id, agent := range k.Agents {
		out[id] = agent.Vertex
	}
	for id, ext := range k.Externals {
		out[id] = ext.Vertex
	}
	return out
}

func (k *ContinuityKernel) snapshotAgents() map[string]AgentSnapshot {
	k.mu.Lock()
	defer k.mu.Unlock()
	out := make(map[string]AgentSnapshot, len(k.Agents))
	for id, agent := range k.Agents {
		out[id] = AgentSnapshot{
			ID:        agent.ID,
			Address5D: agent.Address5D,
			Vertex:    agent.Vertex,
			Epoch:     agent.Epoch,
		}
	}
	return out
}

func Compute8NN(selfID string, vertices map[string]TesseractVertex) []string {
	neighbors := make([]string, 0, 8)
	for id := range vertices {
		if id == selfID {
			continue
		}
		neighbors = append(neighbors, id)
		if len(neighbors) >= 8 {
			break
		}
	}
	return neighbors
}

type InMemorySubstrate struct{ lastEpoch int64 }
func (s *InMemorySubstrate) GetLastEpoch(c string) int64 { return s.lastEpoch }
func (s *InMemorySubstrate) StoreVertexMatrix(c string, e int64, v map[string]TesseractVertex) error {
	s.lastEpoch = e
	return nil
}
func (s *InMemorySubstrate) AppendLineage(r LineageRecord) error { return nil }

type LogMeshQuorum struct{}
func (m *LogMeshQuorum) BroadcastEpoch(c string, e int64) error { return nil }
func (m *LogMeshQuorum) EmitContinuity(ev SovereignContinuityEvent) error { return nil }

type LogSRRP struct{}
func (s *LogSRRP) UpdateRoutingTables(a map[string]*AgentState) error { return nil }
func (s *LogSRRP) RebuildCircuits(e int64, a map[string]*AgentState) error { return nil }

type SimpleTesseract struct{}
func (t *SimpleTesseract) Assign5D(id string) Address5D {
	return Address5D{X: 1, Y: 2, Z: 3, Phi: 4, Lambda: 5}
}
func (t *SimpleTesseract) InitialVertex(a Address5D) TesseractVertex {
	return TesseractVertex{V0: a.X, V1: a.Y, V2: a.Z, V3: a.Phi}
}
func (t *SimpleTesseract) Evolve(a Address5D, v TesseractVertex, e int64) TesseractVertex {
	return TesseractVertex{
		V0: v.V0 + int(e)%5, 
		V1: v.V1 - int(e)%3, 
		V2: v.V2 + int(e)%2, 
		V3: v.V3 - 1,
	}
}

// -----------------------------------------------------------------------------
// 5D Protocol BGP Bridge
// -----------------------------------------------------------------------------

func startMeshQuorumHTTP(kernel *ContinuityKernel) {
	mux := http.NewServeMux()
	
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	mux.HandleFunc("/5d/announce", func(w http.ResponseWriter, r *http.Request) {
		var req AnnouncePayload
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		
		addrType := "ipv6"
		addr := req.IPv6Prefix
		if req.Onion != "" {
			addrType = "onion"
			addr = req.Onion
		}

		addr5D := kernel.Mapper.MapDestination(addr, addrType)
		
		kernel.mu.Lock()
		endpoint := &ExternalEndpoint{
			ASN:       req.ASN,
			Type:      addrType,
			Address:   addr,
			Address5D: addr5D,
			Vertex:    kernel.Tesseract.InitialVertex(addr5D),
			Epoch:     kernel.Epoch,
			Hostname:  req.Hostname,
			Status:    "alive",
		}
		kernel.Externals[addr] = endpoint
		kernel.mu.Unlock()

		kernel.Events <- KernelEvent{
			Type:      Event5DAnnounce,
			Timestamp: time.Now().UnixMilli(),
			Payload:   req,
		}

		log.Printf("[5D-BGP] Received 5D Announce: AS%d Prefix=%s Coordinate=%+v", req.ASN, addr, addr5D)
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("/5d/withdraw", func(w http.ResponseWriter, r *http.Request) {
		var req AnnouncePayload
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		
		kernel.mu.Lock()
		delete(kernel.Externals, req.IPv6Prefix)
		kernel.mu.Unlock()

		kernel.Events <- KernelEvent{
			Type:      Event5DWithdraw,
			Timestamp: time.Now().UnixMilli(),
			Payload:   req,
		}
		log.Printf("[5D-BGP] Withdrew AS%d Prefix=%s", req.ASN, req.IPv6Prefix)
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("/5d/heartbeat", func(w http.ResponseWriter, r *http.Request) {
		var req AnnouncePayload
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("[5D-BGP] Heartbeat from AS%d", req.ASN)
		w.WriteHeader(http.StatusOK)
	})

	go func() {
		log.Printf("[HTTP] 5D Announcement Daemon listening on 127.0.0.3:1111")
		if err := http.ListenAndServe("127.0.0.3:1111", mux); err != nil {
			log.Printf("[HTTP] daemon error: %v", err)
		}
	}()
}

func main() {
	log.Printf("=== Sovereign Continuity Kernel starting ===")

	kernel := &ContinuityKernel{
		ClusterID:  "sovereign-alpha",
		Mapper:     NewCoordinateMapper(),
		Substrate:  &InMemorySubstrate{},
		MeshQuorum: &LogMeshQuorum{},
		SRRP:       &LogSRRP{},
		Tesseract:  &SimpleTesseract{},
		DeltaCycle: 5 * time.Second,
		Quit:       make(chan struct{}),
		Events:     make(chan KernelEvent, 1024),
	}

	startMeshQuorumHTTP(kernel)
	startVisualizerBridge(kernel)
	
	go func() {
		time.Sleep(5 * time.Second)
		routers := []AnnouncePayload{
			{
				Type:       "5d_announce",
				ASN:        64512,
				IPv6Prefix: "2001:db8:abcd::/48",
				Hostname:   "router-lon.example.net",
			},
			{
				Type:       "5d_announce",
				ASN:        64513,
				IPv6Prefix: "2600:1f18:4a22::/48",
				Hostname:   "router-lax.example.net",
			},
			{
				Type:       "5d_announce",
				ASN:        64514,
				IPv6Prefix: "2a02:26f0:199::/48",
				Hostname:   "router-ams.example.net",
			},
		}

		for _, r := range routers {
			b, _ := json.Marshal(r)
			_, _ = http.Post("http://127.0.0.3:1111/5d/announce", "application/json", bytes.NewReader(b))
		}
	}()

	kernel.Start()

	select {}
}
