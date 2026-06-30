#!/bin/bash
# Sovereign Swarm Mesh Startup Initialization
# Source this file in your ~/.bashrc

MESH_ROOT="/home/aellok/sovereign-mesh"
VENV_PATH="$MESH_ROOT/.venv"

# 1. Aesthetics
BOLD="\033[1m"
CYAN="\033[96m"
GREEN="\033[92m"
GOLD="\033[93m"
RESET="\033[0m"

echo -e "${BOLD}${CYAN}=== SOVEREIGN SWARM: MESH CONTROL PLANE ===${RESET}"

# 2. Activate Environment
if [ -d "$VENV_PATH" ]; then
    source "$VENV_PATH/bin/activate"
    echo -e "${GREEN}✅ Mesh environment activated.${RESET}"
else
    echo -e "${GOLD}⚠️  Mesh virtual environment not found. Please run 'cd $MESH_ROOT && python3 -m venv .venv'.${RESET}"
fi

# 3. Path Exports
export PYTHONPATH="$MESH_ROOT:$PYTHONPATH"
export PATH="$MESH_ROOT/cmd:$PATH"

# 4. Status Check
if [ -f "$MESH_ROOT/mudd_interface.py" ]; then
    echo -e "${CYAN}Mesh services ready. Type 'gmudd' to ignite the Swarm Matrix.${RESET}"
else
    echo -e "${GOLD}Mesh components not detected in $MESH_ROOT${RESET}"
fi
