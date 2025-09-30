package systeminfo

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

// OsStats represents the OS statistics that will be returned to the client
type OsStats struct {
	OSName          string  `json:"osName"`
	OSVersion       string  `json:"osVersion"`
	Kernel          string  `json:"kernel"`
	Uptime          string  `json:"uptime"`
	BackendIPAddr   string  `json:"backendIpAddress"`
	CPUUsage        float64 `json:"cpuUsage"`
	MemoryUsage     float64 `json:"memoryUsage"`
	DiskUsage       float64 `json:"diskUsage"`
	LastRefreshTime string  `json:"lastRefreshTime"`
}

var (
	osStatsMutex    sync.RWMutex
	cachedOsStats   *OsStats
	lastRefreshTime time.Time
	// Set cache duration to 1 minute
	cacheDuration = 60 * time.Second
)

// CPUStats holds CPU time statistics
type CPUStats struct {
	idle  uint64
	total uint64
}

// refreshCachedStats updates the OS stats in the cache
func RefreshCachedStats() (*OsStats, error) {
	osStatsMutex.Lock()
	defer osStatsMutex.Unlock()

	// Check if cache is still valid
	if cachedOsStats != nil && time.Since(lastRefreshTime) < cacheDuration {
		return cachedOsStats, nil
	}

	// Gather all OS statistics
	stats := &OsStats{}

	// Get OS name and version
	switch runtime.GOOS {
	case "windows":
		stats.OSName = "Windows"
		cmd := exec.Command("powershell", "-Command", "(Get-WmiObject -class Win32_OperatingSystem).Caption")
		output, err := cmd.Output()
		if err == nil {
			stats.OSVersion = strings.TrimSpace(string(output))
		} else {
			stats.OSVersion = "Unknown Windows Version"
		}
	case "linux":
		// Try to get OS name from /etc/os-release
		if data, err := os.ReadFile("/etc/os-release"); err == nil {
			lines := strings.Split(string(data), "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "NAME=") {
					stats.OSName = strings.Trim(strings.TrimPrefix(line, "NAME="), "\"")
				} else if strings.HasPrefix(line, "VERSION=") {
					stats.OSVersion = strings.Trim(strings.TrimPrefix(line, "VERSION="), "\"")
				}
			}
		}

		if stats.OSName == "" {
			stats.OSName = "Linux"
		}

		if stats.OSVersion == "" {
			// Try lsb_release if available
			cmd := exec.Command("lsb_release", "-d")
			output, err := cmd.Output()
			if err == nil {
				stats.OSVersion = strings.TrimSpace(strings.TrimPrefix(string(output), "Description:"))
			} else {
				stats.OSVersion = "Unknown Linux Version"
			}
		}
	default:
		stats.OSName = runtime.GOOS
		stats.OSVersion = "Unknown Version"
	}

	// Get kernel version
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("powershell", "-Command", "(Get-WmiObject Win32_OperatingSystem).Version")
		output, err := cmd.Output()
		if err == nil {
			stats.Kernel = strings.TrimSpace(string(output))
		} else {
			stats.Kernel = "Unknown"
		}
	case "linux":
		cmd := exec.Command("uname", "-r")
		output, err := cmd.Output()
		if err == nil {
			stats.Kernel = strings.TrimSpace(string(output))
		} else {
			stats.Kernel = "Unknown"
		}
	default:
		stats.Kernel = "Unknown"
	}

	// Get uptime
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("powershell", "-Command", "(Get-WmiObject Win32_OperatingSystem).LastBootUpTime")
		output, err := cmd.Output()
		if err == nil {
			bootTime := strings.TrimSpace(string(output))
			bootTimeFormat := "20060102150405.000000-070"
			t, err := time.Parse(bootTimeFormat, bootTime)
			if err == nil {
				uptime := time.Since(t)
				days := int(uptime.Hours() / 24)
				hours := int(uptime.Hours()) % 24
				mins := int(uptime.Minutes()) % 60

				if days > 0 {
					stats.Uptime = fmt.Sprintf("%d days, %d hours, %d minutes", days, hours, mins)
				} else if hours > 0 {
					stats.Uptime = fmt.Sprintf("%d hours, %d minutes", hours, mins)
				} else {
					stats.Uptime = fmt.Sprintf("%d minutes", mins)
				}
			} else {
				stats.Uptime = "Unknown"
			}
		} else {
			stats.Uptime = "Unknown"
		}
	case "linux":
		if data, err := os.ReadFile("/proc/uptime"); err == nil {
			fields := strings.Fields(string(data))
			if len(fields) > 0 {
				uptimeSec, err := strconv.ParseFloat(fields[0], 64)
				if err == nil {
					uptime := time.Duration(uptimeSec) * time.Second
					days := int(uptime.Hours() / 24)
					hours := int(uptime.Hours()) % 24
					mins := int(uptime.Minutes()) % 60

					if days > 0 {
						stats.Uptime = fmt.Sprintf("%d days, %d hours, %d minutes", days, hours, mins)
					} else if hours > 0 {
						stats.Uptime = fmt.Sprintf("%d hours, %d minutes", hours, mins)
					} else {
						stats.Uptime = fmt.Sprintf("%d minutes", mins)
					}
				} else {
					stats.Uptime = "Unknown"
				}
			} else {
				stats.Uptime = "Unknown"
			}
		} else {
			stats.Uptime = "Unknown"
		}
	default:
		stats.Uptime = "Unknown"
	}

	// Get backend IP address - try to get a non-loopback IP
	stats.BackendIPAddr = getIPAddress()

	// Get CPU usage
	stats.CPUUsage = getCPUUsage()

	// Get memory usage
	stats.MemoryUsage = getMemoryUsage()

	// Get disk usage
	stats.DiskUsage = getDiskUsage()

	// Update timestamp
	now := time.Now()
	stats.LastRefreshTime = now.Format(time.RFC3339)
	lastRefreshTime = now
	cachedOsStats = stats

	return stats, nil
}

// getIPAddress returns a non-loopback IP address of the machine
func getIPAddress() string {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("powershell", "-Command", "Get-NetIPAddress | Where-Object {$_.AddressFamily -eq 'IPv4' -and $_.IPAddress -ne '127.0.0.1'} | Select-Object -ExpandProperty IPAddress -First 1")
		output, err := cmd.Output()
		if err == nil && len(output) > 0 {
			return strings.TrimSpace(string(output))
		}
	case "linux":
		cmd := exec.Command("hostname", "-I")
		output, err := cmd.Output()
		if err == nil && len(output) > 0 {
			ips := strings.Fields(string(output))
			if len(ips) > 0 {
				return ips[0]
			}
		}
	}
	return "127.0.0.1" // Default to loopback if we can't find anything else
}

// getCPUUsage returns the CPU usage as a percentage
func getCPUUsage() float64 {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("powershell", "-Command", `(1..5 | ForEach-Object { (Get-CimInstance Win32_PerfFormattedData_PerfOS_Processor -Filter "Name='_Total'").PercentProcessorTime; Start-Sleep 1 } | Measure-Object -Average).Average`)
		output, err := cmd.Output()
		if err == nil {
			if usage, err := strconv.ParseFloat(strings.TrimSpace(string(output)), 64); err == nil {
				return usage
			}
		}
	case "linux":
		// First try using /proc/stat for accurate measurement across environments including containers
		// We need to sample twice with a small delay to calculate CPU usage
		prev, err := readCPUStats()
		if err == nil {
			// Small sleep to get a delta
			time.Sleep(200 * time.Millisecond)
			current, err := readCPUStats()
			if err == nil {
				// Calculate the delta between measurements
				idleDelta := current.idle - prev.idle
				totalDelta := current.total - prev.total

				if totalDelta > 0 {
					// CPU usage is the percentage of non-idle time
					return 100.0 * (1.0 - float64(idleDelta)/float64(totalDelta))
				}
			}
		}

		// Fallback to other methods if /proc/stat didn't work
		methods := []func() (float64, error){
			// Method 1: Using top
			func() (float64, error) {
				cmd := exec.Command("sh", "-c", "top -bn1 | grep 'Cpu(s)' | awk '{print $2 + $4}'")
				output, err := cmd.Output()
				if err != nil {
					return 0.0, err
				}
				return strconv.ParseFloat(strings.TrimSpace(string(output)), 64)
			},
			// Method 2: Using mpstat if available
			func() (float64, error) {
				cmd := exec.Command("sh", "-c", "command -v mpstat >/dev/null 2>&1 && mpstat 1 1 | awk '/Average:/ && $12 ~ /[0-9.]+/ {print 100 - $12}'")
				output, err := cmd.Output()
				if err != nil {
					return 0.0, err
				}
				return strconv.ParseFloat(strings.TrimSpace(string(output)), 64)
			},
			// Method 3: Using /proc/loadavg as a rough indicator
			func() (float64, error) {
				data, err := os.ReadFile("/proc/loadavg")
				if err != nil {
					return 0.0, err
				}
				fields := strings.Fields(string(data))
				if len(fields) < 1 {
					return 0.0, fmt.Errorf("invalid format in /proc/loadavg")
				}
				loadavg, err := strconv.ParseFloat(fields[0], 64)
				if err != nil {
					return 0.0, err
				}
				// Get number of CPUs to normalize the load
				var numCPU int
				if cpuinfo, err := os.ReadFile("/proc/cpuinfo"); err == nil {
					processors := strings.Count(string(cpuinfo), "processor")
					if processors > 0 {
						numCPU = processors
					} else {
						numCPU = runtime.NumCPU()
					}
				} else {
					numCPU = runtime.NumCPU()
				}
				// Convert load average to a percentage (capped at 100%)
				cpuUsage := (loadavg / float64(numCPU)) * 100
				if cpuUsage > 100 {
					cpuUsage = 100
				}
				return cpuUsage, nil
			},
		}

		// Try each method until one works
		for _, method := range methods {
			if usage, err := method(); err == nil {
				return usage
			}
		}
	}
	return 0.0 // Default to 0 if we can't determine the CPU usage
}

// readCPUStats reads the current CPU statistics from /proc/stat
func readCPUStats() (CPUStats, error) {
	data, err := os.ReadFile("/proc/stat")
	if err != nil {
		return CPUStats{}, err
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 5 && fields[0] == "cpu" {
			// CPU line format: cpu user nice system idle iowait irq softirq steal guest guest_nice
			// We need to sum all fields except idle (field 4) to get total CPU time
			var total uint64
			idle, _ := strconv.ParseUint(fields[4], 10, 64)

			for i := 1; i < len(fields); i++ {
				val, _ := strconv.ParseUint(fields[i], 10, 64)
				total += val
			}

			return CPUStats{idle: idle, total: total}, nil
		}
	}

	return CPUStats{}, fmt.Errorf("failed to parse CPU stats from /proc/stat")
}

// getMemoryUsage returns the systems memory usage in MB
func getMemoryUsage() float64 {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("powershell", "-Command", "[math]::Round((((Get-WmiObject Win32_OperatingSystem).TotalVisibleMemorySize - (Get-WmiObject Win32_OperatingSystem).FreePhysicalMemory) /1024 ) )")
		output, err := cmd.Output()
		if err == nil {
			if usage, err := strconv.ParseFloat(strings.TrimSpace(string(output)), 64); err == nil {
				return usage
			}
		}
	case "linux":
		cmd := exec.Command("sh", "-c", "free | grep Mem | awk '{print ($3) / 1024}'")
		output, err := cmd.Output()
		if err == nil {
			if usage, err := strconv.ParseFloat(strings.TrimSpace(string(output)), 64); err == nil {
				return usage
			}
		}
	}
	return 0.0 // Default to 0 if we can't determine the memory usage
}

// getDiskUsage returns the disk usage as a percentage
func getDiskUsage() float64 {
	var path string
	if runtime.GOOS == "windows" {
		path = "C:\\"
	} else {
		path = "/"
	}

	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("powershell", "-Command", `[math]::Round(((Get-CimInstance Win32_LogicalDisk -Filter "DeviceID='$((Get-Location).Drive.Name):'").Size - (Get-CimInstance Win32_LogicalDisk -Filter "DeviceID='$((Get-Location).Drive.Name):'").FreeSpace) / (Get-CimInstance Win32_LogicalDisk -Filter "DeviceID='$((Get-Location).Drive.Name):'").Size * 100, 2)`)
		output, err := cmd.Output()
		if err == nil {
			if usage, err := strconv.ParseFloat(strings.TrimSpace(string(output)), 64); err == nil {
				return usage
			}
		}
	case "linux":
		cmd := exec.Command("df", "-h", path)
		output, err := cmd.Output()
		if err == nil {
			lines := strings.Split(string(output), "\n")
			if len(lines) >= 2 {
				fields := strings.Fields(lines[1])
				if len(fields) >= 5 {
					usage := strings.TrimSuffix(fields[4], "%")
					if parsed, err := strconv.ParseFloat(usage, 64); err == nil {
						return parsed
					}
				}
			}
		}
	}
	return 0.0 // Default to 0 if we can't determine the disk usage
}
