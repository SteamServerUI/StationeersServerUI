package diagnostics

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

func platformEnvironment(env *Environment) {
	if data, err := os.ReadFile("/proc/sys/kernel/osrelease"); err == nil {
		env.Kernel = strings.TrimSpace(string(data))
	}
	if data, err := os.ReadFile("/proc/meminfo"); err == nil {
		for _, line := range strings.Split(string(data), "\n") {
			if value, ok := strings.CutPrefix(line, "MemTotal:"); ok {
				kilobytes, _ := strconv.ParseInt(strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(value), "kB")), 10, 64)
				env.MemoryMB = kilobytes / 1024
			}
		}
	}
}

func describeVolume(volume *Volume) {
	path, err := filepath.Abs(volume.Path)
	if err != nil {
		return
	}
	if resolved, err := filepath.EvalSymlinks(path); err == nil {
		path = resolved
	}
	var stat syscall.Stat_t
	if syscall.Stat(path, &stat) == nil {
		uid, gid := int(stat.Uid), int(stat.Gid)
		volume.OwnerUID, volume.OwnerGID = &uid, &gid
		volume.Mode = fmt.Sprintf("%04o", stat.Mode&0o7777)
	}
	var fsStat syscall.Statfs_t
	if syscall.Statfs(path, &fsStat) == nil {
		volume.FreeMB = int64(fsStat.Bavail) * int64(fsStat.Bsize) >> 20
		volume.TotalMB = int64(fsStat.Blocks) * int64(fsStat.Bsize) >> 20
	}
	data, err := os.ReadFile("/proc/self/mountinfo")
	if err != nil {
		return
	}
	for _, line := range strings.Split(string(data), "\n") {
		fields := strings.Fields(line)
		separator := -1
		for i, field := range fields {
			if field == "-" {
				separator = i
				break
			}
		}
		if len(fields) < 6 || separator < 0 || separator+2 >= len(fields) {
			continue
		}
		mountPoint := unescapeMount(fields[4])
		if !withinMount(path, mountPoint) || len(mountPoint) < len(volume.MountPoint) {
			continue
		}
		volume.MountPoint = mountPoint
		volume.Filesystem = fields[separator+1]
		volume.Source = fields[separator+2]
		volume.ReadOnly = false
		for _, option := range strings.Split(fields[5], ",") {
			if option == "ro" {
				volume.ReadOnly = true
			}
		}
	}
}

func withinMount(path, mountPoint string) bool {
	return mountPoint == "/" || path == mountPoint || strings.HasPrefix(path, mountPoint+"/")
}

func unescapeMount(value string) string {
	return strings.NewReplacer(`\040`, " ", `\011`, "\t", `\012`, "\n", `\134`, `\`).Replace(value)
}
