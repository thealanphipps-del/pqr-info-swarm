# 🌐 Multi-Node Sovereign Mesh

To achieve true failover and redundancy across physical machines (**Alienware** and **201.mh**), we utilize a distributed architecture where nodes collaborate over the network.

## 🏗️ Architecture Options

### Option A: Shared Genesis (Easiest)
*   **Alienware (Genesis Node)**: Runs DB, Vault, Nginx, and PQR Server.
*   **201.mh (Worker Node)**: Runs PQR Server only, connecting back to Alienware's DB.
*   *Pros*: Simple setup.
*   *Cons*: Alienware is a single point of failure for the DB.

### Option B: Distributed Fabric (Recommended)
*   **Both Nodes**: Run CockroachDB in a cluster and PQR Servers.
*   **Cloudflare Tunnel**: Run on both machines with the same token for global load balancing.
*   *Pros*: High availability. If one machine goes down, the other takes over seamlessly.

---

## 🚀 Setup Instructions (Distributed Fabric)

### 1. Preparation
Ensure both machines can ping each other. Let's assume:
*   **Alienware**: `192.168.1.100`
*   **201.mh**: `192.168.1.201`

### 2. Configure Alienware (Node 1)
Update your `.env` or set environment variables:
```powershell
$env:NODE_ID = "alienware"
$env:JOIN_IPS = "192.168.1.100,192.168.1.201"
.\start_pqr.ps1
```

### 3. Configure 201.mh (Node 2)
Clone the repository to the second machine and run:
```powershell
$env:NODE_ID = "201.mh"
$env:JOIN_IPS = "192.168.1.100,192.168.1.201"
.\start_pqr.ps1
```

## 🧬 Updated Docker Compose for Multi-Node
I have updated the `docker-compose.yml` to use a `${NODE_IP}` variable for the database listener, allowing nodes to find each other.

## ☁️ Cloudflare Redundancy
By running the `tunnel` service on both machines with the **same token**, Cloudflare automatically creates a high-availability "anycast" setup. Traffic to `pqr.info` will be routed to whichever machine is closest and healthiest.

---

## 📋 Verification
1.  **DB Cluster**: Visit `http://localhost:8081` on either machine. You should see **2 Nodes** in the CockroachDB console.
2.  **API**: Test `https://pqr.info/REST/2.0/health`. You should get a response even if you stop one machine's Docker stack.
