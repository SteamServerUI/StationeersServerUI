#!/bin/bash
set -e

# Check if running as root to prevent installing a service as root
if [[ $(id -u) = 0 ]]; then
  echo "For security reasons, it is not recommended to run this service as a root user."
  exit 1
fi

# Check if systemd is available
if ! command -v systemctl &> /dev/null; then
  echo "Error: systemd is not available on this system."
  exit 1
fi

# Determine the full path of this script
SCRIPT_PATH=$(readlink -f "$0")

# Determine the base directory and locate the StationeersServerControl binary
BASEDIR=$(dirname "$SCRIPT_PATH")
if [[ -z "$BASEDIR" || ! -d "$BASEDIR" ]]; then
  echo "Error: Could not determine base directory."
  exit 1
fi

# Find the last modified SSUI binary if multiple exist
SSUI_BINARY=$(ls -t "$BASEDIR"/StationeersServerControlv* 2>/dev/null | head -n 1)
if [[ -z "$SSUI_BINARY" || ! -x "$SSUI_BINARY" ]]; then
  echo "Error: Could not find executable StationeersServerControl binary in $BASEDIR."
  exit 1
fi

# If the systemd service file already exists, just exec the SSUI binary
if [[ -f /etc/systemd/system/ssui.service ]]; then
  echo "Service already installed. Starting SSUI..."
  exec "$SSUI_BINARY"
fi

# Create the systemd service file pointing to this script
sudo tee /etc/systemd/system/ssui.service > /dev/null <<EOF
[Unit]
Description=Stationeers Server UI
After=network.target

[Service]
Type=simple
Restart=always
RestartSec=5s
User=$(whoami)
WorkingDirectory=$BASEDIR
ExecStart=$SCRIPT_PATH

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
sudo systemctl enable ssui.service
sudo systemctl start ssui.service

echo "Success! Service installed in '/etc/systemd/system/ssui.service'"