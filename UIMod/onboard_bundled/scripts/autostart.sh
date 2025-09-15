#!/bin/bash

if [[ $(id -u) = 0 ]]; then
  echo "For security reasons, it is not recommended to run this service as a root user."
  exit 1
fi

BASEDIR=$(dirname $(readlink -f "$0"))
BINARY=$(find $BASEDIR/* -maxdepth 0 -name "StationeersServerControlv*" -print -quit)

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

echo "Success! Service installed in \`/etc/systemd/system/ssui.service\`"