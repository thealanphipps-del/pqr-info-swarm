# pqr.info Firebase App Architecture

## Overview
This Firebase application serves as the public-facing frontend for the Sovereign DAO and the 128-Agent Swarm. It allows users to interact with the emergent personalities, view relationship profiles, and rent Sovereign Dev Environments.

## Features
*   **Swarm Dashboard**: Live visualization of the 64 Human Design and 64 Game Theory agent pairs.
*   **Agent Interaction**: Real-time chat interface connecting users to specific agents (backed by the Sovereign-27 neural core).
*   **DAO Minting Portal**: Interface for new citizens to register and mint their own identities via SURFGO burns.
*   **Sovereign Rentals**: A marketplace to provision and manage temporary GCP/Hetzner dev environments.

## Firebase Integration
*   **Firestore**: Real-time syncing of agent vitality, chat logs, and active dev environment leases.
*   **Authentication**: Passwordless auth linked to the citizen's generated `CitizenPassport` (JWT).
*   **Functions**: Serverless handlers bridging the frontend to the backend gRPC `SovereignCity` service.
