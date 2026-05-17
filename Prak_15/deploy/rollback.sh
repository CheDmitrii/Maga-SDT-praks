#!/usr/bin/env bash
# rollback.sh — reverts to the previous binary version on VPS.
# Usage: ./rollback.sh <VPS_USER> <VPS_IP>

set -euo pipefail

VPS_USER="${1:-user}"
VPS_IP="${2:?VPS IP required}"

ssh "${VPS_USER}@${VPS_IP}" bash << 'REMOTE'
  set -e
  if [ ! -f /opt/tasks/tasks.old ]; then
    echo "No previous version found at /opt/tasks/tasks.old"
    exit 1
  fi
  sudo systemctl stop tasks
  sudo mv /opt/tasks/tasks.old /opt/tasks/tasks
  sudo systemctl start tasks
  sudo systemctl status tasks --no-pager
REMOTE

echo "==> Rollback done."
