package diagnostics

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/SteamServerUI/StationeersServerUI/v6/src/config"
)

const (
	lowDiskBytes     = 2 << 30
	accessScanLimit  = 20000
	accessIssueLimit = 10
)

type Environment struct {
	UID              int      `json:"uid"`
	GID              int      `json:"gid"`
	Groups           []int    `json:"groups,omitempty"`
	CPUs             int      `json:"cpus"`
	MemoryMB         int64    `json:"memoryMB,omitempty"`
	Kernel           string   `json:"kernel,omitempty"`
	ContainerSignals []string `json:"containerSignals,omitempty"`
	Volumes          []Volume `json:"volumes"`
}

type Volume struct {
	Name       string `json:"name"`
	Path       string `json:"path"`
	Status     string `json:"status"`
	MountPoint string `json:"mountPoint,omitempty"`
	Filesystem string `json:"filesystem,omitempty"`
	Source     string `json:"source,omitempty"`
	ReadOnly   bool   `json:"readOnly,omitempty"`
	FreeMB     int64  `json:"freeMB,omitempty"`
	TotalMB    int64  `json:"totalMB,omitempty"`
	OwnerUID   *int   `json:"ownerUid,omitempty"`
	OwnerGID   *int   `json:"ownerGid,omitempty"`
	Mode       string `json:"mode,omitempty"`
}

type AccessScan struct {
	Checked   int      `json:"checked"`
	Truncated bool     `json:"truncated,omitempty"`
	Problems  int      `json:"problems"`
	Issues    []string `json:"issues"`
}

func collectEnvironment() Environment {
	env := Environment{UID: -1, GID: -1, CPUs: runtime.NumCPU(), ContainerSignals: ContainerSignals(), Volumes: []Volume{}}
	if runtime.GOOS != "windows" {
		env.UID, env.GID = os.Getuid(), os.Getgid()
		env.Groups, _ = os.Getgroups()
	}
	platformEnvironment(&env)
	for _, entry := range []struct{ name, path string }{
		{"workdir", "."},
		{"ssui", config.GetSSUIFolder()},
		{"saves", "./saves"},
		{"logs", config.GetLogFolder()},
	} {
		volume := Volume{Name: entry.name, Path: entry.path, Status: inspectDirectory(entry.path)}
		if volume.Status != "missing" && volume.Status != "not configured" {
			describeVolume(&volume)
		}
		env.Volumes = append(env.Volumes, volume)
	}
	return env
}

func ContainerSignals() []string {
	signals := []string{}
	if value := strings.TrimSpace(os.Getenv("SSUI_CONTAINER")); value == "1" || strings.EqualFold(value, "true") {
		signals = append(signals, "SSUI_CONTAINER="+value)
	}
	if _, err := os.Stat("/.dockerenv"); err == nil {
		signals = append(signals, "/.dockerenv")
	}
	if _, err := os.Stat("/run/.containerenv"); err == nil {
		signals = append(signals, "/run/.containerenv")
	}
	if data, err := os.ReadFile("/proc/1/cgroup"); err == nil {
		for _, runtimeName := range []string{"docker", "containerd", "kubepods", "crio", "libpod"} {
			if strings.Contains(string(data), runtimeName) {
				signals = append(signals, "cgroup:"+runtimeName)
			}
		}
	}
	return signals
}

func inspectDirectory(path string) string {
	status := inspectPath(path, true)
	if status != "ok" {
		return status
	}
	if writeTest(path) != nil {
		return "not writable"
	}
	return "ok"
}

func writeTest(path string) error {
	file, err := os.CreateTemp(path, ".ssui-write-test-*")
	if err != nil {
		return err
	}
	name := file.Name()
	if err := file.Close(); err != nil {
		os.Remove(name)
		return err
	}
	return os.Remove(name)
}

func (report *Report) checkEnvironment() {
	for _, volume := range report.Environment.Volumes {
		if volume.ReadOnly && volume.Status == "ok" {
			volume.Status = "on a read-only mount"
		}
		if volume.Status != "ok" && volume.Status != "missing" && volume.Status != "not configured" {
			report.add(codeReadOnly, StatusFail, fmt.Sprintf("The %s folder (%s) is %s.", volume.Name, volume.Path, volume.Status), "Check the folder owner and permissions, and the container volume mount.")
		}
		if volume.TotalMB > 0 && volume.FreeMB*(1<<20) < lowDiskBytes {
			report.add(codeDisk, StatusWarn, fmt.Sprintf("Only %d MB are free for the %s folder.", volume.FreeMB, volume.Name), "Free up disk space or remove old backups.")
		}
		if volume.OwnerUID != nil && report.Environment.UID >= 0 && *volume.OwnerUID != report.Environment.UID && report.Environment.UID != 0 && volume.Status != "ok" {
			report.add(codeOwner, StatusWarn, fmt.Sprintf("The %s folder belongs to UID %d, but SSUI runs as UID %d.", volume.Name, *volume.OwnerUID, report.Environment.UID), "Change the folder owner, or run SSUI or the container as that user.")
		}
	}
}

func scanAccess(limit int) AccessScan {
	scan := AccessScan{Issues: []string{}}
	scan.tree(config.GetSSUIFolder(), true, limit)
	scan.tree("./saves", config.GetIsDockerContainer(), limit)
	return scan
}

func (scan *AccessScan) tree(root string, writable bool, limit int) {
	if strings.TrimSpace(root) == "" {
		return
	}
	if _, err := os.Stat(root); errors.Is(err, fs.ErrNotExist) {
		return
	}
	filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if limit > 0 && scan.Checked >= limit {
			scan.Truncated = true
			return filepath.SkipAll
		}
		if walkErr != nil {
			if !errors.Is(walkErr, fs.ErrNotExist) {
				scan.note(path, "cannot be read")
			}
			return nil
		}
		if entry.Type()&fs.ModeSymlink != 0 || (!entry.IsDir() && !entry.Type().IsRegular()) {
			return nil
		}
		scan.Checked++
		if entry.IsDir() {
			status := inspectPath(path, true)
			if status == "ok" && writable && writeTest(path) != nil {
				status = "not writable"
			}
			if status != "ok" {
				scan.note(path, status)
			}
			return nil
		}
		flag, problem := os.O_RDONLY, "cannot be read"
		if writable {
			flag, problem = os.O_RDWR, "cannot be opened for reading and writing"
		}
		if file, err := os.OpenFile(path, flag, 0); err != nil {
			scan.note(path, problem)
		} else {
			file.Close()
		}
		return nil
	})
}

func (scan *AccessScan) note(path, problem string) {
	scan.Problems++
	if len(scan.Issues) < accessIssueLimit {
		scan.Issues = append(scan.Issues, path+" "+problem)
	}
}
