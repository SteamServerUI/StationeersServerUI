#!/usr/bin/env bash

# This script serves two purposes:
# 1. Installation: Creates and configures a systemd service (ssui.service) to run the StationeersServerControl (StationeersServerUI) application.
# 2. Runtime: When executed and Service already installed, finds and runs the latest StationeersServerControl binary (matching StationeersServerControlv*).
# The systemd service uses ExecStart=$SCRIPT_PATH to run this script, which then dynamically selects the latest binary version of SSUI to run.
# Check if running as root to prevent installing a service as root

if [[ $(id -u) = 0 ]]; then
  echo "For security reasons, it is not recommended to run this service as a root user."
  exit 10
fi

# Check if systemd is available
if [[ ! -d /run/systemd/system ]]; then
  echo "Error: systemd is not the active init system."
  exit 2
fi

# Determine the full path of this script
SCRIPT_PATH=$(readlink -f "$0")

# Determine the base directory
BASEDIR=$(dirname "$SCRIPT_PATH")
if [[ -z "$BASEDIR" || ! -d "$BASEDIR" ]]; then
  echo "Error: Could not determine base directory from SCRIPT_PATH: '$SCRIPT_PATH'."
  exit 3
fi


USERNAME=$(whoami)

# Create the systemd service file pointing to this script
sudo tee /etc/systemd/system/stationeersserverui.service > /dev/null <<EOF
[Unit]
Description=Stationeers Server UI
After=network.target

[Service]
Type=simple
Restart=always
RestartSec=5s
User=$USERNAME
WorkingDirectory=$BASEDIR
ExecStart=$BASEDIR/StationeersServerUI

[Install]
WantedBy=multi-user.target
EOF

sudo chmod 0644 /etc/systemd/system/ssui.service
if [[ $? -ne 0 ]]; then
  echo "Error: Failed to set permissions on /etc/systemd/system/ssui.service."
  exit 6
fi

# Reload systemd daemon
sudo systemctl daemon-reload
if [[ $? -ne 0 ]]; then
  echo "Error: Failed to reload systemd daemon."
  exit 7
fi

# Enable the service
sudo systemctl enable ssui.service
if [[ $? -ne 0 ]]; then
  echo "Error: Failed to enable ssui.service."
  exit 8
fi

# Start the service
sudo systemctl start ssui.service
if [[ $? -ne 0 ]]; then
  echo "Error: Failed to start ssui.service."
  exit 9
fi

echo "Success! Service installed in '/etc/systemd/system/ssui.service'"