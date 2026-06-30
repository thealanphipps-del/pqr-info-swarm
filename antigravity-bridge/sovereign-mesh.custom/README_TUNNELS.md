# Sovereign Mesh Tunneling & Ticketing Guide

## Accessing Ticketing System
The RTGO ticketing system is accessible via the internal Mesh portal at:
  - Gateway: 192.168.12.201:8081
  - Bridge API: http://192.168.12.201/api/bridge

## SSH Tunneling Strategy
Tunnels are established using ~/.ssh/config leveraging 39.mh (204.168.138.60) as the jump host.
Nodes 38.mh (62.238.2.240) and 201.mh (89.167.91.81) are reached via:
  ProxyJump 39.mh

## gRPC Bridging
Local gRPC services on ports 1111/11111 are mapped via:
  ssh -L 1111:localhost:1111 -L 11111:localhost:11111 39.mh
