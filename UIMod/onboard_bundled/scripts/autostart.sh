#!/usr/bin/env bash
set -e

# This script serves two purposes:
# 1. Installation: Creates and configures a systemd service (ssui.service) to run the StationeersServerControl (StationeersServerUI) application.
# 2. Runtime: When executed and Service already installed, finds and runs the latest StationeersServerControl binary (matching StationeersServerControlv*).
# The systemd service uses ExecStart=$SCRIPT_PATH to run this script, which then dynamically selects the latest binary version of SSUI to run.
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
SSUI_BINARY=$(find "$BASEDIR" -maxdepth 1 -name 'StationeersServerControlv*' -type f -executable -print0 | xargs -0 ls -t | head -n 1)
if [[ -z "$SSUI_BINARY" || ! -x "$SSUI_BINARY" ]]; then
  echo "Error: Could not find executable StationeersServerControl binary in $BASEDIR."
  exit 1
fi

# Create a link to the latest binary in /usr/local/bin for easier access
sudo ln -sf "$SSUI_BINARY" /usr/local/bin/StationeersServerControl
if ! sudo chmod 0755 /usr/local/bin/StationeersServerControl; then
  echo "Error: Failed to create symlink or set executable permission."
  exit 1
fi

# If the service is already running, stop it before updating
if systemctl is-active --quiet ssui.service; then
  echo "Stopping ssui.service for the update..."
  if ! sudo systemctl stop ssui.service; then
    echo "Error: Failed to stop ssui.service."
    exit 1
  fi
fi

# Create the systemd service file pointing to this script
if ! sudo tee /etc/systemd/system/ssui.service > /dev/null <<EOF
[Unit]
Description=Stationeers Server UI
After=network.target

[Service]
Type=simple
Restart=always
RestartSec=5s
User=$(whoami)
WorkingDirectory=$BASEDIR
ExecStart=/usr/local/bin/StationeersServerControl

[Install]
WantedBy=multi-user.target
EOF
then
  echo "Error: Failed to create service file."
  exit 1
fi

# Set service file permissions
sudo chmod 0644 /etc/systemd/system/ssui.service
if ! sudo systemctl daemon-reload; then
  echo "Error: Failed to reload systemd daemon."
  exit 1
fi
if ! sudo systemctl enable ssui.service; then
  echo "Error: Failed to enable ssui.service."
  exit 1
fi
if ! sudo systemctl start ssui.service; then
  echo "Error: Failed to start ssui.service."
  exit 1
fi

echo "Success! Service installed in '/etc/systemd/system/ssui.service'"