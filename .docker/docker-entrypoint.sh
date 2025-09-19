#!/bin/sh
set -e

UMASK_VAL="022"

APP_USER="app"
APP_GROUP="app"
APP_HOME="/app"

detect_target_ids() {
    # Try to detect from common bind-mounted paths
    TARGET_UID=""
    TARGET_GID=""
    for p in \
        "$APP_HOME/UIMod" \
        "$APP_HOME/saves" \
        "$APP_HOME"; do
        if [ -e "$p" ]; then
            set -- $(stat -c '%u %g' "$p" 2>/dev/null || echo "0 0")
            u="$1"; g="$2"
            # Prefer non-root IDs when available
            if [ "$u" != "0" ]; then TARGET_UID="$u"; fi
            if [ "$g" != "0" ]; then TARGET_GID="$g"; fi
            # If both found, stop early
            if [ -n "$TARGET_UID" ] && [ -n "$TARGET_GID" ]; then
                break
            fi
        fi
    done

    # Fallback if detection failed
    if [ -z "$TARGET_UID" ]; then TARGET_UID="1000"; fi
    if [ -z "$TARGET_GID" ]; then TARGET_GID="1000"; fi
}

is_path_on_mount() {
    # True if the given path is a mountpoint or inside a mountpoint
    # Compare against the mount points (field 5) in /proc/self/mountinfo
    p="$1"
    while read -r _ _ _ _ mp _; do
        # Exact mountpoint
        if [ "$p" = "$mp" ]; then
            return 0
        fi
        # Under the mountpoint (mp is a prefix followed by '/')
        case "$p" in
            "$mp"/*) return 0 ;;
        esac
    done < /proc/self/mountinfo
    return 1
}

ensure_group() {
    if getent group "$APP_GROUP" >/dev/null 2>&1; then
        # Group exists; update GID if necessary
        CURRENT_GID=$(getent group "$APP_GROUP" | cut -d: -f3)
        if [ "$CURRENT_GID" != "$TARGET_GID" ]; then
            groupmod -o -g "$TARGET_GID" "$APP_GROUP" 2>/dev/null || true
        fi
    else
        # Try to reuse existing group with same GID, otherwise create
        if getent group "$TARGET_GID" >/dev/null 2>&1; then
            APP_GROUP=$(getent group "$TARGET_GID" | cut -d: -f1)
        else
            groupadd -o -g "$TARGET_GID" "$APP_GROUP" || true
        fi
    fi
}

ensure_user() {
    if id -u "$APP_USER" >/dev/null 2>&1; then
        CURRENT_UID=$(id -u "$APP_USER")
        if [ "$CURRENT_UID" != "$TARGET_UID" ]; then
            usermod -o -u "$TARGET_UID" "$APP_USER" 2>/dev/null || true
        fi
    # Ensure user's primary group is APP_GROUP (with target GID)
        usermod -g "$APP_GROUP" "$APP_USER" 2>/dev/null || true
        usermod -d "$APP_HOME" -s /bin/sh "$APP_USER" 2>/dev/null || true
    else
        useradd -m -d "$APP_HOME" -s /bin/sh -u "$TARGET_UID" -g "$APP_GROUP" "$APP_USER" 2>/dev/null || true
    fi
}

fix_perms() {
    # Create expected directories if missing (bind mounts may provide them)
    mkdir -p \
        "$APP_HOME" \
        "$APP_HOME/saves" \
        "$APP_HOME/UIMod" \
        "$APP_HOME/UIMod/config" \
        "$APP_HOME/UIMod/tls"

    # Only adjust ownership for paths that are not on/under a bind mount
    for p in \
        "$APP_HOME" \
        "$APP_HOME/saves" \
        "$APP_HOME/UIMod"; do
        if ! is_path_on_mount "$p"; then
            chown "$APP_USER:$APP_GROUP" "$p" 2>/dev/null || true
        fi
    done
}

main() {
    # Only attempt privilege drop if running as root
    if [ "$(id -u)" = "0" ]; then
        detect_target_ids
        ensure_group
        ensure_user
        fix_perms
        umask "$UMASK_VAL"

        # Execute as the app user
        if [ $# -eq 0 ]; then
            exec gosu "$APP_USER:$APP_GROUP" StationeersServerControl
        else
            exec gosu "$APP_USER:$APP_GROUP" "$@"
        fi
    else
        # Already running as non-root
        umask "$UMASK_VAL"
        if [ $# -eq 0 ]; then
            exec StationeersServerControl
        else
            exec "$@"
        fi
    fi
}

main "$@"
