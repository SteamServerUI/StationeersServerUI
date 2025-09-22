#!/usr/bin/env sh

cp /opt/SSUIBuildFiles/StationeersServerUI /app/StationeersServerUI
cp /opt/SSUIBuildFiles/LICENSE /app/LICENSE
chmod +x /app/StationeersServerUI
exec /app/StationeersServerUI "$@"