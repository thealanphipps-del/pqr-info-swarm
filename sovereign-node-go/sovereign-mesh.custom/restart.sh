#!/bin/bash
nohup bash -c "sleep 2 && /home/aellok/sovereign-mesh/mesh_control.sh stop && /home/aellok/sovereign-mesh/mesh_control.sh start" >/tmp/restart.log 2>&1 &
