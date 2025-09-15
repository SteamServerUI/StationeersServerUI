#!/bin/bash

if [[ $(id -u) = 0 ]]; then
  echo "For security reasons, it is not recommended to run this service as a root user."
  exit 1
fi

BASEDIR=$(dirname $(readlink -f "$0"))
if [[ -z "$BASEDIR" || ! -d "$BASEDIR" ]]; then
  echo "Error: Could not determine base directory."
  exit 1
fi
BINARY=$(find $BASEDIR/* -maxdepth 0 -name "StationeersServerControlv*" -quit)
if [[ -z "$BINARY" || ! -x "$BINARY" ]]; then
  echo "Error: Could not find executable StationeersServerControl binary in $BASEDIR."
  exit 1
fi

sudo cat <<EOF > /etc/systemd/system/ssui.service
[Unit]
Description=Stationeers Server UI
After=network.target

[Service]
Type=simple
Restart=on-failure
RestartSec=5s
User=$(whoami)
WorkingDirectory=$BASEDIR
ExecStart=$BINARY

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
sudo systemctl enable ssui.service
sudo systemctl start ssui.service

echo "Success! Service installed in '/etc/systemd/system/ssui.service'"