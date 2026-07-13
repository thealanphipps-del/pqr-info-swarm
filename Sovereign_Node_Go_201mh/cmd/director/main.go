package main

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"
	"github.com/miekg/dns"
)

const (
	DNSPort  = "5353"
	LDAPPort = "1389"
	MeshZone = "mesh."
)

var (
	nodeMap = map[string]string{
		"helsinki.mesh":  "10.0.39.1",
		"nuremburg.mesh": "10.0.0.1",
		"s25fe.mesh":     "10.0.39.203",
	}
	mu sync.RWMutex
)

func handleForward(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	mu.RLock()
	if ip, ok := nodeMap[r.Question[0].Name]; ok {
		m.Answer = append(m.Answer, &dns.A{
			Hdr: dns.RR_Header{Name: r.Question[0].Name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
			A:   net.ParseIP(ip),
		})
	}
	mu.RUnlock()
	w.WriteMsg(m)
}

func handleRecursive(w dns.ResponseWriter, r *dns.Msg) {
	c := new(dns.Client)
	in, _, err := c.Exchange(r, "8.8.8.8:53")
	if err == nil {
		w.WriteMsg(in)
	}
}

func startDNS() {
	mux := dns.NewServeMux()
	mux.HandleFunc(MeshZone, handleForward)
	mux.HandleFunc(".", handleRecursive)
	server := &dns.Server{Addr: "127.0.0.1:" + DNSPort, Net: "udp"}
	fmt.Printf("[DIRECTOR] DNS Active on %s\n", DNSPort)
	server.ListenAndServe()
}

func startVirtualLLDP() {
	ticker := time.NewTicker(30 * time.Second)
	for range ticker.C {
		ts := time.Now().Format(time.RFC3339)
		msg := fmt.Sprintf("%s NODE:S25FE IP:10.0.39.203 ROLE:CONTROLLER STATUS:ONLINE\n", ts)
		f, _ := os.OpenFile("/data/data/com.termux/files/home/Sovereign_Node_Go/logs/lldp_mesh.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		f.WriteString(msg)
		f.Close()
	}
}

func startLDAP() {
	fmt.Printf("[DIRECTOR] LDAP Auth Bridge Active on %s\n", LDAPPort)
	l, _ := net.Listen("tcp", "127.0.0.1:"+LDAPPort)
	for {
		conn, _ := l.Accept()
		go func(c net.Conn) {
			fmt.Println("[DIRECTOR] LDAP Identity Query Intercepted -> Tunneling to Hub")
			c.Close()
		}(conn)
	}
}

func main() {
	go startDNS()
	go startLDAP()
	go startVirtualLLDP()
	select {}
}
