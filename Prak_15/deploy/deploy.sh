#!/usr/bin/env bash
# deploy.sh — copies a new binary to VPS and restarts the service.
# Usage: ./deploy.sh <VPS_USER> <VPS_IP>

set -euo pipefail

VPS_USER="${1:-user}"
VPS_IP="${2:?VPS IP required}"

echo "==> Building Linux binary..."
GOOS=linux GOARCH=amd64 go build -o bin/tasks ./cmd/server

echo "==> Copying binary to VPS..."
scp bin/tasks "${VPS_USER}@${VPS_IP}:/tmp/tasks"

echo "==> Deploying on VPS..."
ssh "${VPS_USER}@${VPS_IP}" bash << 'REMOTE'
  set -e
  sudo systemctl stop tasks
  sudo mv /opt/tasks/tasks /opt/tasks/tasks.old 2>/dev/null || true
  sudo mv /tmp/tasks /opt/tasks/tasks
  sudo chown tasksuser:tasksuser /opt/tasks/tasks
  sudo chmod 755 /opt/tasks/tasks
  sudo systemctl start tasks
  sudo systemctl status tasks --no-pager
REMOTE

echo "==> Done."
